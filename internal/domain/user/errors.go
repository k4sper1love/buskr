package user

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrAccessDenied      = errors.New("access denied")
	ErrInvalidStatus     = errors.New("invalid status for this action")
	ErrInviteNotFound    = errors.New("invite not found")
	ErrInviteAlreadyUsed = errors.New("invite already used")
	ErrInviteExpired     = errors.New("invite has expired")
)
