import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import { Transport, MicroserviceOptions } from '@nestjs/microservices';
import { DocumentBuilder, SwaggerModule } from '@nestjs/swagger';

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

  app.enableCors({
    origin: '*',
    methods: '*',
    allowedHeaders: '*',
    credentials: true,
  });

  const config = new DocumentBuilder()
    .setTitle('Chimera Flora Upstream Service')
    .setDescription(
      'This API provides data uploads and updates (or) upstream for Chimera Flora data',
    )
    .setVersion('1.0.0')
    .setContact(
      'Naveen R',
      'https://naveen2070.github.io/portfolio',
      'naveenrameshcud@gmail.com',
    )
    .setLicense('Apache 2.0', 'http://www.apache.org/licenses/LICENSE-2.0')
    .addServer('http://localhost:3030/', 'Local Server')
    .addServer('http://localhost:8080/flora-upstream/', 'Gateway Server')
    .addBearerAuth()
    .build();

  const documentFactory = () => SwaggerModule.createDocument(app, config);
  SwaggerModule.setup('api', app, documentFactory, {
    jsonDocumentUrl: 'swagger/v1/swagger.json',
  });

  await app.startAllMicroservices();
  await app.listen(process.env.PORT ?? 3030, '0.0.0.0');
}
bootstrap()
  .then(() => {
    console.log('Server started on port', process.env.PORT ?? 3030);
  })
  .catch((err) => console.error(err));
