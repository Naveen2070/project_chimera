import { PartialType } from '@nestjs/mapped-types';
import { CreateFloraUpstreamDto } from './create-flora_upstream.dto';

export class UpdateFloraUpstreamDto extends PartialType(
  CreateFloraUpstreamDto,
) {
  id: number;
}
