import { $Enums } from '@prisma/client';

export class FloraUpstream {
  id?: string;
  common_name!: string;
  scientific_name!: string;
  user_id!: string;
  Image!: Uint8Array;
  Description!: string;
  Origin!: string;
  OtherDetails!: object;
  type!: $Enums.PostType;
}
