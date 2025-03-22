//	Copyright 2025 Naveen R
//
//		Licensed under the Apache License, Version 2.0 (the "License");
//		you may not use this file except in compliance with the License.
//		You may obtain a copy of the License at
//
//		http://www.apache.org/licenses/LICENSE-2.0
//
//		Unless required by applicable law or agreed to in writing, software
//		distributed under the License is distributed on an "AS IS" BASIS,
//		WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//		See the License for the specific language governing permissions and
//		limitations under the License.
import { FloraPg } from 'src/flora_upstream/dto/create-flora_upstream.dto';
import { FloraUpstream } from 'src/flora_upstream/entities/flora_upstream.entity';
import { Flora } from 'src/flora_upstream/schema/flora.schema';

export function toFloraPg(upstream: FloraUpstream): FloraPg {
  return {
    common_name: upstream.common_name,
    scientific_name: upstream.scientific_name,
    user_id: upstream.user_id,
    type: upstream.type,
  };
}

export function toFloraMongo(id: string, upstream: FloraUpstream): Flora {
  return {
    flora_id: id,
    Image: Buffer.from(upstream.Image),
    Description: upstream.Description,
    Origin: upstream.Origin,
    OtherDetails: upstream.OtherDetails,
  };
}

export function toFloraUpstream(
  pgData: FloraPg,
  mongoData: Flora,
): FloraUpstream {
  return {
    id: mongoData.flora_id,
    common_name: pgData.common_name,
    scientific_name: pgData.scientific_name,
    user_id: pgData.user_id,
    type: pgData.type,
    Image: mongoData.Image,
    Description: mongoData.Description,
    Origin: mongoData.Origin,
    OtherDetails: mongoData.OtherDetails,
  };
}
