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
	CommonName     string                 `bson:"common_name" json:"common_name"`
	ScientificName string                 `bson:"scientific_name" json:"scientific_name"`
	UserID         string                 `bson:"user_id" json:"user_id"`
	Type           string                 `bson:"type" json:"type"`
	Image          string                 `bson:"image" json:"image"`
	Description    string                 `bson:"description" json:"description"`
	Origin         string                 `bson:"origin" json:"origin"`
	OtherDetails   map[string]interface{} `bson:"other_details" json:"OtherDetails"`
	CreatedAt      time.Time              `bson:"created_at" json:"createdAt"`
}

type FloraResponse struct {
	Data struct {
		Code   int    `json:"code"`
		Status string `json:"status"`
		Type   string `json:"type"`
		Data   struct {
			Values FloraData `json:"values"`
			ID     string    `json:"id"`
			Error  string    `json:"error"`
		} `json:"data"`
	} `json:"data"`
	Pattern string `json:"pattern"`
}
