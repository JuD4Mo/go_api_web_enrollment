package enrollment

import "errors"

var ErrUserIdRequired = errors.New("user id is required")
var ErrCourseIdRequired = errors.New("course id is required")
