import {
  Injectable,
  OnModuleInit,
  BeforeApplicationShutdown,
} from '@nestjs/common';
import Consul from 'consul';

@Injectable()
export class ConsulService implements OnModuleInit, BeforeApplicationShutdown {
  private consul: Consul;
  private readonly serviceId = 'flora-upstream-service';

  constructor() {
    this.consul = new Consul({
      host: 'localhost',
      port: 8500,
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
        address: 'localhost',
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
