import { z } from 'zod';

export const NotificationResponseSchema = z.object({
  type: z.string(),
  status: z.string(),
  code: z.number(),
  data: z.string(),
});

export type Notification = z.infer<typeof NotificationResponseSchema>;

export class NotificationResponse {
  constructor(data: unknown) {
    try {
      const parsedData = NotificationResponseSchema.parse(data);
      Object.assign(this, parsedData);
    } catch (error) {
      console.error('Invalid notification data:', error);
      throw new Error('Invalid notification data', { cause: error });
    }
  }
}
