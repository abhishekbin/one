package onecache

import (
	"context"
	"time"
)

type CancelFn func()

type Scheduler interface {
	// RunOnceAfter Schedules a function to be run one time after a dely.  Returns a cancellation function
	// that prevents fn from being run.
	RunOnceAfter(ctx context.Context, duration time.Duration, fn func()) (CancelFn, error)
}

type DefaultScheduler struct{}

func (d DefaultScheduler) RunOnceAfter(ctx context.Context, duration time.Duration, fn func()) (CancelFn, error) {
	timer := time.AfterFunc(duration, fn)

	cancelFn := func() { timer.Stop() }
	return cancelFn, nil
}
