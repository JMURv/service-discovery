package repo

import "errors"

var ErrNotFound = errors.New("not found")
var ErrAlreadyExists = errors.New("address already registered")
