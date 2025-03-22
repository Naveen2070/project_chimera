import { Controller } from '@nestjs/common';
import {
  Ctx,
  MessagePattern,
  Payload,
  RmqContext,
} from '@nestjs/microservices';
import { FloraUpstreamService } from './flora_upstream.service';
import { RabbitMqPayload } from './dto/rabbit-payload';
import { FloraUpstream } from './entities/flora_upstream.entity';

@Controller()
export class FloraUpstreamController {
  constructor(private readonly floraUpstreamService: FloraUpstreamService) {}

  @MessagePattern({ cmd: 'add_flora' })
  async create(@Payload() data: RabbitMqPayload, @Ctx() context: RmqContext) {
    const channel = context.getChannelRef();
    const originalMsg = context.getMessage();

    const createFloraUpstreamDto: FloraUpstream = {
      common_name: data.CommonName,
      scientific_name: data.ScientificName,
      user_id: data.UserId,
      type: data.Type,
      Image: data.Image,
      Description: data.Description,
      Origin: data.Origin,
      OtherDetails: data.OtherDetails,
    };
    try {
      await this.floraUpstreamService.create(createFloraUpstreamDto);
    } catch (error) {
      console.log(error);
      channel.nack(originalMsg);
    } finally {
      channel.ack(originalMsg);
    }
  }

  @MessagePattern({ cmd: 'update_flora' })
  async update(@Payload() data: RabbitMqPayload, @Ctx() context: RmqContext) {
    const channel = context.getChannelRef();
    const originalMsg = context.getMessage();

    const updateFloraUpstreamDto: FloraUpstream = {
      common_name: data.CommonName,
      scientific_name: data.ScientificName,
      user_id: data.UserId,
      type: data.Type,
      Image: data.Image,
      Description: data.Description,
      Origin: data.Origin,
      OtherDetails: data.OtherDetails,
    };
    try {
      await this.floraUpstreamService.update(
        data.id as string,
        updateFloraUpstreamDto,
      );
    } catch (error) {
      console.log(error);
      channel.nack(originalMsg);
    } finally {
      channel.ack(originalMsg);
    }
  }
}
