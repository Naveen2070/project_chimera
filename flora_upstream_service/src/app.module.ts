import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { ConsulService } from './consul/consul.service';
import { FloraUpstreamModule } from './flora_upstream/flora_upstream.module';
import { ActuatorModule } from './actuator/actuator.module';

@Module({
  imports: [FloraUpstreamModule, ActuatorModule],
  controllers: [AppController],
  providers: [AppService, ConsulService],
})
export class AppModule {}
