/*
  Warnings:

  - A unique constraint covering the columns `[scientific_name]` on the table `Flora` will be added. If there are existing duplicate values, this will fail.

*/
-- CreateIndex
CREATE UNIQUE INDEX "Flora_scientific_name_key" ON "Flora"("scientific_name");
