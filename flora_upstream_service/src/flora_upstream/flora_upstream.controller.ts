import { Controller, Headers } from '@nestjs/common';
import {
  Ctx,
  MessagePattern,
  Payload,
  RmqContext,
} from '@nestjs/microservices';
import { FloraUpstreamService } from './flora_upstream.service';
import { UpdateFloraUpstreamDto } from './dto/update-flora_upstream.dto';
import { RabbitMqPayload } from './dto/rabbit-payload';
import { FloraUpstream } from './entities/flora_upstream.entity';

@Controller()
export class FloraUpstreamController {
  constructor(private readonly floraUpstreamService: FloraUpstreamService) {}

  @MessagePattern({ cmd: 'add_flora' })
  create(
    @Headers('X-Auth-User') headers,
    @Payload() data: RabbitMqPayload,
    @Ctx() context: RmqContext,
  ) {
    const channel = context.getChannelRef();
    const originalMsg = context.getMessage();

    const createFloraUpstreamDto: FloraUpstream = {
      common_name: data.CommonName,
      scientific_name: data.ScientificName,
      user_id: 'aaaaaa',
      type: data.Type,
      Image: data.Image,
      Description: data.Description,
      Origin: data.Origin,
      OtherDetails: data.OtherDetails,
    };
    channel.ack(originalMsg);

    return this.floraUpstreamService.create(createFloraUpstreamDto);
  }

  @MessagePattern('updateFloraUpstream')
  update(@Payload() updateFloraUpstreamDto: UpdateFloraUpstreamDto) {
    return this.floraUpstreamService.update(
      updateFloraUpstreamDto.id,
      updateFloraUpstreamDto,
    );
  }
}
