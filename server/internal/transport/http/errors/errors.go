package errors

import "errors"

var (
	ErrInvalidMethod = errors.New("invalid request method")
	ErrNotFoundID    = errors.New("выражение с таким id не найдено")
	ErrNotTask       = errors.New("нету Task для решения")
)
