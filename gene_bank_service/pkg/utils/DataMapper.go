package utils

import "project_chimera/gene_bank_service/internal/dto"

func CreateFloraDataMap(payload dto.FloraUpdateRequest, userId string, imageBytes []byte) map[string]interface{} {
	data := make(map[string]interface{})

	if payload.ID != "" {
		data["ID"] = payload.ID
	}

	if payload.CommonName != "" {
		data["CommonName"] = payload.CommonName
	}
	if payload.ScientificName != "" {
		data["ScientificName"] = payload.ScientificName
	}
	if len(imageBytes) > 0 {
		data["Image"] = imageBytes
	}
	if payload.Description != "" {
		data["Description"] = payload.Description
	}
	if payload.Origin != "" {
		data["Origin"] = payload.Origin
	}
	if len(payload.OtherDetails) > 0 {
		data["OtherDetails"] = payload.OtherDetails
	}
	if payload.Type != "" {
		data["Type"] = payload.Type
	}
	if userId != "" {
		data["UserId"] = userId
	}

	return data
}
