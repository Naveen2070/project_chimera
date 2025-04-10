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
import {
  Injectable,
  OnModuleInit,
  BeforeApplicationShutdown,
} from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import Consul from 'consul';

@Injectable()
export class ConsulService implements OnModuleInit, BeforeApplicationShutdown {
  private consul: Consul;
  private readonly serviceId = 'flora-upstream-service';
  private readonly consul_host: string;
  private readonly consul_port: number;

  constructor(private readonly configService: ConfigService) {
    this.consul_host =
      this.configService.get<string>('CONSUL_HOST') || 'localhost';
    this.consul_port = this.configService.get<number>('CONSUL_PORT') || 8500;
    this.consul = new Consul({
      host: this.consul_host,
      port: this.consul_port,
      secure: false,
    });
  }

  async onModuleInit() {
    await this.registerService();
  }

  async beforeApplicationShutdown() {
    // This will be called when NestJS is preparing to shut down
    console.log('Shutting down...');
    await this.deregisterService();
  }

  private async registerService() {
    try {
      await this.consul.agent.service.register({
        name: this.serviceId,
        address: 'host.docker.internal',
        port: 3030,
        tags: ['rabbitmq', 'flora', 'upstream'],
        meta: {
          protocol: 'rabbitmq',
          description: 'RabbitMQ producer for flora data',
        },
        check: {
          status: 'passing',
          name: 'Flora upstream service health check',
          http: 'http://host.docker.internal:3030/actuator/health',
          interval: '15s',
          timeout: '10s',
          deregistercriticalserviceafter: '5m',
        },
      });

      console.log(`Service "${this.serviceId}" registered in Consul.`);
    } catch (err) {
      if (err instanceof Error) {
        console.error('Failed to register service in Consul:', err.message);
      } else {
        console.error(
          'Failed to register service in Consul with unknown error:',
          err,
        );
      }
    }
  }

  private async deregisterService() {
    try {
      await this.consul.agent.service.deregister(this.serviceId);
      console.log(`Service "${this.serviceId}" deregistered from Consul.`);
    } catch (err) {
      if (err instanceof Error) {
        console.error('Failed to deregister service from Consul:', err.message);
      } else {
        console.error(
          'Failed to deregister service from Consul with unknown error:',
          err,
        );
      }
    }
  }
}
