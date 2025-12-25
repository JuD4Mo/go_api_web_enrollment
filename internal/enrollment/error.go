package enrollment

import (
	"errors"
	"fmt"
)

var ErrUserIdRequired = errors.New("user id is required")
var ErrCourseIdRequired = errors.New("course id is required")
var ErrStatusRequired = errors.New("status is required")

type ErrNotFound struct {
	EnrollmentId string
}

type ErrInvalidStatus struct {
	Status string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("enrollment '%s' does not exist", e.EnrollmentId)
}

func (e ErrInvalidStatus) Error() string {
	return fmt.Sprintf("invalid '%s' status", e.Status)
}
