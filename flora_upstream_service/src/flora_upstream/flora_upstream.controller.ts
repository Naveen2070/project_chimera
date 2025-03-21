import { Controller } from '@nestjs/common';
import {
  Ctx,
  MessagePattern,
  Payload,
  RmqContext,
} from '@nestjs/microservices';
import { FloraUpstreamService } from './flora_upstream.service';
import { CreateFloraUpstreamDto } from './dto/create-flora_upstream.dto';
import { UpdateFloraUpstreamDto } from './dto/update-flora_upstream.dto';

@Controller()
export class FloraUpstreamController {
  constructor(private readonly floraUpstreamService: FloraUpstreamService) {}

  @MessagePattern({ cmd: 'add_flora' })
  create(@Payload() data: any, @Ctx() context: RmqContext) {
    // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
    const channel = context.getChannelRef();
    const originalMsg = context.getMessage();
    // eslint-disable-next-line @typescript-eslint/no-unsafe-call, @typescript-eslint/no-unsafe-member-access
    channel.ack(originalMsg);

    const createFloraUpstreamDto: CreateFloraUpstreamDto = {
      common_name: data.common_name,
      scientific_name: data.scientific_name,
      user_id: data.user_id,
      type: data.type,
    };

    return this.floraUpstreamService.create(createFloraUpstreamDto);
  }

  @MessagePattern('findAllFloraUpstream')
  findAll() {
    return this.floraUpstreamService.findAll();
  }

  @MessagePattern('findOneFloraUpstream')
  findOne(@Payload() id: string) {
    return this.floraUpstreamService.findOne(id);
  }

  @MessagePattern('updateFloraUpstream')
  update(@Payload() updateFloraUpstreamDto: UpdateFloraUpstreamDto) {
    return this.floraUpstreamService.update(
      updateFloraUpstreamDto.id,
      updateFloraUpstreamDto,
    );
  }

  @MessagePattern('removeFloraUpstream')
  remove(@Payload() id: string) {
    return this.floraUpstreamService.remove(id);
  }
}
