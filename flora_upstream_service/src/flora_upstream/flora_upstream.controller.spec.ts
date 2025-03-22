//	Copyright 2025 Naveen R
//
//		Licensed under the Apache License, Version 2.0 (the "License");
//		you may not use this file except in compliance with the License.
//		You may obtain a copy of the License at
//
//		http://www.apache.org/licenses/LICENSE-2.0
//
//		Unless required by applicable law or agreed to in writing, software
//		distributed under the License is distributed on an "AS IS" BASIS,
//		WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//		See the License for the specific language governing permissions and
//		limitations under the License.
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
