package test

import (
	"net/http"
	"testing"

	"github.com/JuD4Mo/go_api_web_domain/domain"
	"github.com/JuD4Mo/go_api_web_enrollment/internal/enrollment"
	"github.com/stretchr/testify/assert"
)

type dataResponse struct {
	Message string      `json:"message"`
	Status  int         `json:"status"`
	Data    interface{} `json:"data"`
	Meta    interface{} `json:"meta"`
}

func TestEnrollments(t *testing.T) {
	t.Run("should create and enrollment and get it", func(t *testing.T) {

		bodyReq := enrollment.CreateReq{
			UserId:   "123-test",
			CourseId: "222-test",
		}

		resp := cli.Post("/enrollments", bodyReq)
		assert.Nil(t, resp.Err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		dataCreated := domain.Enrollment{}
		dRespCreated := dataResponse{Data: &dataCreated}
		err := resp.FillUp(&dRespCreated)
		assert.Nil(t, err)

		assert.Equal(t, "success", dRespCreated.Message)
		assert.Equal(t, http.StatusCreated, dRespCreated.Status)

		assert.NotEmpty(t, dataCreated.ID)
		assert.Equal(t, "123-test", dataCreated.UserID)
		assert.Equal(t, "222-test", dataCreated.CourseID)

	})
}
