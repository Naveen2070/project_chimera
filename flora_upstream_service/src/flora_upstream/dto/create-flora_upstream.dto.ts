import { $Enums, Prisma } from '@prisma/client';

export class FloraPg implements Prisma.FloraCreateInput {
  id?: string | undefined;
  common_name!: string;
  scientific_name!: string;
  user_id!: string;
  type!: $Enums.PostType;
}
