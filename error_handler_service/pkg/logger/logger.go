// Copyright 2025 Naveen R
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package customlogger

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/middleware/logger"
)

// getLogFilePath generates a file path like "logs/routes/2025-04-03.log"
func getLogFilePath(logType string) string {
	// Get current date in YYYY-MM-DD format
	currentDate := time.Now().Format("2006-01-02")

	// Ensure logs directory exists
	dirPath := fmt.Sprintf("logs/%s", logType)

	// Create logs directory if it doesn't exist
	ensureDirectoryExists(dirPath)

	// Return log file path
	return fmt.Sprintf("%s/%s.log", dirPath, currentDate)
}

// InitLogger initializes and returns the Fiber logger middleware
func InitLogger() logger.Config {
	// Ensure logs directory exists
	createBaseLogDirectories()

	// Open daily log file for routes
	routeLogFile := openLogFile(getLogFilePath("routes"))

	// Create multi-writer for console + route log file
	routeWriter := io.MultiWriter(os.Stdout, routeLogFile)

	// Create Fiber logger configuration
	return logger.Config{
		Format:     "[${time}] ${ip}:${port} | ${status} | ${method} | ${path} | ${latency}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
		Output:     routeWriter,
	}
}

// LogInfo writes info logs to a daily logs/info/YYYY-MM-DD.log
func LogInfo(message string) {
	infoLogFile := openLogFile(getLogFilePath("info"))
	defer infoLogFile.Close()

	infoWriter := io.MultiWriter(os.Stdout, infoLogFile)
	log.SetOutput(infoWriter)
	log.Println("[INFO]", message)
}

// LogError writes error logs to a daily logs/error/YYYY-MM-DD.log
func LogError(message string) {
	errorLogFile := openLogFile(getLogFilePath("error"))
	defer errorLogFile.Close()

	errorWriter := io.MultiWriter(os.Stdout, errorLogFile)
	log.SetOutput(errorWriter)
	log.Println("[ERROR]", message)
}

// LogWarning writes warning logs to a daily logs/warning/YYYY-MM-DD.log
func LogWarning(message string) {
	warningLogFile := openLogFile(getLogFilePath("warning"))
	defer warningLogFile.Close()

	warningWriter := io.MultiWriter(os.Stdout, warningLogFile)
	log.SetOutput(warningWriter)
	log.Println("[WARNING]", message)
}

// LogFatal writes fatal logs to a daily logs/fatal/YYYY-MM-DD.log
func LogFatal(message string) {
	fatalLogFile := openLogFile(getLogFilePath("fatal"))
	defer fatalLogFile.Close()

	fatalWriter := io.MultiWriter(os.Stdout, fatalLogFile)
	log.SetOutput(fatalWriter)

	log.Fatalf("[FATAL] %s", message)
}

// createBaseLogDirectories ensures the logs/ directory and its subdirectories exist
func createBaseLogDirectories() {
	baseDirs := []string{"logs/routes", "logs/info", "logs/error", "logs/fatal"}
	for _, dir := range baseDirs {
		ensureDirectoryExists(dir)
	}
}

// ensureDirectoryExists creates a directory if it does not exist
func ensureDirectoryExists(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0755)
	}
}

// openLogFile opens (or creates) a daily log file
func openLogFile(filePath string) *os.File {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file %s: %v", filePath, err)
	}
	return file
}

const dateLayout = "2006-01-02"

func ArchiveOldLogs() {
	logTypes := []string{"info", "error", "warning", "fatal", "routes"}
	cutoff := time.Now().AddDate(0, 0, -15)

	for _, logType := range logTypes {
		dir := filepath.Join("logs", logType)
		archiveDir := filepath.Join(dir, "archive")
		ensureDirectoryExists(archiveDir)

		filesGrouped := make(map[string][]string)

		_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() || !strings.HasSuffix(info.Name(), ".log") {
				return nil
			}

			baseName := strings.TrimSuffix(info.Name(), ".log")
			fileDate, err := time.Parse(dateLayout, baseName)
			if err != nil || !fileDate.Before(cutoff) {
				return nil
			}

			bucketStart := fileDate.AddDate(0, 0, -((fileDate.Day() - 1) % 15))
			bucketEnd := bucketStart.AddDate(0, 0, 14)
			rangeKey := fmt.Sprintf("%s-%s", bucketStart.Format("02-01-2006"), bucketEnd.Format("02-01-2006"))

			filesGrouped[rangeKey] = append(filesGrouped[rangeKey], path)
			return nil
		})

		for rangeKey, paths := range filesGrouped {
			zipFilePath := filepath.Join(archiveDir, fmt.Sprintf("%s.zip", rangeKey))

			err := zipFiles(zipFilePath, paths)
			if err != nil {
				LogError(fmt.Sprintf("Failed to zip logs for %s: %v", rangeKey, err))
				continue
			}

			for _, f := range paths {
				if err := os.Remove(f); err != nil {
					LogWarning(fmt.Sprintf("Failed to remove original file %s: %v", f, err))
				}
			}
			LogInfo(fmt.Sprintf("Archived logs to %s", zipFilePath))
		}
	}
}

func zipFiles(zipPath string, files []string) error {
	newZipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	sort.Strings(files)

	for _, file := range files {
		if err := addFileToZip(zipWriter, file); err != nil {
			return err
		}
	}
	return nil
}

func addFileToZip(zipWriter *zip.Writer, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, filename := filepath.Split(filePath)
	writer, err := zipWriter.Create(filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, file)
	return err
}
