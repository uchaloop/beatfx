package beatfx

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"go.uber.org/fx"
)

// fakeShutdowner counts the requests it receives and how many options each
// carried. The options themselves are opaque, so the count is as far as a test
// can check that a code was given.
type fakeShutdowner struct {
	calls   atomic.Int64
	options atomic.Int64
}

func (f *fakeShutdowner) Shutdown(opts ...fx.ShutdownOption) error {
	f.calls.Add(1)
	f.options.Add(int64(len(opts)))

	return nil
}

// A loop that ends without being asked to takes the application down with it:
// a scheduler whose loop has left has nothing further to do, and leaving the
// process up would hide the failure.
func TestWatchLoop_ShutsDownWhenTheLoopEndsByItself(t *testing.T) {
	done, stopping := make(chan struct{}), make(chan struct{})
	var shutdowner fakeShutdowner

	finished := make(chan struct{})
	go func() {
		watchLoop(done, stopping, &shutdowner)
		close(finished)
	}()

	close(done)
	<-finished

	if n := shutdowner.calls.Load(); n != 1 {
		t.Errorf("shutdown requested %d times, want 1", n)
	}
	// fx.ExitCode(1), so anyone waiting on fx.App.Wait learns why.
	if n := shutdowner.options.Load(); n != 1 {
		t.Errorf("the request carried %d options, want 1 (the exit code)", n)
	}
}

// The ordinary path closes stopping before Stop, so the Done that follows must
// not be read as a failure. A false shutdown here would be worse than a missing
// one: it would fire on every clean stop.
func TestWatchLoop_QuietOnAnOrdinaryShutdown(t *testing.T) {
	done, stopping := make(chan struct{}), make(chan struct{})
	var shutdowner fakeShutdowner

	finished := make(chan struct{})
	go func() {
		watchLoop(done, stopping, &shutdowner)
		close(finished)
	}()

	close(stopping)
	close(done)
	<-finished

	if n := shutdowner.calls.Load(); n != 0 {
		t.Errorf("shutdown requested %d times during an ordinary stop, want 0", n)
	}
}

// Both notifications arriving together still count as an ordinary shutdown:
// Fx is already on its way out, so asking again would add nothing.
func TestWatchLoop_StoppingWinsWhenBothAreReady(t *testing.T) {
	done, stopping := make(chan struct{}), make(chan struct{})
	var shutdowner fakeShutdowner

	close(stopping)
	close(done)

	watchLoop(done, stopping, &shutdowner)

	// Whichever branch the select took, at most one request may follow.
	if n := shutdowner.calls.Load(); n > 1 {
		t.Errorf("shutdown requested %d times, want at most 1", n)
	}
}

// The watcher must not outlive the application when the loop never ends.
func TestWatchLoop_ReturnsOnStopping(t *testing.T) {
	done, stopping := make(chan struct{}), make(chan struct{})
	var shutdowner fakeShutdowner

	finished := make(chan struct{})
	go func() {
		watchLoop(done, stopping, &shutdowner)
		close(finished)
	}()

	close(stopping)

	select {
	case <-finished:
	case <-time.After(2 * time.Second):
		t.Fatal("the watcher did not return after stopping")
	}
}

// Exercise the real Fx signal so ExitCode(0) cannot satisfy this regression.
func TestWatchLoop_ReportsFailureToFx(t *testing.T) {
	var shutdowner fx.Shutdowner
	app := fx.New(fx.NopLogger, fx.Populate(&shutdowner))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := app.Start(ctx); err != nil {
		t.Fatal(err)
	}
	defer func() {
		cleanup, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if err := app.Stop(cleanup); err != nil {
			t.Error(err)
		}
	}()
	signals := app.Wait()
	done := make(chan struct{})
	close(done)
	watchLoop(done, make(chan struct{}), shutdowner)
	select {
	case signal := <-signals:
		if signal.ExitCode != 1 {
			t.Fatalf("exit code = %d, want 1", signal.ExitCode)
		}
	case <-ctx.Done():
		t.Fatal("Fx did not receive the shutdown request")
	}
}
