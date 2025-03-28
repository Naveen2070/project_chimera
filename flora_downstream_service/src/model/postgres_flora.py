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

import uuid
from sqlalchemy import Column, String, Enum
from sqlalchemy.dialects.postgresql import UUID
from sqlalchemy.ext.declarative import declarative_base
from enum import Enum as en

Base = declarative_base()


class Type(str, en):
    public = "public"
    private = "private"


class FloraPG(Base):
    __tablename__ = "Flora"
    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    user_id = Column(String, unique=True, index=True)
    common_name = Column(String, unique=True, index=True)
    scientific_name = Column(String, unique=True, index=True)
    type = Column(Enum(Type), nullable=False)

    def __repr__(self):
        return (
            f"FloraPG(id={self.id}, user_id={self.user_id}, "
            f"common_name={self.common_name}, scientific_name={self.scientific_name}, "
            f"type={self.type})"
        )
