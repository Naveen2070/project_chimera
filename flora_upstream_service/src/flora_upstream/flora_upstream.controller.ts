//	Copyright 2025 Naveen R
//
//		Licensed under the Apache License, Version 2.0 (the "License");
//		you may not use this file except in compliance with the License.
//		You may obtain a copy of the License at
//
//		http://www.apache.org/licenses/LICENSE-2.0
//
//		Unless required by applicable law or agreed to in writing, software
//		distributed under the License is distributed on an "AS IS" BASIS,
//		WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//		See the License for the specific language governing permissions and
//		limitations under the License.
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
        data.ID as string,
        updateFloraUpstreamDto,
      );
    } catch (error) {
      console.log(error);
    } finally {
      channel.ack(originalMsg);
    }
  }
}
