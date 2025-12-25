package enrollment

import (
	"context"
	"errors"
	"log"

	"github.com/JuD4Mo/go_api_web_meta/meta"
	courseSDK "github.com/JuD4Mo/go_api_web_sdk/course"
	userSDK "github.com/JuD4Mo/go_api_web_sdk/user"
	"github.com/JuD4Mo/go_lib_response/response"
)

type (
	Controller func(ctx context.Context, request interface{}) (response interface{}, err error)

	Endpoints struct {
		Create Controller
		GetAll Controller
		Update Controller
	}

	CreateReq struct {
		UserId   string `json:"user_id"`
		CourseId string `json:"course_id"`
	}

	GetAllReq struct {
		UserID   string
		CourseID string
		Limit    int
		Page     int
	}

	UpdateReq struct {
		ID     string
		Status *string `json:"status"`
	}

	Response struct {
		Status int         `json:"status"`
		Data   interface{} `json:"data,omitempty"`
		Err    string      `json:"error,omitempty"`
		Meta   *meta.Meta  `json:"meta,omitempty"`
	}

	Config struct {
		LimitPage string
	}
)

func MakeEndpoints(s Service, config Config) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoint(s),
		GetAll: makeGetAllEndpoint(s, config),
		Update: makeUpdateEndpoint(s),
	}
}

func makeCreateEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateReq)

		if req.UserId == "" {
			return nil, response.BadRequest(ErrUserIdRequired.Error())
		}

		if req.CourseId == "" {
			return nil, response.BadRequest(ErrCourseIdRequired.Error())
		}

		enroll, err := s.Create(ctx, req.UserId, req.CourseId)
		if err != nil {
			log.Println(err)
			if errors.As(err, &userSDK.ErrNotFound{}) ||
				errors.As(err, &courseSDK.ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}

			return nil, response.InternalServerError(err.Error())
		}

		return response.Created("success", enroll, nil), nil
	}
}

func makeGetAllEndpoint(s Service, config Config) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetAllReq)

		filters := Filters{
			UserId:   req.UserID,
			CourseId: req.CourseID,
		}

		count, err := s.Count(ctx, filters)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		meta, err := meta.New(req.Page, req.Limit, count, config.LimitPage)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		enrollments, err := s.GetAll(ctx, filters, meta.Offset(), meta.Limit())
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", enrollments, meta), nil
	}
}

func makeUpdateEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateReq)

		if req.Status != nil && *req.Status == "" {
			return nil, response.BadRequest(ErrStatusRequired.Error())
		}

		if err := s.Update(ctx, req.ID, req.Status); err != nil {

			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}

			if errors.As(err, &ErrInvalidStatus{}) {
				return nil, response.BadRequest(err.Error())
			}

			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", nil, nil), nil
	}
}
