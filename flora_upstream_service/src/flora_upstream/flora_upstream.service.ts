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
import { Injectable } from '@nestjs/common';
import { FloraPg } from './dto/create-flora_upstream.dto';
import { PrismaService } from 'src/prisma_client/prisma.service';
import { Prisma } from '@prisma/client';
import { InjectModel } from '@nestjs/mongoose';
import { Flora as FloraMongo } from './schema/flora.schema';
import { Model } from 'mongoose';
import { FloraUpstream } from './entities/flora_upstream.entity';
import {
  toFloraMongo,
  toFloraPg,
  toFloraUpstream,
} from 'src/utils/data-mapper';

@Injectable()
export class FloraUpstreamService {
  constructor(
    private readonly prisma: PrismaService,
    @InjectModel(FloraMongo.name) private floraModel: Model<FloraMongo>,
  ) {}

  async create(data: FloraUpstream): Promise<FloraUpstream | Error> {
    let id: string | null = null;
    try {
      const pgData: FloraPg = toFloraPg(data);
      const pgResult = await this.prisma.$transaction(async (prisma) => {
        return prisma.flora.create({
          data: pgData,
        });
      });

      id = pgResult.id;

      const mongoData: FloraMongo = toFloraMongo(pgResult.id, data);
      console.log(mongoData);

      const mongoResult = new this.floraModel(mongoData);
      await mongoResult.save();

      const result: FloraUpstream = toFloraUpstream(pgResult, mongoData);

      return result;
    } catch (error) {
      // Rollback PostgreSQL changes if necessary
      if (id) {
        await this.prisma.flora.delete({
          where: { id },
        });
        console.log('PostgreSQL transaction rolled back.');
      }

      console.log(error);
      throw error;
    }
  }

  async findOne(condition: Prisma.FloraWhereUniqueInput) {
    return await this.prisma.flora.findUnique({
      where: condition,
    });
  }

  async update(id: string, data: FloraUpstream) {
    try {
      const pgData: FloraPg = toFloraPg(data);
      const pgResult = await this.prisma.$transaction(async (prisma) => {
        return prisma.flora.update({
          where: { id },
          data: pgData,
        });
      });

      const mongoData: FloraMongo = toFloraMongo(data.id as string, data);

      const mongoResult = await this.floraModel.findOneAndUpdate(
        { flora_id: id },
        { $set: mongoData },
        { new: true },
      );

      if (!mongoResult) {
        console.log('No document found with the given ID.');
        return new Error('No document found with the given ID.');
      }

      const updatedFlora: FloraUpstream = toFloraUpstream(
        pgResult,
        mongoResult as FloraMongo,
      );
      return updatedFlora;
    } catch (error) {
      console.log(error);
      throw error;
    }
  }
}
