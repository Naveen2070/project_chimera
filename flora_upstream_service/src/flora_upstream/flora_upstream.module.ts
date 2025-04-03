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
import { Module } from '@nestjs/common';
import { FloraUpstreamService } from './flora_upstream.service';
import { FloraUpstreamController } from './flora_upstream.controller';
import { PrismaService } from 'src/prisma_client/prisma.service';
import { MongooseModule } from '@nestjs/mongoose';
import { Flora, FloraSchema } from './schema/flora.schema';
import { RmqClientModule } from 'src/client/Rmq.client';

@Module({
  imports: [
    MongooseModule.forFeature([{ name: Flora.name, schema: FloraSchema }]),
    RmqClientModule.register('notification', 'notification_queue'),
    RmqClientModule.register('error', 'error_dump_queue'),
  ],
  controllers: [FloraUpstreamController],
  providers: [FloraUpstreamService, PrismaService],
})
export class FloraUpstreamModule {}
