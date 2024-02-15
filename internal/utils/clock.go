package utils

import "time"

type ClockInterface interface {
	Now() time.Time
}

type Clock struct {
	//
}

func (c *Clock) Now() time.Time {
	return time.Now()
}
