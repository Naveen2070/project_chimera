import { Injectable } from '@nestjs/common';
import { CreateFloraUpstreamDto } from './dto/create-flora_upstream.dto';
import { UpdateFloraUpstreamDto } from './dto/update-flora_upstream.dto';
import { PrismaService } from 'src/prisma_client/prisma.service';
import { Prisma } from '@prisma/client';

@Injectable()
export class FloraUpstreamService {
  constructor(private prisma: PrismaService) {}
  async create(createFloraUpstreamDto: CreateFloraUpstreamDto) {
    return await this.prisma.flora.create({
      data: createFloraUpstreamDto,
    });
  }

  async findAll(params?: {
    skip?: number;
    take?: number;
    cursor?: Prisma.FloraWhereUniqueInput;
    where?: Prisma.FloraWhereUniqueInput;
    orderBy?: Prisma.FloraOrderByWithRelationInput;
  }) {
    return await this.prisma.flora.findMany({
      ...params,
    });
  }

  async findOne(id: string) {
    return await this.prisma.flora.findUnique({
      where: { id },
    });
  }

  async update(id: string, updateFloraUpstreamDto: UpdateFloraUpstreamDto) {
    return await this.prisma.flora.update({
      where: { id },
      data: updateFloraUpstreamDto,
    });
  }

  async remove(id: string) {
    return await this.prisma.flora.delete({
      where: { id },
    });
  }
}
