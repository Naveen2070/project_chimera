// Copyright 2025 Naveen R
//
//	Licensed under the Apache License, Version 2.0 (the "License");
//	you may not use this file except in compliance with the License.
//	You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
//	Unless required by applicable law or agreed to in writing, software
//	distributed under the License is distributed on an "AS IS" BASIS,
//	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//	See the License for the specific language governing permissions and
//	limitations under the License.
package utils

import (
	"encoding/base64"
	"fmt"
	"project_chimera/gene_bank_service/internal/dto"
)

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

func MapToFlora(data map[string]interface{}) (dto.Flora, error) {
	var flora dto.Flora

	// ID
	if id, ok := data["id"].(string); ok {
		flora.ID = id
	}

	// CommonName
	if commonName, ok := data["common_name"].(string); ok {
		flora.CommonName = commonName
	}

	// ScientificName
	if scientificName, ok := data["scientific_name"].(string); ok {
		flora.ScientificName = scientificName
	}

	// Image
	if imageBase64, ok := data["image"].(string); ok {
		image, err := base64.StdEncoding.DecodeString(imageBase64)
		if err != nil {
			return flora, fmt.Errorf("failed to decode image: %v", err)
		}
		flora.Image = image
	}

	// Description
	if description, ok := data["description"].(string); ok {
		flora.Description = description
	}

	// Origin
	if origin, ok := data["origin"].(string); ok {
		flora.Origin = origin
	}

	// OtherDetails
	if otherDetails, ok := data["other_details"].(string); ok {
		flora.OtherDetails = otherDetails
	}

	// Type
	if typeData, ok := data["type"].(map[string]interface{}); ok {
		var floraType dto.Type
		if typeName, ok := typeData["name"].(string); ok {
			if typeName == string(dto.Public) || typeName == string(dto.Private) {
				floraType = dto.Type(typeName)
			} else {
				return flora, fmt.Errorf("invalid type name: expected 'public' or 'private', got '%s'", typeName)
			}
		}
		flora.Type = floraType
	}

	return flora, nil
}
