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
