package misc

import "errors"

var (
	ErrNotFound       = errors.New("entity not found")
	ErrIDEmpty        = errors.New("empty id")
	ErrUsernameEmpty  = errors.New("empty username")
	ErrConig          = errors.New("incorrect config")
	ErrArgSuccessExit = errors.New("arg success exit")
	ErrRuntime        = errors.New("runtime error")
)
