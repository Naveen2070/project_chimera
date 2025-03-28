# 	Copyright 2025 Naveen R
#
# 		Licensed under the Apache License, Version 2.0 (the "License");
# 		you may not use this file except in compliance with the License.
# 		You may obtain a copy of the License at
#
# 		http://www.apache.org/licenses/LICENSE-2.0
#
# 		Unless required by applicable law or agreed to in writing, software
# 		distributed under the License is distributed on an "AS IS" BASIS,
# 		WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# 		See the License for the specific language governing permissions and
# 		limitations under the License.

import base64
from typing import Any, Dict
from pydantic import BaseModel
from enum import Enum

from src.model.postgres_flora import Type


# Unified Flora Class
class Flora(BaseModel):
    id: str
    user_id: str
    common_name: str
    scientific_name: str
    type: Type
    image: bytes
    description: str
    origin: str
    other_details: Dict[str, Any]

    class Config:
        json_encoders = {
            Type: lambda t: t.value,
            bytes: lambda b: base64.b64encode(b).decode("utf-8"),
        }

    def __repr__(self):
        return (
            f"Flora(id={self.id}, user_id={self.user_id}, "
            f"common_name={self.common_name}, scientific_name={self.scientific_name}, "
            f"type={self.type}, description={self.description}, "
            f"origin={self.origin}, other_details={self.other_details})"
        )
