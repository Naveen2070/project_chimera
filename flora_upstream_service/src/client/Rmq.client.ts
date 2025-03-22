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
import { Module, DynamicModule } from '@nestjs/common';
import { ClientsModule, Transport } from '@nestjs/microservices';

@Module({})
export class RmqClientModule {
  static register(name: string, queue: string): DynamicModule {
    return {
      module: RmqClientModule,
      imports: [
        ClientsModule.register([
          {
            name: `${name.toUpperCase()}_SERVICE`,
            transport: Transport.RMQ,
            options: {
              urls: ['amqp://admin:naveen@2007@localhost:5672'],
              queue,
            },
          },
        ]),
      ],
      exports: [ClientsModule],
    };
  }
}
