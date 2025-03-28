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

import json
from typing import Dict, Optional


class Pattern:
    def __init__(self, cmd: str):
        """
        Initialize the Pattern object.
        :param cmd: The command for routing.
        """
        self.cmd = cmd


class RpcMessage:
    def __init__(self, pattern: Pattern, data: Optional[Dict] = None):
        """
        Initialize the RpcMessage object.
        :param pattern: The routing pattern (as a Pattern object).
        :param data: The payload data.
        """
        self.pattern = pattern
        self.data = data or {}

    @classmethod
    def from_json(cls, json_str: str):
        """
        Create an RpcMessage object from a JSON string.
        :param json_str: The JSON string representing the message.
        :return: An RpcMessage object.
        """
        parsed_data = json.loads(json_str)
        pattern_data = parsed_data.get("pattern", {})
        pattern = Pattern(cmd=pattern_data.get("cmd", ""))
        return cls(pattern=pattern, data=parsed_data.get("data", {}))

    def to_json(self) -> str:
        """ "
        Convert the RpcMessage object to a JSON string.
        :return: The JSON string representation of the message.
        """
        return json.dumps({"pattern": {"cmd": self.pattern.cmd}, "data": self.data})

    def __str__(self):
        return self.to_json()

    def __repr__(self):
        return f"RpcMessage(pattern={self.pattern}, data={self.data})"
