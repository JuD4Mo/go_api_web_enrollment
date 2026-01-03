package enrollment_test

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/JuD4Mo/go_api_web_domain/domain"
	"github.com/JuD4Mo/go_api_web_enrollment/internal/enrollment"
	"github.com/JuD4Mo/go_api_web_sdk/course"
	courseSdk "github.com/JuD4Mo/go_api_web_sdk/course"
	courseSdkMock "github.com/JuD4Mo/go_api_web_sdk/course/mock"
	"github.com/JuD4Mo/go_api_web_sdk/user"
	userSdk "github.com/JuD4Mo/go_api_web_sdk/user"
	userSdkMock "github.com/JuD4Mo/go_api_web_sdk/user/mock"
	"github.com/JuD4Mo/go_lib_response/response"
	"github.com/stretchr/testify/assert"
)

func TestCreateEndpoint(t *testing.T) {
	l := log.New(io.Discard, "", 0)

	t.Run("should return bad request when user id is empty", func(t *testing.T) {
		endpoint := enrollment.MakeEndpoints(nil, enrollment.Config{})
		_, err := endpoint.Create(context.Background(), enrollment.CreateReq{})
		assert.Error(t, err)

		resp := err.(response.Response)

		assert.EqualError(t, enrollment.ErrUserIdRequired, resp.Error())
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode())
	})

	t.Run("should return bad request when course id is empty", func(t *testing.T) {
		endpoint := enrollment.MakeEndpoints(nil, enrollment.Config{})
		_, err := endpoint.Create(context.Background(), enrollment.CreateReq{UserId: "1234"})
		assert.Error(t, err)

		resp := err.(response.Response)

		assert.EqualError(t, enrollment.ErrCourseIdRequired, resp.Error())
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode())
	})

	obj := []struct {
		tag              string
		repositoryMock   enrollment.Repository
		userSdkMock      user.Transport
		courseSdkMock    course.Transport
		expectedErr      error
		expectedStatus   int
		expectedResponse *domain.Enrollment
	}{
		{

			tag: "should return an error if user skd returns an unexpected error",
			userSdkMock: &userSdkMock.UserSdkMock{
				GetMock: func(id string) (*domain.User, error) {
					return nil, errors.New("unexpected error")
				},
			},
			expectedErr:    errors.New("unexpected error"),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			tag: "should return an error if user does not exist",
			userSdkMock: &userSdkMock.UserSdkMock{
				GetMock: func(id string) (*domain.User, error) {
					return nil, userSdk.ErrNotFound{Message: "user not found"}
				},
			},
			expectedErr:    userSdk.ErrNotFound{Message: "user not found"},
			expectedStatus: http.StatusNotFound,
		},
		{
			tag: "should return an error if course skd returns an unexpected error",
			userSdkMock: &userSdkMock.UserSdkMock{
				GetMock: func(id string) (*domain.User, error) {
					return nil, nil
				},
			},
			courseSdkMock: &courseSdkMock.CourseSdkMock{
				GetMock: func(id string) (*domain.Course, error) {
					return nil, errors.New("unexpected error")
				},
			},
			expectedErr:    errors.New("unexpected error"),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			tag: "should return an error if course does not exist",
			userSdkMock: &userSdkMock.UserSdkMock{
				GetMock: func(id string) (*domain.User, error) {
					return nil, nil
				},
			},
			courseSdkMock: &courseSdkMock.CourseSdkMock{
				GetMock: func(id string) (*domain.Course, error) {
					return nil, courseSdk.ErrNotFound{Message: "course not found"}
				},
			},
			expectedErr:    courseSdk.ErrNotFound{Message: "course not found"},
			expectedStatus: http.StatusNotFound,
		},
		{
			tag: "should return an error if repository returns an unexpected error",
			userSdkMock: &userSdkMock.UserSdkMock{
				GetMock: func(id string) (*domain.User, error) {
					return nil, nil
				},
			},
			courseSdkMock: &courseSdkMock.CourseSdkMock{
				GetMock: func(id string) (*domain.Course, error) {
					return nil, nil
				},
			},
			repositoryMock: &mockRepository{
				CreateMock: func(ctx context.Context, enrollment *domain.Enrollment) error {
					return errors.New("unexpected error")
				},
			},
			expectedErr:    errors.New("unexpected error"),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			tag: "should return the enrollment",
			userSdkMock: &userSdkMock.UserSdkMock{
				GetMock: func(id string) (*domain.User, error) {
					return nil, nil
				},
			},
			courseSdkMock: &courseSdkMock.CourseSdkMock{
				GetMock: func(id string) (*domain.Course, error) {
					return nil, nil
				},
			},
			repositoryMock: &mockRepository{
				CreateMock: func(ctx context.Context, enrollment *domain.Enrollment) error {
					enrollment.ID = "10010"
					return nil
				},
			},
			expectedStatus: http.StatusCreated,
			expectedResponse: &domain.Enrollment{
				ID:       "10010",
				UserID:   "1",
				CourseID: "4",
				Status:   "P",
			},
		},
	}
	for _, obj := range obj {
		t.Run(obj.tag, func(t *testing.T) {
			service := enrollment.NewService(l, obj.repositoryMock, obj.userSdkMock, obj.courseSdkMock)
			endpoint := enrollment.MakeEndpoints(service, enrollment.Config{})
			resp, err := endpoint.Create(context.Background(), enrollment.CreateReq{UserId: "1", CourseId: "4"})

			if obj.expectedErr != nil {
				assert.NotNil(t, err)
				assert.Nil(t, resp)

				r := err.(response.Response)
				assert.EqualError(t, obj.expectedErr, r.Error())
				assert.Equal(t, obj.expectedStatus, r.StatusCode())
			} else {
				assert.NotNil(t, resp)
				assert.Nil(t, err)

				r := resp.(response.Response)
				assert.Equal(t, obj.expectedStatus, r.StatusCode())
				assert.Empty(t, r.Error())

				enrollment := r.GetData().(*domain.Enrollment)
				assert.Equal(t, obj.expectedResponse.ID, enrollment.ID)
				assert.Equal(t, obj.expectedResponse.UserID, enrollment.UserID)
				assert.Equal(t, obj.expectedResponse.CourseID, enrollment.CourseID)
				assert.Equal(t, obj.expectedResponse.Status, enrollment.Status)
			}

		})
	}
}

