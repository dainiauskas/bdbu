package app

import (
	"fmt"
	"time"
)

// Duration structure
type Duration struct {
	Start time.Time
}

// NewDuration create duration structure with start time
func NewDuration() *Duration {
	return &Duration{
		Start: time.Now(),
	}
}

// Completed show duration to screen by template
func (d *Duration) Completed(tpl string) {
	fmt.Printf(tpl, d.end())
}

func (d *Duration) end() time.Duration {
	return time.Now().Sub(d.Start)
}
