import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import { Transport, MicroserviceOptions } from '@nestjs/microservices';

async function bootstrap() {
  const app = await NestFactory.create(AppModule);

  // Connect RabbitMQ microservice (consumer)
  app.connectMicroservice<MicroserviceOptions>({
    transport: Transport.RMQ,
    options: {
      urls: ['amqp://admin:naveen@2007@localhost:5672'],
      queue: 'flora_upstream_queue',
      queueOptions: { durable: true, exclusive: false },
    },
  });

  // Enable graceful shutdown
  app.enableShutdownHooks();

  await app.startAllMicroservices();
  await app.listen(process.env.PORT ?? 3030);
}
bootstrap()
  .then(() => {
    console.log('Server started on port', process.env.PORT ?? 3030);
  })
  .catch((err) => console.error(err));
