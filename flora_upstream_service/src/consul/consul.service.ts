import { Injectable, OnModuleInit, OnModuleDestroy } from '@nestjs/common';
import Consul from 'consul';

@Injectable()
export class ConsulService implements OnModuleInit, OnModuleDestroy {
  private consul: Consul;
  private readonly serviceId = 'chimera-flora-upstream';

  constructor() {
    this.consul = new Consul({
      host: 'localhost',
      port: 8500,
    });
  }

  async onModuleInit() {
    await this.registerService();
  }

  async onModuleDestroy() {
    await this.deregisterService();
  }

  private async registerService() {
    try {
      await this.consul.agent.service.register({
        id: this.serviceId,
        name: 'chimera-flora-upstream',
        address: '10.0.0.5',
        port: 3000,
        tags: ['rabbitmq', 'flora', 'upstream'],
        meta: {
          protocol: 'rabbitmq',
          description: 'RabbitMQ producer for flora data',
        },
        check: {
          http: `http://10.0.0.5:3000/actuator/health`,
          interval: '10s',
          timeout: '5s',
          deregistercriticalserviceafter: '1m',
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
