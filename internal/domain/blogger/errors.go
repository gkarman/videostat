package blogger

import "errors"

var (
	ErrUrlInvalid = errors.New("url is invalid")
	ErrUrlExist = errors.New("url already exists")
	ErrBloggerNotFound = errors.New("blogger not found")
    ErrConcurrentUpdate = errors.New("concurrent status update")
	ErrVideoNotFound = errors.New("video not found")
)
