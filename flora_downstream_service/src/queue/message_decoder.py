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
