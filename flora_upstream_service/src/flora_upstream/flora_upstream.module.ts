import { Module } from '@nestjs/common';
import { FloraUpstreamService } from './flora_upstream.service';
import { FloraUpstreamController } from './flora_upstream.controller';

@Module({
  controllers: [FloraUpstreamController],
  providers: [FloraUpstreamService],
})
export class FloraUpstreamModule {}
