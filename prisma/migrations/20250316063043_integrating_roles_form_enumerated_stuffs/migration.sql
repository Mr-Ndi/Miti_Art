/*
  Warnings:

  - Changed the type of `tax_pin` on the `Vendor` table. No cast exists, the column would be dropped and recreated, which cannot be done if there is data, since the column is required.

*/
-- AlterTable
ALTER TABLE "Vendor" DROP COLUMN "tax_pin",
ADD COLUMN     "tax_pin" INTEGER NOT NULL;

-- CreateIndex
CREATE UNIQUE INDEX "Vendor_tax_pin_key" ON "Vendor"("tax_pin");
