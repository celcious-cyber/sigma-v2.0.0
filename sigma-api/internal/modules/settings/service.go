package settings

import (
	"gorm.io/gorm"
	"sigma-api/internal/core/models"
)

type SettingsService interface {
	GetSettings() (models.Setting, error)
	UpdateSettings(s models.Setting) error
	GetMe(userID uint) (models.User, error)
}

type settingsService struct {
	db *gorm.DB
}

func NewSettingsService(db *gorm.DB) SettingsService {
	return &settingsService{db}
}

func (s *settingsService) GetSettings() (models.Setting, error) {
	var setting models.Setting
	err := s.db.First(&setting).Error
	return setting, err
}

func (s *settingsService) UpdateSettings(st models.Setting) error {
	return s.db.Save(&st).Error
}

func (s *settingsService) GetMe(userID uint) (models.User, error) {
	var user models.User
	err := s.db.First(&user, userID).Error
	return user, err
}
