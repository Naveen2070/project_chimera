import { Inject, Injectable } from '@nestjs/common';
import { ClientProxy } from '@nestjs/microservices';

@Injectable()
export class ActuatorService {
  constructor(
    @Inject('ACTUATOR_SERVICE') private readonly RmqClient: ClientProxy,
  ) {}

  async isHealthy() {
    try {
      await this.RmqClient.connect();

      return {
        status: 'up',
        message: 'RabbitMQ connection is healthy',
      };
    } catch (err) {
      console.error(err);
      return {
        status: 'down',
        message: 'RabbitMQ connection is down',
      };
    }
  }
}
