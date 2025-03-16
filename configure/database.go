package configure

import (
	"MITI_ART/prisma/miti_art"
	"log"
)

func InitDB() *miti_art.PrismaClient {
	prisma := miti_art.NewClient()

	if err := prisma.Connect(); err != nil {
		log.Fatalf("\n\n❌ Failed to connect to Prisma: %v", err)
	}
	log.Println("\n\n✅ Connected to Prisma successfully")
	return prisma
}
