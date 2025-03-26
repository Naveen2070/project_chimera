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
