package database

import (
	"fmt"
	"log"
	"sigma-api/internal/core/models"
)

func SeedMedicalData() {
	// 1. Create Facility
	facility := models.MedicalFacility{
		Name: "UKS Utama",
		Type: "WARD",
	}
	
	// Check if already exists
	var count int64
	DB.Model(&models.MedicalFacility{}).Where("name = ?", "UKS Utama").Count(&count)
	if count > 0 {
		return
	}

	if err := DB.Create(&facility).Error; err != nil {
		log.Printf("Failed to seed facility: %v", err)
		return
	}

	// 2. Create 12 Beds
	for i := 1; i <= 12; i++ {
		bed := models.Bed{
			FacilityID: facility.ID,
			Number:     fmt.Sprintf("BED-%02d", i),
			IsOccupied: false,
		}
		DB.Create(&bed)
	}
	log.Println("Medical beds seeded successfully!")
}
