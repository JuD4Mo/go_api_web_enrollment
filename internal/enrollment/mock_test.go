package enrollment_test

import (
	"context"

	"github.com/JuD4Mo/go_api_web_domain/domain"
	"github.com/JuD4Mo/go_api_web_enrollment/internal/enrollment"
)

type mockRepository struct {
	CreateMock func(ctx context.Context, enroll *domain.Enrollment) error
	GetAllMock func(ctx context.Context, filters enrollment.Filters, offset, limit int) ([]domain.Enrollment, error)
	UpdateMock func(ctx context.Context, id string, status *string) error
	CountMock  func(ctx context.Context, filter enrollment.Filters) (int, error)
}

func (mock *mockRepository) Create(ctx context.Context, enroll *domain.Enrollment) error {
	return mock.CreateMock(ctx, enroll)
}

func (mock *mockRepository) GetAll(ctx context.Context, filters enrollment.Filters, offset, limit int) ([]domain.Enrollment, error) {
	return mock.GetAllMock(ctx, filters, offset, limit)
}

func (mock *mockRepository) Update(ctx context.Context, id string, status *string) error {
	return mock.UpdateMock(ctx, id, status)
}

func (mock *mockRepository) Count(ctx context.Context, filter enrollment.Filters) (int, error) {
	return mock.CountMock(ctx, filter)
}
