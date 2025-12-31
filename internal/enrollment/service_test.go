package enrollment_test

import (
	"context"
	"errors"
	"io"
	"log"
	"testing"

	"github.com/JuD4Mo/go_api_web_domain/domain"
	"github.com/JuD4Mo/go_api_web_enrollment/internal/enrollment"
	"github.com/stretchr/testify/assert"
)

func TestServiceGetAll(t *testing.T) {

	//Grupo de pruebas dentro del test
	l := log.New(io.Discard, "", 0)
	count := 0
	expectedCounter := 1
	t.Run("should return an error", func(t *testing.T) {
		expectedErr := errors.New("some error")
		repo := &mockRepository{
			GetAllMock: func(ctx context.Context, filters enrollment.Filters, offset, limit int) ([]domain.Enrollment, error) {
				count++
				return nil, errors.New("some error")
			},
		}

		service := enrollment.NewService(l, repo, nil, nil)

		enrollments, err := service.GetAll(context.Background(), enrollment.Filters{}, 0, 10)
		assert.Error(t, err)
		assert.EqualError(t, expectedErr, err.Error())
		assert.Equal(t, expectedCounter, count)
		assert.Nil(t, enrollments)
	})

	t.Run("get all enrollments", func(t *testing.T) {
		expectedData := []domain.Enrollment{
			{
				ID:       "1",
				UserID:   "11",
				CourseID: "22",
				Status:   "P",
			},
		}

		repo := &mockRepository{
			GetAllMock: func(ctx context.Context, filters enrollment.Filters, offset, limit int) ([]domain.Enrollment, error) {
				return []domain.Enrollment{
					{
						ID:       "1",
						UserID:   "11",
						CourseID: "22",
						Status:   "P",
					},
				}, nil
			},
		}

		service := enrollment.NewService(l, repo, nil, nil)

		enrollments, err := service.GetAll(context.Background(), enrollment.Filters{}, 0, 10)

		assert.Nil(t, err)
		assert.NotNil(t, enrollments)
		assert.Equal(t, expectedData, enrollments)
	})
}

func TestService_Update(t *testing.T) {
	l := log.New(io.Discard, "", 0)

	t.Run("should return an error", func(t *testing.T) {
		expectedErr := errors.New("some error")
		repo := &mockRepository{
			UpdateMock: func(ctx context.Context, id string, status *string) error {
				return errors.New("some error")
			},
		}

		service := enrollment.NewService(l, repo, nil, nil)

		status := "A"
		err := service.Update(context.Background(), "11", &status)

		assert.NotNil(t, err)
		assert.Equal(t, expectedErr, err)
		assert.EqualError(t, expectedErr, err.Error())
	})

	t.Run("should update an enrollment", func(t *testing.T) {
		expectedCounter := 1
		count := 0
		expetectedId := "1"
		expetectedStatus := "A"
		repo := &mockRepository{
			UpdateMock: func(ctx context.Context, id string, status *string) error {
				count++
				assert.Equal(t, expetectedId, id)
				assert.NotNil(t, status)
				assert.Equal(t, expetectedStatus, *status)
				return nil
			},
		}

		service := enrollment.NewService(l, repo, nil, nil)

		status := "A"
		err := service.Update(context.Background(), "1", &status)

		assert.Nil(t, err)
		assert.Equal(t, expectedCounter, count)
	})
}

func TestService_Count(t *testing.T) {
	l := log.New(io.Discard, "", 0)

	t.Run("should return an error", func(t *testing.T) {
		expCount := 1
		counter := 0
		expError := errors.New("some error")

		repo := &mockRepository{
			CountMock: func(ctx context.Context, filter enrollment.Filters) (int, error) {
				counter++
				return 0, errors.New("some error")
			},
		}

		service := enrollment.NewService(l, repo, nil, nil)

		count, err := service.Count(context.Background(), enrollment.Filters{})

		assert.NotNil(t, err)
		assert.Equal(t, expError, err)
		assert.Equal(t, expCount, counter)
		assert.Zero(t, count)
	})

	t.Run("should return the count of enrollments", func(t *testing.T) {
		expCounter := 1
		counter := 0
		expTotal := 10
		repo := &mockRepository{
			CountMock: func(ctx context.Context, filter enrollment.Filters) (int, error) {
				counter++
				return 10, nil
			},
		}

		service := enrollment.NewService(l, repo, nil, nil)
		total, err := service.Count(context.Background(), enrollment.Filters{})

		assert.Nil(t, err)
		assert.Equal(t, expTotal, total)
		assert.Equal(t, expCounter, counter)
	})
}
