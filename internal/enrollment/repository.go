package enrollment

import (
	"log"

	"github.com/JuD4Mo/go_api_web_domain/domain"
	"gorm.io/gorm"
)

type (
	Repository interface {
		Create(enroll *domain.Enrollment) error
	}

	repo struct {
		db  *gorm.DB
		log *log.Logger
	}
)

func NewRepo(db *gorm.DB, log *log.Logger) Repository {
	return &repo{
		db:  db,
		log: log,
	}

}

func (repo *repo) Create(enroll *domain.Enrollment) error {
	result := repo.db.Create(enroll)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
