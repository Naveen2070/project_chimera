package flora

import (
	"errors"
	"log"
	"project_chimera/error_handle_service/pkg/common"
	logger "project_chimera/error_handle_service/pkg/logger"
	"project_chimera/error_handle_service/pkg/models"
)

func AutoFixFlora(floraResp models.FloraResponse) (models.FloraResponse, error) {
	fixedFloraData := floraResp

	eventType := floraResp.Pattern
	log.Print(eventType)
	switch eventType {
	case "flora.created":
		field := common.ExtractFieldNameFromError(floraResp.Data.Data.Error)
		log.Print(field)
		if field == "type" {
			logger.LogInfo("Fixing type field in flora.created event")
			fixedFloraData.Data.Data.Values.Type = "offline"
			return fixedFloraData, nil
		}
	}
	return fixedFloraData, errors.New("error")
}
