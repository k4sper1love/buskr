package booking

import "errors"

var (
	ErrBookingNotFound        = errors.New("booking not found")
	ErrSlotTaken              = errors.New("time slot is already taken")
	ErrNoiseExceeded          = errors.New("equipment noise level exceeds location limits")
	ErrUserNotActive          = errors.New("booking is only available for active users")
	ErrLocationInactive       = errors.New("location is temporarily unavailable for booking")
	ErrInvalidTime            = errors.New("invalid booking time")
	ErrInvalidStatus          = errors.New("invalid booking status for this action")
	ErrAlreadyCheckedIn       = errors.New("booking is already checked in")
	ErrAccessDenied           = errors.New("user is not the owner of this booking")
	ErrMaxActiveBookings      = errors.New("max active bookings reached")
	ErrMaxBookingsPerLocation = errors.New("max bookings per location at same day reached")
	ErrTimeOverlap            = errors.New("time overlap with another booking")
	ErrTooFarInFuture         = errors.New("booking is too far in future")
	ErrNoisyNeighbor          = errors.New("adjacent location has a loud artist booked at this time")
)
