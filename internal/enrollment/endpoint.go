package enrollment

import (
	"context"

	"github.com/JuD4Mo/go_api_web_meta/meta"
	"github.com/JuD4Mo/go_lib_response/response"
)

type (
	Controller func(ctx context.Context, request interface{}) (response interface{}, err error)

	Endpoints struct {
		Create Controller
	}

	CreateReq struct {
		UserId   string `json:"user_id"`
		CourseId string `json:"course_id"`
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

		enroll, err := s.Create(req.UserId, req.CourseId)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.Created("success", enroll, nil), nil
	}
}
