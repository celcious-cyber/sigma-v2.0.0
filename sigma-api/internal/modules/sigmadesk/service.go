package sigmadesk

import (
	"gorm.io/gorm"
	"sigma-api/internal/core/models"
)

type DeskService interface {
	GetAnnouncements(activeOnly bool) ([]models.Announcement, error)
	GetVisitors() ([]models.Visitor, error)
	GetIncomingLetters() ([]models.IncomingLetter, error)
	GetOutgoingLetters() ([]models.OutgoingLetter, error)
	CreateAnnouncement(a models.Announcement) error
}

type deskService struct {
	db *gorm.DB
}

func NewDeskService(db *gorm.DB) DeskService {
	return &deskService{db}
}

func (s *deskService) GetAnnouncements(activeOnly bool) ([]models.Announcement, error) {
	var announcements []models.Announcement
	query := s.db.Order("created_at desc")
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}
	err := query.Find(&announcements).Error
	return announcements, err
}

func (s *deskService) GetVisitors() ([]models.Visitor, error) {
	var visitors []models.Visitor
	err := s.db.Order("created_at desc").Find(&visitors).Error
	return visitors, err
}

func (s *deskService) GetIncomingLetters() ([]models.IncomingLetter, error) {
	var letters []models.IncomingLetter
	err := s.db.Order("date desc").Find(&letters).Error
	return letters, err
}

func (s *deskService) GetOutgoingLetters() ([]models.OutgoingLetter, error) {
	var letters []models.OutgoingLetter
	err := s.db.Order("date desc").Find(&letters).Error
	return letters, err
}

func (s *deskService) CreateAnnouncement(a models.Announcement) error {
	return s.db.Create(&a).Error
}
