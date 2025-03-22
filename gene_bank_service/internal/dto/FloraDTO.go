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
package dto

type Flora struct {
	ID             string `json:"id,omitempty"`              // Unique identifier for the plant
	CommonName     string `json:"common_name,omitempty"`     // Common name of the plant
	ScientificName string `json:"scientific_name,omitempty"` // Scientific name of the plant
	Image          []byte `json:"image,omitempty"`           // Image data (bytes)
	Description    string `json:"description,omitempty"`     // Description of the plant
	Origin         string `json:"origin,omitempty"`          // Origin of the plant
	OtherDetails   string `json:"other_details,omitempty"`   // Additional details about the plant
	Type           string `json:"type,omitempty"`            // Type of post
}

type FloraResponse struct {
	Flora []Flora `json:"flora,omitempty"`
}

type FloraRequest struct {
	CommonName     string                 `json:"common_name,omitempty"`     // Common name of the plant
	ScientificName string                 `json:"scientific_name,omitempty"` // Scientific name of the plant
	ImageURL       string                 `json:"image_url,omitempty"`       // Image URL or file reference
	ImagePath      string                 `json:"image_path,omitempty"`      // Image URL or file reference
	Description    string                 `json:"description,omitempty"`     // Description of the plant
	Origin         string                 `json:"origin,omitempty"`          // Origin of the plant
	OtherDetails   map[string]interface{} `json:"other_details,omitempty"`   // Additional details about the plant
	Type           string                 `json:"type,omitempty"`            // Type of post
}

type FloraUpdateRequest struct {
	ID             string                 `json:"id"`                        // Unique identifier for the plant
	CommonName     string                 `json:"common_name,omitempty"`     // Common name of the plant
	ScientificName string                 `json:"scientific_name,omitempty"` // Scientific name of the plant
	ImageURL       string                 `json:"image_url,omitempty"`       // Image URL or file reference
	ImagePath      string                 `json:"image_path,omitempty"`      // Image URL or file reference
	Image          []byte                 `json:"image,omitempty"`           // Image data (bytes)
	Description    string                 `json:"description,omitempty"`     // Description of the plant
	Origin         string                 `json:"origin,omitempty"`          // Origin of the plant
	OtherDetails   map[string]interface{} `json:"other_details,omitempty"`   // Additional details about the plant
	Type           string                 `json:"type,omitempty"`            // Type of post
}
