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
	"fmt"
	"io"
	"log"
	"os"
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
