package enrollment

import (
	"log"

	"github.com/JuD4Mo/go_api_web_domain/domain"
)

type (
	Service interface {
		Create(userId, courseId string) (*domain.Enrollment, error)
	}

	service struct {
		log  *log.Logger
		repo Repository
	}
)

func NewService(log *log.Logger, repo Repository) Service {
	return &service{
		log:  log,
		repo: repo,
	}
}

func (s service) Create(userId, courseId string) (*domain.Enrollment, error) {
	enroll := &domain.Enrollment{
		UserId:   userId,
		CourseId: courseId,
		Status:   "P",
	}

	err := s.repo.Create(enroll)
	if err != nil {
		return nil, err
	}

	return enroll, nil
}
