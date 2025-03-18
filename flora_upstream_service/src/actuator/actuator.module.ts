import { Module } from '@nestjs/common';
import { ActuatorService } from './actuator.service';
import { ActuatorController } from './actuator.controller';
import { RmqClientModule } from 'src/client/Rmq.client';

@Module({
  imports: [RmqClientModule.register('ACTUATOR', 'flora_upstream_queue')],
  controllers: [ActuatorController],
  providers: [ActuatorService],
})
export class ActuatorModule {}
