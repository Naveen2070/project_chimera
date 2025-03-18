import { Test, TestingModule } from '@nestjs/testing';
import { FloraUpstreamController } from './flora_upstream.controller';
import { FloraUpstreamService } from './flora_upstream.service';

describe('FloraUpstreamController', () => {
  let controller: FloraUpstreamController;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      controllers: [FloraUpstreamController],
      providers: [FloraUpstreamService],
    }).compile();

    controller = module.get<FloraUpstreamController>(FloraUpstreamController);
  });

  it('should be defined', () => {
    expect(controller).toBeDefined();
  });
});
