package enrollment

import (
	"context"
	"log"

	"github.com/JuD4Mo/go_api_web_domain/domain"
	courseSDK "github.com/JuD4Mo/go_api_web_sdk/course"
	userSDK "github.com/JuD4Mo/go_api_web_sdk/user"
)

type (
	Service interface {
		Create(ctx context.Context, userId, courseId string) (*domain.Enrollment, error)
		GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Enrollment, error)
		Update(ctx context.Context, id string, status *string) error
		Count(ctx context.Context, filters Filters) (int, error)
	}

	service struct {
		log             *log.Logger
		repo            Repository
		userTransport   userSDK.Transport
		courseTransport courseSDK.Transport
	}

	Filters struct {
		UserId   string
		CourseId string
	}
)

func NewService(log *log.Logger, repo Repository, userTransport userSDK.Transport, courseTransport courseSDK.Transport) Service {
	return &service{
		log:             log,
		repo:            repo,
		userTransport:   userTransport,
		courseTransport: courseTransport,
	}
}

func (s service) Create(ctx context.Context, userId, courseId string) (*domain.Enrollment, error) {
	enroll := &domain.Enrollment{
		UserID:   userId,
		CourseID: courseId,
		Status:   domain.Pending,
	}

	_, err := s.userTransport.Get(userId)
	if err != nil {
		return nil, err
	}

	_, err = s.courseTransport.Get(courseId)
	if err != nil {
		s.log.Println(err)
		return nil, err
	}

	err = s.repo.Create(ctx, enroll)
	if err != nil {
		return nil, err
	}

	return enroll, nil
}

func (s service) GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Enrollment, error) {
	enrollments, err := s.repo.GetAll(ctx, filters, offset, limit)
	if err != nil {
		return nil, err
	}
	return enrollments, nil
}

func (s service) Update(ctx context.Context, id string, status *string) error {

	if status != nil {
		switch domain.EnrollStatus(*status) {
		case domain.Pending, domain.Active, domain.Studying:
		default:
			return ErrInvalidStatus{*status}
		}
	}

	if err := s.repo.Update(ctx, id, status); err != nil {
		return err
	}
	return nil
}

func (s service) Count(ctx context.Context, filters Filters) (int, error) {
	return s.repo.Count(ctx, filters)
}
