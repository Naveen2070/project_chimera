import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { ConsulService } from './consul/consul.service';
import { FloraUpstreamModule } from './flora_upstream/flora_upstream.module';
import { ActuatorModule } from './actuator/actuator.module';
import { PrismaService } from './prisma_client/prisma.service';
import { ConfigModule, ConfigService } from '@nestjs/config';
import { MongooseModule } from '@nestjs/mongoose';

@Module({
  imports: [
    FloraUpstreamModule,
    ActuatorModule,
    ConfigModule.forRoot({ isGlobal: true }),
    MongooseModule.forRootAsync({
      imports: [ConfigModule],
      useFactory: (configService: ConfigService) => ({
        uri: configService.get<string>('DATABASE_URL_MONGO'),
      }),
      inject: [ConfigService],
    }),
  ],
  controllers: [AppController],
  providers: [AppService, ConsulService, PrismaService],
})
export class AppModule {}
