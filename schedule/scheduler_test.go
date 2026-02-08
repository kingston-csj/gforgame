package schedule_test

import (
	"io/github/gforgame/schedule"
	"testing"
	"time"
)

func TestDefaultTaskScheduler_Schedule(t *testing.T) {
	scheduler := schedule.NewDefaultTaskScheduler()

	// 1. Test normal schedule
	done := make(chan struct{})
	start := time.Now()
	delayMs := int64(100)
	
	_, err := scheduler.Schedule(func() {
		close(done)
	}, delayMs)
	
	if err != nil {
		t.Fatalf("Schedule failed: %v", err)
	}

	select {
	case <-done:
		elapsed := time.Since(start)
		if elapsed < time.Duration(delayMs)*time.Millisecond {
			t.Errorf("Task executed too early: %v", elapsed)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Task timeout")
	}

	// 2. Test cancel
	done2 := make(chan struct{})
	cancellable, err := scheduler.Schedule(func() {
		close(done2)
	}, 200)
	if err != nil {
		t.Fatalf("Schedule failed: %v", err)
	}

	if !cancellable.Cancel() {
		t.Error("Cancel failed")
	}

	select {
	case <-done2:
		t.Error("Task should be cancelled")
	case <-time.After(300 * time.Millisecond):
		// Expected
	}
}
