//	Copyright 2025 Naveen R
//
//		Licensed under the Apache License, Version 2.0 (the "License");
//		you may not use this file except in compliance with the License.
//		You may obtain a copy of the License at
//
//		http://www.apache.org/licenses/LICENSE-2.0
//
//		Unless required by applicable law or agreed to in writing, software
//		distributed under the License is distributed on an "AS IS" BASIS,
//		WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//		See the License for the specific language governing permissions and
//		limitations under the License.

package common

type SuccessResponse struct {
	Status string `json:"status"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// Define the top-level structure for the message
type MessageRequest struct {
	Pattern Pattern     `json:"pattern"`
	Data    interface{} `json:"data"`
}

type MessageResponse struct {
	Status string        `json:"status"`
	Data   []interface{} `json:"data"`
}

// Define the "Pattern" structure
type Pattern struct {
	Cmd string `json:"cmd"`
}
