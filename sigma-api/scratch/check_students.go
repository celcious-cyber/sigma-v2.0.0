package main

import (
	"fmt"
	"log"
	"sigma-api/internal/core/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("../data/sigma.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	var count int64
	db.Model(&models.Student{}).Count(&count)
	fmt.Printf("Total Students: %d\n", count)
    
    var students []models.Student
    db.Find(&students)
    for _, s := range students {
        fmt.Printf("- %s (ID: %d, Active: %v)\n", s.Name, s.ID, s.IsActive)
    }
}
