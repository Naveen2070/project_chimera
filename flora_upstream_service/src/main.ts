import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import { Transport, MicroserviceOptions } from '@nestjs/microservices';
import { DocumentBuilder, SwaggerModule } from '@nestjs/swagger';
import { ConfigService } from '@nestjs/config';

async function bootstrap() {
  const app = await NestFactory.create(AppModule);

  const configService = app.get(ConfigService);

  const port = configService.get<number>('PORT') || 3030;
  const RABBIT_MQ_URL = configService.get<string>('RABBIT_MQ_URL') as string;
  const RABBIT_QUEUE_FLORA =
    configService.get<string>('RABBIT_QUEUE_FLORA') || 'flora_upstream_queue';

  // Connect RabbitMQ microservice (consumer)
  app.connectMicroservice<MicroserviceOptions>({
    transport: Transport.RMQ,
    options: {
      urls: [RABBIT_MQ_URL],
      queue: RABBIT_QUEUE_FLORA,
      queueOptions: { durable: true, exclusive: false },
      noAck: false,
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
  await app.listen(port, '0.0.0.0').then(() => {
    console.log(`Server started on port ${port}`);
  });
}

bootstrap().catch((err) => console.error(err));
