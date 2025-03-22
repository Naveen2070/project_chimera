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

export function toFloraMongo(upstream: FloraUpstream): Flora {
  return {
    flora_id: upstream.id as string,
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
