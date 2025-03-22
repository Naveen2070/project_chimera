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
import { Inject, Injectable } from '@nestjs/common';
import { ClientProxy } from '@nestjs/microservices';

@Injectable()
export class ActuatorService {
  constructor(
    @Inject('ACTUATOR_SERVICE') private readonly RmqClient: ClientProxy,
  ) {}

  async isHealthy() {
    try {
      await this.RmqClient.connect();

      return {
        status: 'up',
        message: 'RabbitMQ connection is healthy',
      };
    } catch (err) {
      console.error(err);
      return {
        status: 'down',
        message: 'RabbitMQ connection is down',
      };
    }
  }
}
