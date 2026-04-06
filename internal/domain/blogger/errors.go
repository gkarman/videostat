package blogger

import "errors"

var (
	ErrUrlInvalid = errors.New("url is invalid")
	ErrUrlExist = errors.New("url already exists")
)
