import { Module } from '@nestjs/common';
import { FloraUpstreamService } from './flora_upstream.service';
import { FloraUpstreamController } from './flora_upstream.controller';
import { PrismaService } from 'src/prisma_client/prisma.service';
import { MongooseModule } from '@nestjs/mongoose';
import { Flora, FloraSchema } from './schema/flora.schema';

@Module({
  imports: [
    MongooseModule.forFeature([{ name: Flora.name, schema: FloraSchema }]),
  ],
  controllers: [FloraUpstreamController],
  providers: [FloraUpstreamService, PrismaService],
})
export class FloraUpstreamModule {}
