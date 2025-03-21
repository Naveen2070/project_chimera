import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { ConsulService } from './consul/consul.service';
import { FloraUpstreamModule } from './flora_upstream/flora_upstream.module';
import { ActuatorModule } from './actuator/actuator.module';
import { PrismaService } from './prisma_client/prisma.service';

@Module({
  imports: [FloraUpstreamModule, ActuatorModule],
  controllers: [AppController],
  providers: [AppService, ConsulService, PrismaService],
})
export class AppModule {}