func TestGetAllEndpoint(t *testing.T) {
	l := log.New(io.Discard, "", 0)

	t.Run("should return an error if Count returns an unexpected error", func(t *testing.T) {
		wantErr := errors.New("unexpected error")
		service := enrollment.NewService(l, &mockRepository{
			CountMock: func(ctx context.Context, filters enrollment.Filters) (int, error) {
				return 0, errors.New("unexpected error")
			},
		}, nil, nil)
		endpoint := enrollment.MakeEndpoints(service, enrollment.Config{})
		_, err := endpoint.GetAll(context.Background(), enrollment.GetAllReq{})
		assert.Error(t, err)

		resp := err.(response.Response)
		assert.EqualError(t, wantErr, resp.Error())
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode())
	})

	t.Run("should return an error if meta returns a parsing error", func(t *testing.T) {
		wantErr := errors.New("strconv.Atoi: parsing \"invalid number\": invalid syntax")
		service := enrollment.NewService(l, &mockRepository{
			CountMock: func(ctx context.Context, filters enrollment.Filters) (int, error) {
				return 3, nil
			},
		}, nil, nil)
		endpoint := enrollment.MakeEndpoints(service, enrollment.Config{LimitPage: "invalid number"})
		_, err := endpoint.GetAll(context.Background(), enrollment.GetAllReq{})
		assert.Error(t, err)

		resp := err.(response.Response)
		assert.EqualError(t, wantErr, resp.Error())
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode())
	})

	t.Run("should return an error if GetAll repository returns an unexpected error", func(t *testing.T) {
		wantErr := errors.New("unexpected error")
		service := enrollment.NewService(l, &mockRepository{
			CountMock: func(ctx context.Context, filters enrollment.Filters) (int, error) {
				return 3, nil
			},
			GetAllMock: func(ctx context.Context, filters enrollment.Filters, offset, limit int) ([]domain.Enrollment, error) {
				return nil, errors.New("unexpected error")
			},
		}, nil, nil)
		endpoint := enrollment.MakeEndpoints(service, enrollment.Config{LimitPage: "10"})
		_, err := endpoint.GetAll(context.Background(), enrollment.GetAllReq{})
		assert.Error(t, err)

		resp := err.(response.Response)
		assert.EqualError(t, wantErr, resp.Error())
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode())
	})

	t.Run("should return the enrollments", func(t *testing.T) {
		wantEnrollments := []domain.Enrollment{
			{ID: "1", UserID: "11", CourseID: "111", Status: "P"},
			{ID: "2", UserID: "22", CourseID: "222", Status: "P"},
			{ID: "3", UserID: "33", CourseID: "333", Status: "P"},
		}
		service := enrollment.NewService(l, &mockRepository{
			CountMock: func(ctx context.Context, filters enrollment.Filters) (int, error) {
				return 3, nil
			},
			GetAllMock: func(ctx context.Context, filters enrollment.Filters, offset, limit int) ([]domain.Enrollment, error) {
				return []domain.Enrollment{
					{ID: "1", UserID: "11", CourseID: "111", Status: "P"},
					{ID: "2", UserID: "22", CourseID: "222", Status: "P"},
					{ID: "3", UserID: "33", CourseID: "333", Status: "P"},
				}, nil
			},
		}, nil, nil)
		endpoint := enrollment.MakeEndpoints(service, enrollment.Config{LimitPage: "10"})
		resp, err := endpoint.GetAll(context.Background(), enrollment.GetAllReq{})
		assert.Nil(t, err)

		r := resp.(response.Response)
		assert.Equal(t, http.StatusOK, r.StatusCode())
		assert.Empty(t, r.Error())

		enrollments := r.GetData().([]domain.Enrollment)
		assert.Equal(t, wantEnrollments, enrollments)

	})
}

