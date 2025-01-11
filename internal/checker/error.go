package checker

import (
	"errors"
	"fmt"
)

var ErrCheckService = errors.New("service health check failed")
var ErrNotServing = errors.New("service is not in a serving state")

type ErrUnexpectedStatusCode struct{ StatusCode int }

func (e ErrUnexpectedStatusCode) Error() string {
	return fmt.Sprintf("unexpected status code: %v", e.StatusCode)
}
