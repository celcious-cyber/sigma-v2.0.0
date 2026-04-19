package database

import (
	"log"
	"sigma-api/internal/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("sigma.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto Migrate the schema
	err = DB.AutoMigrate(&models.Santri{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database connection established and migrated")
}
