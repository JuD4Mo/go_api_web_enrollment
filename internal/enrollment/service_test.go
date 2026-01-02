package enrollment_test

import (
	"context"
	"errors"
	"io"
	"log"
	"testing"

	"github.com/JuD4Mo/go_api_web_domain/domain"
	"github.com/JuD4Mo/go_api_web_enrollment/internal/enrollment"
	courseSdk "github.com/JuD4Mo/go_api_web_sdk/course/mock"
	userSdk "github.com/JuD4Mo/go_api_web_sdk/user/mock"
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

func TestService_Create(t *testing.T) {
	l := log.New(io.Discard, "", 0)
	t.Run("should return an error in user sdk", func(t *testing.T) {
		expectedErr := errors.New("some error")
		expectedCounter := 1
		counter := 0

		userSdk := &userSdk.UserSdkMock{
			GetMock: func(id string) (*domain.User, error) {
				counter++
				return nil, errors.New("some error")
			},
		}

		service := enrollment.NewService(l, nil, userSdk, nil)

		enrollment, err := service.Create(context.Background(), "11", "22")

		assert.NotNil(t, err)
		assert.Equal(t, expectedCounter, counter)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, enrollment)
	})

	t.Run("should return an error in course sdk", func(t *testing.T) {
		expectedErr := errors.New("some error")
		expectedCounter := 2
		counter := 0
		userSdk := &userSdk.UserSdkMock{
			GetMock: func(id string) (*domain.User, error) {
				counter++
				return nil, nil
			},
		}

		courseSdk := &courseSdk.CourseSdkMock{
			GetMock: func(id string) (*domain.Course, error) {
				counter++
				return nil, errors.New("some error")
			},
		}

		service := enrollment.NewService(l, nil, userSdk, courseSdk)

		enrollment, err := service.Create(context.Background(), "11", "22")

		assert.NotNil(t, err)
		assert.Equal(t, expectedCounter, counter)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, enrollment)
	})

	t.Run("should return an error in repository", func(t *testing.T) {
		expectedErr := errors.New("some error")
		expectedCounter := 3
		counter := 0
		userSdk := &userSdk.UserSdkMock{
			GetMock: func(id string) (*domain.User, error) {
				counter++
				return nil, nil
			},
		}

		courseSdk := &courseSdk.CourseSdkMock{
			GetMock: func(id string) (*domain.Course, error) {
				counter++
				return nil, nil
			},
		}

		repo := &mockRepository{
			CreateMock: func(ctx context.Context, enroll *domain.Enrollment) error {
				counter++
				return errors.New("some error")
			},
		}

		service := enrollment.NewService(l, repo, userSdk, courseSdk)

		enrollment, err := service.Create(context.Background(), "11", "22")

		assert.NotNil(t, err)
		assert.Equal(t, expectedCounter, counter)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, enrollment)
	})

	t.Run("should create enrollment", func(t *testing.T) {
		expectedCounter := 3
		counter := 0
		expectedUserId := "11"
		expectedCourseId := "22"
		expectedStatus := domain.Pending
		expectedId := "123"
		userSdk := &userSdk.UserSdkMock{
			GetMock: func(id string) (*domain.User, error) {
				counter++
				assert.Equal(t, expectedUserId, id)
				return nil, nil
			},
		}

		courseSdk := &courseSdk.CourseSdkMock{
			GetMock: func(id string) (*domain.Course, error) {
				counter++
				assert.Equal(t, expectedCourseId, id)
				return nil, nil
			},
		}

		repo := &mockRepository{
			CreateMock: func(ctx context.Context, enroll *domain.Enrollment) error {
				counter++
				enroll.ID = "123"
				return nil
			},
		}

		service := enrollment.NewService(l, repo, userSdk, courseSdk)

		enrollment, err := service.Create(context.Background(), "11", "22")

		assert.Nil(t, err)
		assert.Equal(t, expectedCounter, counter)
		assert.NotNil(t, enrollment)
		assert.Equal(t, expectedUserId, enrollment.UserID)
		assert.Equal(t, expectedCourseId, enrollment.CourseID)
		assert.Equal(t, expectedStatus, enrollment.Status)
		assert.Equal(t, expectedId, enrollment.ID)
	})
}
