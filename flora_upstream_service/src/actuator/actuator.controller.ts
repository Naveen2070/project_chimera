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
  async getHealth() {
    const RmqStatus = await this.actuatorService.isHealthy();
    return {
      status: RmqStatus.status === 'up' ? 'UP' : 'DOWN',
    };
  }
}
