from typing import Any, Dict
from odmantic import Model


class FloraMongo(Model):
    flora_id: str
    Image: bytes  # Representing binary data
    Description: str
    Origin: str
    OtherDetails: Dict[str, Any]
    model_config = {"collection": "floras"}
