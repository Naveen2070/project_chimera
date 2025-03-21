import { Prop, Schema, SchemaFactory } from '@nestjs/mongoose';
import { HydratedDocument } from 'mongoose';

export type FloraDocument = HydratedDocument<Flora>;

@Schema()
export class Flora {
  @Prop({ required: true })
  flora_id!: string;

  @Prop({ type: Buffer }) // Use Buffer for binary data
  Image!: Buffer;

  @Prop({ required: true })
  Description!: string;

  @Prop({ required: true })
  Origin!: string;

  @Prop({ type: Object })
  OtherDetails!: Record<string, any>;
}

export const FloraSchema = SchemaFactory.createForClass(Flora);
