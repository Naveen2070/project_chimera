package common

import (
	"project_chimera/error_handle_service/pkg/models"
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
