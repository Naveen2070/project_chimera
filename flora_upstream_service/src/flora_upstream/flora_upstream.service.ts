import { Injectable } from '@nestjs/common';
import { CreateFloraUpstreamDto } from './dto/create-flora_upstream.dto';
import { UpdateFloraUpstreamDto } from './dto/update-flora_upstream.dto';

@Injectable()
export class FloraUpstreamService {
  create(createFloraUpstreamDto: any) {
    return console.log(createFloraUpstreamDto);
  }

  findAll() {
    return `This action returns all floraUpstream`;
  }

  findOne(id: number) {
    return `This action returns a #${id} floraUpstream`;
  }

  update(id: number, updateFloraUpstreamDto: UpdateFloraUpstreamDto) {
    return `This action updates a #${id} floraUpstream`;
  }

  remove(id: number) {
    return `This action removes a #${id} floraUpstream`;
  }
}
