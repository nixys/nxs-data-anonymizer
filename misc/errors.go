package misc

import "errors"

var (

	// ErrNotFound defines an error when entity not found
	ErrNotFound = errors.New("entity not found")

	// ErrIDEmpty defines an error when id is empty
	ErrIDEmpty = errors.New("empty id")

	// ErrUsernameEmpty defines an error when username is empty
	ErrUsernameEmpty = errors.New("empty username")
)
