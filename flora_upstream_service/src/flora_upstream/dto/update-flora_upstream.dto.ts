import { PartialType } from '@nestjs/mapped-types';
import { FloraPg } from './create-flora_upstream.dto';

export class UpdateFloraUpstreamDto extends PartialType(FloraPg) {
  id!: string;
}
