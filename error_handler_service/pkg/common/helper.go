package common

import (
	"project_chimera/error_handle_service/pkg/models"
	"regexp"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func FloraResponseToBson(body models.FloraResponse) bson.D {
	floraData := body.Data.Data.Values

	return bson.D{
		{Key: "common_name", Value: floraData.CommonName},
		{Key: "scientific_name", Value: floraData.ScientificName},
		{Key: "user_id", Value: floraData.UserID},
		{Key: "type", Value: floraData.Type},
		{Key: "image", Value: floraData.Image},
		{Key: "description", Value: floraData.Description},
		{Key: "origin", Value: floraData.Origin},
		{Key: "other_details", Value: floraData.OtherDetails},
		{Key: "created_at", Value: time.Now().UTC()},

		// Include additional metadata from FloraResponse
		{Key: "status", Value: body.Data.Status},
		{Key: "error", Value: body.Data.Data.Error},
		{Key: "id", Value: body.Data.Data.ID},
		{Key: "response_type", Value: body.Data.Type},
		{Key: "pattern", Value: body.Pattern},
	}
}

// Function to extract field name from an error string
// ExtractFieldNameFromError attempts to find the field name causing a Prisma error
func ExtractFieldNameFromError(errorString string) string {
	// Try to find 'Invalid value for argument ...'
	reMain := regexp.MustCompile(`(?i)Invalid value for argument [\` + "`" + `"']?(\w+)[\` + "`" + `"']?`)
	match := reMain.FindStringSubmatch(errorString)
	if len(match) > 1 {
		return match[1]
	}

	// Fallback: try to locate the line where the type/value mismatch occurs
	lines := strings.Split(errorString, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasSuffix(line, "~") || strings.Contains(line, "~") {
			// Try to extract field name from line
			colonSplit := strings.Split(line, ":")
			if len(colonSplit) > 0 {
				fieldLine := strings.TrimSpace(colonSplit[0])
				fieldParts := strings.Split(fieldLine, "\"")
				if len(fieldParts) >= 2 {
					return fieldParts[1] // e.g., `"type": ...` → "type"
				}
			}
		}
	}

	return ""
}
