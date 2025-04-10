import { z } from 'zod';

// Define the NotificationResponse schema
export const NotificationResponseSchema = z.object({
  type: z.string(),
  status: z.string(),
  code: z.number(),
  data: z.record(z.any()), // Allows flexible data properties, replace with a stricter schema if needed
});

// Define the Notification type
export type Notification = z.infer<typeof NotificationResponseSchema>;

// NotificationResponse class
export class NotificationResponse {
  type!: string;
  status!: string;
  code!: number;
  data!: Record<string, unknown>;

  constructor(data: unknown) {
    try {
      // Validate and parse the input data
      const parsedData = NotificationResponseSchema.parse(data);
      Object.assign(this, parsedData);
    } catch (error) {
      console.error('Invalid notification data:', error);
      throw new Error(
        `Invalid notification data: ${error instanceof Error ? error.message : String(error)}`,
      );
    }
  }
}