func TestUpdateEndpoint(t *testing.T) {
	l := log.New(io.Discard, "", 0)

	t.Run("should return an error if status is empty", func(t *testing.T) {
		endpoint := enrollment.MakeEndpoints(nil, enrollment.Config{})
		status := ""
		_, err := endpoint.Update(context.Background(), enrollment.UpdateReq{Status: &status})
		assert.Error(t, err)

		resp := err.(response.Response)
		assert.EqualError(t, enrollment.ErrStatusRequired, resp.Error())
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode())
	})

	t.Run("should return an error if repository returns a not found error", func(t *testing.T) {
		service := enrollment.NewService(l, &mockRepository{
			UpdateMock: func(ctx context.Context, id string, status *string) error {
				return enrollment.ErrNotFound{EnrollmentId: id}
			},
		}, nil, nil)
		endpoint := enrollment.MakeEndpoints(service, enrollment.Config{})
		status := "A"
		_, err := endpoint.Update(context.Background(), enrollment.UpdateReq{ID: "20", Status: &status})
		assert.Error(t, err)

		resp := err.(response.Response)
		assert.EqualError(t, enrollment.ErrNotFound{EnrollmentId: "20"}, resp.Error())
		assert.Equal(t, http.StatusNotFound, resp.StatusCode())
	})

	t.Run("should return an error if repository returns a unexpected error", func(t *testing.T) {
		wantErr := errors.New("unexpected error")
		service := enrollment.NewService(l, &mockRepository{
			UpdateMock: func(ctx context.Context, id string, status *string) error {
				return errors.New("unexpected error")
			},
		}, nil, nil)
		endpoint := enrollment.MakeEndpoints(service, enrollment.Config{})
		status := "A"
		_, err := endpoint.Update(context.Background(), enrollment.UpdateReq{ID: "20", Status: &status})
		assert.Error(t, err)

		resp := err.(response.Response)
		assert.EqualError(t, wantErr, resp.Error())
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode())
	})

	t.Run("should return success", func(t *testing.T) {
		service := enrollment.NewService(l, &mockRepository{
			UpdateMock: func(ctx context.Context, id string, status *string) error {
				assert.Equal(t, "20", id)
				assert.NotNil(t, status)
				assert.Equal(t, "A", *status)
				return nil
			},
		}, nil, nil)
		endpoint := enrollment.MakeEndpoints(service, enrollment.Config{})
		status := "A"
		resp, err := endpoint.Update(context.Background(), enrollment.UpdateReq{ID: "20", Status: &status})
		assert.Nil(t, err)

		r := resp.(response.Response)
		assert.Equal(t, http.StatusOK, r.StatusCode())
		assert.Empty(t, r.Error())
		assert.Nil(t, r.GetData())
	})
}
