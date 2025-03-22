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
import { Controller, Get, HttpCode, HttpStatus, Header } from '@nestjs/common';
import { ActuatorService } from './actuator.service';

@Controller('actuator')
export class ActuatorController {
  constructor(private readonly actuatorService: ActuatorService) {}

  @Get('health')
  @HttpCode(HttpStatus.OK)
  @Header(
    'Cache-Control',
    'no-store, no-cache, must-revalidate, proxy-revalidate',
  )
  @Header('Pragma', 'no-cache')
  @Header('Content-Type', 'application/json')
  async getHealth() {
    const RmqStatus = await this.actuatorService.isHealthy();
    return {
      status: RmqStatus.status === 'up' ? 'UP' : 'DOWN',
      components: {
        RmqStatus,
      },
    };
  }
}
