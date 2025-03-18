import { Controller, Get } from '@nestjs/common';
import { ActuatorService } from './actuator.service';

@Controller('actuator')
export class ActuatorController {
  constructor(private readonly actuatorService: ActuatorService) {}

  @Get('health')
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
