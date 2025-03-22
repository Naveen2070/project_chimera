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
    try {
      const pgData: FloraPg = toFloraPg(data);
      const pgResult = await this.prisma.flora.create({
        data: pgData,
      });

      const mongoData: FloraMongo = toFloraMongo(data);
      const mongoResult = new this.floraModel(mongoData);
      await mongoResult.save();

      const result: FloraUpstream = toFloraUpstream(pgResult, mongoData);

      return result;
    } catch (error) {
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
      const pgResult = await this.prisma.flora.update({
        where: { id },
        data: pgData,
      });

      const mongoData: FloraMongo = toFloraMongo(data);

      const mongoResult = await this.floraModel.findByIdAndUpdate(
        data.id,
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
