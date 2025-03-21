-- CreateEnum
CREATE TYPE "PostType" AS ENUM ('public', 'private');

-- CreateTable
CREATE TABLE "Flora" (
    "id" UUID NOT NULL,
    "common_name" TEXT NOT NULL,
    "scientific_name" TEXT NOT NULL,
    "user_id" TEXT NOT NULL,
    "type" "PostType" NOT NULL,

    CONSTRAINT "Flora_pkey" PRIMARY KEY ("id")
);
