package tz

import "time"

var loc *time.Location = time.UTC

func Init(timezone string) error {
	l, err := time.LoadLocation(timezone)
	if err != nil {
		return err
	}

	loc = l
	return nil
}

func Location() *time.Location {
	return loc
}

func Now() time.Time {
	return time.Now().In(loc)
}

func In(t time.Time) time.Time {
	return t.In(loc)
}
