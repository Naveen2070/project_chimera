from typing import Any, Dict
from odmantic import Model


class FloraMongo(Model):
    flora_id: str
    Image: bytes  # Representing binary data
    Description: str
    Origin: str
    OtherDetails: Dict[str, Any]
    model_config = {"collection": "floras"}

    def __repr__(self):
        return (
            f"FloraMongo(flora_id={self.flora_id}, Image={self.Image}, "
            f"Description={self.Description}, Origin={self.Origin}, "
            f"OtherDetails={self.OtherDetails})"
        )
