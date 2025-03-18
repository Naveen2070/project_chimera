import { Module, DynamicModule } from '@nestjs/common';
import { ClientsModule, Transport } from '@nestjs/microservices';

@Module({})
export class RmqClientModule {
  static register(name: string, queue: string): DynamicModule {
    return {
      module: RmqClientModule,
      imports: [
        ClientsModule.register([
          {
            name: `${name.toUpperCase()}_SERVICE`,
            transport: Transport.RMQ,
            options: {
              urls: ['amqp://admin:naveen@2007@localhost:5672'],
              queue,
            },
          },
        ]),
      ],
      exports: [ClientsModule],
    };
  }
}
