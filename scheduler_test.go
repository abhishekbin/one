package onecache

import (
	"context"
	"testing"
	"time"
)

func TestDefaultScheduler(t *testing.T) {
	t.Parallel()

	// TODO: stop sleeping in tests

	var (
		ctx = context.Background()

		scheduler  DefaultScheduler
		iterations int
	)

	run := func() { iterations++ }

	// Schedule work and ensure our function has not run yet

	_, err := scheduler.RunOnceAfter(ctx, 20*time.Millisecond, run)
	assertNoError(t, err)
	assertEquals(t, iterations, 0)

	// Sleep a bit to ensure our work runs

	time.Sleep(30 * time.Millisecond)
	assertEquals(t, iterations, 1)

	// Schedule some work then stop it

	close, err := scheduler.RunOnceAfter(ctx, 30*time.Millisecond, run)
	assertNoError(t, err)
	assertEquals(t, iterations, 1)

	// Run the close function to stop the work

	close()

	// Sleep a bit and ensure the work didn't run

	time.Sleep(30 * time.Millisecond)
	assertEquals(t, iterations, 1)
}
