import { $Enums } from '@prisma/client';

export class RabbitMqPayload {
  id?: string; // Unique identifier for the plant
  CommonName!: string; // Common name of the plant
  ScientificName!: string; // Scientific name of the plant
  Image!: Uint8Array; // Image data (bytes)
  Description!: string; // Description of the plant
  Origin!: string; // Origin of the plant
  OtherDetails!: object; // Additional details about the plant
  Type!: $Enums.PostType; // Type of post
  UserId!: string; // User ID
}
