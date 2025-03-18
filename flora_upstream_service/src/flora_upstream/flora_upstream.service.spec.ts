import { Test, TestingModule } from '@nestjs/testing';
import { FloraUpstreamService } from './flora_upstream.service';

describe('FloraUpstreamService', () => {
  let service: FloraUpstreamService;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [FloraUpstreamService],
    }).compile();

    service = module.get<FloraUpstreamService>(FloraUpstreamService);
  });

  it('should be defined', () => {
    expect(service).toBeDefined();
  });
});
