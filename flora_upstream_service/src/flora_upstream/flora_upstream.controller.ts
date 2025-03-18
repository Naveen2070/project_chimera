import { Controller } from '@nestjs/common';
import { MessagePattern, Payload } from '@nestjs/microservices';
import { FloraUpstreamService } from './flora_upstream.service';
import { CreateFloraUpstreamDto } from './dto/create-flora_upstream.dto';
import { UpdateFloraUpstreamDto } from './dto/update-flora_upstream.dto';

@Controller()
export class FloraUpstreamController {
  constructor(private readonly floraUpstreamService: FloraUpstreamService) {}

  @MessagePattern('createFloraUpstream')
  create(@Payload() createFloraUpstreamDto: CreateFloraUpstreamDto) {
    return this.floraUpstreamService.create(createFloraUpstreamDto);
  }

  @MessagePattern('findAllFloraUpstream')
  findAll() {
    return this.floraUpstreamService.findAll();
  }

  @MessagePattern('findOneFloraUpstream')
  findOne(@Payload() id: number) {
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
  remove(@Payload() id: number) {
    return this.floraUpstreamService.remove(id);
  }
}
