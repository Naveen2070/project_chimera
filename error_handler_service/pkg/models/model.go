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

package models

import (
	"time"
)

// FloraData represents the structure of the "values" field in your JSON
type FloraData struct {
	CommonName     string                 `bson:"common_name"`
	ScientificName string                 `bson:"scientific_name"`
	UserID         string                 `bson:"user_id"`
	Type           string                 `bson:"type"`
	Image          string                 `bson:"Image"`
	Description    string                 `bson:"Description"`
	Origin         string                 `bson:"Origin"`
	OtherDetails   map[string]interface{} `bson:"OtherDetails"`
	CreatedAt      time.Time              `bson:"createdAt"` // Automatically generated timestamp
}

// FloraResponse represents the overall structure of the response
type FloraResponse struct {
	Data struct {
		Code int `json:"code"`
		Data struct {
			Values FloraData `json:"values"`
		} `json:"data"`
		Error  string `json:"error"`
		Status string `json:"status"`
		Type   string `json:"type"`
	} `json:"data"`
	Pattern string `json:"pattern"`
}
