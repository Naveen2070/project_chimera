package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"project_chimera/gene_bank_service/internal/dto"
	"project_chimera/gene_bank_service/pkg/common"

	"github.com/gofiber/fiber/v2"
)

// Helper function to process response data
func ProcessFloraData(rawDataList []interface{}) ([]dto.FloraData, error) {
	var floraList []dto.FloraData

	for _, item := range rawDataList {
		if itemStr, ok := item.(string); ok {
			var floraData dto.FloraData
			err := json.Unmarshal([]byte(itemStr), &floraData) // Convert JSON string to struct
			if err != nil {
				log.Printf("Error unmarshaling JSON string to FloraData: %v", err)
				return nil, fmt.Errorf("failed to map JSON string to FloraData: %v", err)
			}
			floraList = append(floraList, floraData)
		} else {
			log.Printf("Unexpected item type in data: %T, Value: %v", item, item)
			return nil, errors.New("invalid item type in RPC response data")
		}
	}

	return floraList, nil
}

// Helper function to handle RPC error responses
func HandleRPCError(res common.MessageResponse) error {
	var msg string
	for _, item := range res.Data {
		if itemStr, ok := item.(string); ok {
			msg = itemStr
		}
	}
	return &fiber.Error{Code: int(res.Code), Message: msg}
}
