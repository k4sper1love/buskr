package location

import "errors"

var (
	ErrLocationNotFound = errors.New("location not found")
	ErrLocationInactive = errors.New("location is currently inactive")
)
