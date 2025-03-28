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

from typing import Any, Dict, List, Optional
import json
from typing import Any, Dict, List, Optional


class FloraResponse:
    def __init__(self, code: int, data: Optional[List[Dict[str, Any]]]):
        self.code = code
        self.data = data

    def to_dict(self) -> Dict[str, Any]:
        return {"code": self.code, "data": self.data}

    def to_json(self) -> str:
        return json.dumps(self.to_dict())

    def __str__(self):
        return self.to_json()

    def __repr__(self):
        return self.to_json()


class ApiResponse:
    def __init__(
        self, status: str, code: int, data: Optional[List[Dict[str, Any]]] = None
    ):
        self.status = status
        self.code = code
        self.data = data or []

    def to_dict(self) -> Dict[str, Any]:
        return {"status": self.status, "code": self.code, "data": self.data}

    def to_json(self) -> str:
        return json.dumps(self.to_dict())

    def __str__(self):
        return self.to_json()

    def __repr__(self):
        return self.to_json()
