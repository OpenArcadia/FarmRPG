package timer

import (
	"time"
)

type Timer struct {
	Duration  time.Duration
	Callback  func()
	startTime time.Time
	active    bool
}

func NewTimer(durationMs int, callback func()) *Timer {
	return &Timer{
		Duration: time.Duration(durationMs) * time.Millisecond,
		Callback: callback,
		active:   false,
	}
}

func (t *Timer) Activate() {
	t.active = true
	t.startTime = time.Now()
}

func (t *Timer) Deactivate() {
	t.active = false
	t.startTime = time.Time{}
}

func (t *Timer) Update() {
	if t.active && time.Since(t.startTime) >= t.Duration {
		if t.Callback != nil && !t.startTime.IsZero() {
			t.Callback()
		}
		t.Deactivate()
	}
}

func (t *Timer) IsActive() bool {
	return t.active
}
