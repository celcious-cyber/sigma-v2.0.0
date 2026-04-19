package sigmalit

import (
	"gorm.io/gorm"
	"sigma-api/internal/core/models"
)

type LibraryService interface {
	GetBooks() ([]models.Book, error)
}

type libraryService struct {
	db *gorm.DB
}

func NewLibraryService(db *gorm.DB) LibraryService {
	return &libraryService{db}
}

func (s *libraryService) GetBooks() ([]models.Book, error) {
	var books []models.Book
	err := s.db.Find(&books).Error
	return books, err
}
