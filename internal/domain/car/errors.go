package car

import "errors"

var (
	ErrNotFound  = errors.New("car not found")
	ErrEmptyName = errors.New("car name is empty")
	ErrAlreadyExists = errors.New("car already exists")
)
