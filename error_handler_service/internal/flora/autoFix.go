package flora

import (
	"errors"
	"log"
	"project_chimera/error_handle_service/pkg/common"
	"project_chimera/error_handle_service/pkg/models"
)

func AutoFixFlora(floraResp models.FloraResponse) (models.FloraResponse, error) {
	fixedFloraData := floraResp

	eventType := floraResp.Pattern

	switch eventType {
	case "flora.created":
		field := common.ExtractFieldNameFromError(floraResp.Data.Data.Error)
		if field == "type" {
			log.Println("type error")
			fixedFloraData.Data.Data.Values.Type = "offline"
			return fixedFloraData, nil
		}
	}
	return fixedFloraData, errors.New("error")
}
