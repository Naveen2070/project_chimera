//	Copyright 2025 Naveen R
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
using System.Text.Json.Serialization;
using System.Text.Json;

namespace user_service.Model
{
    public class ErrorDTO
    {
        public string Pattern { get; set; }
        public ResponseData Data { get; set; }

        public ErrorDTO(string pattern, ResponseData data)
        {
            Pattern = pattern;
            Data = data;
        }

        public class ResponseData
        {
            public int Code { get; set; }
            public string Status { get; set; }
            public string Type { get; set; }
            public object Data { get; set; }

            public ResponseData(int code, string status, string type, object data)
            {
                Code = code;
                Status = status;
                Type = type;
                Data = data;
            }
        }

        public static ErrorDTO CreateErrorDTO(string pattern, int code, string status, string type, object data)
        {
            var responseData = new ResponseData(code, status, type, data);
            return new ErrorDTO(pattern, responseData);
        }

        public override string ToString()
        {
            try
            {
                return JsonSerializer.Serialize(this, new JsonSerializerOptions
                {
                    WriteIndented = false,
                    DefaultIgnoreCondition = JsonIgnoreCondition.WhenWritingNull
                });
            }
            catch (Exception ex)
            {
                Console.Error.WriteLine("Failed to convert ErrorDataDTO to JSON: " + ex.Message);
                return base.ToString() ?? string.Empty;
            }
        }
    }
}
