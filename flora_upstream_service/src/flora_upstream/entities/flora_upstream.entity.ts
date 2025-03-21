import { $Enums, Flora } from '@prisma/client';

export class FloraUpstream implements Flora {
  id: string;
  common_name: string;
  scientific_name: string;
  user_id: string;
  type: $Enums.PostType;
}
