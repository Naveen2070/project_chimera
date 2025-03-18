import { Test, TestingModule } from '@nestjs/testing';
import { ActuatorController } from './actuator.controller';
import { ActuatorService } from './actuator.service';

describe('ActuatorController', () => {
  let controller: ActuatorController;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      controllers: [ActuatorController],
      providers: [ActuatorService],
    }).compile();

    controller = module.get<ActuatorController>(ActuatorController);
  });

  it('should be defined', () => {
    expect(controller).toBeDefined();
  });
});
