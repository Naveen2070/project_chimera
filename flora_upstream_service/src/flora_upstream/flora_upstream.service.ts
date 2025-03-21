import { Injectable } from '@nestjs/common';
import { CreateFloraUpstreamDto } from './dto/create-flora_upstream.dto';
import { UpdateFloraUpstreamDto } from './dto/update-flora_upstream.dto';
import { PrismaService } from 'src/prisma_client/prisma.service';
import { Prisma } from '@prisma/client';
import { InjectModel } from '@nestjs/mongoose';
import { Flora as FloraMongo } from './schema/flora.schema';
import { Model } from 'mongoose';
import { FloraUpstream } from './entities/flora_upstream.entity';

@Injectable()
export class FloraUpstreamService {
  constructor(
    private readonly prisma: PrismaService,
    @InjectModel(FloraMongo.name) private catModel: Model<FloraMongo>,
  ) {}
  async create(data: FloraUpstream): Promise<FloraUpstream> {
    try {
      const pgData: CreateFloraUpstreamDto = {
        common_name: data.common_name,
        scientific_name: data.scientific_name,
        user_id: data.user_id,
        type: data.type,
      };
      const pgResult = await this.prisma.flora.create({
        data: pgData,
      });

      const mongoData: FloraMongo = {
        flora_id: pgData.id as string,
        Image: Buffer.from(data.Image),
        Description: data.Description,
        Origin: data.Origin,
        OtherDetails: data.OtherDetails,
      };
      const mongoResult = new this.catModel(mongoData);
      await mongoResult.save();

      const result: FloraUpstream = {
        id: pgResult.id,
        common_name: pgResult.common_name,
        scientific_name: pgResult.scientific_name,
        user_id: pgResult.user_id,
        type: pgResult.type,
        Image: mongoData.Image,
        Description: mongoData.Description,
        Origin: mongoData.Origin,
        OtherDetails: mongoData.OtherDetails,
      };

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

  async update(id: string, updateFloraUpstreamDto: UpdateFloraUpstreamDto) {
    return await this.prisma.flora.update({
      where: { id },
      data: updateFloraUpstreamDto,
    });
  }
}
