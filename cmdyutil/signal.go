package cmdyutil

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"time"

	"github.com/ArtProcessors/cmdy"
)

const DefaultInterruptTimeout = 2 * time.Second

const ExitInterrupt = 130 // 128 + signal 2

type InterruptRunnerOption func(i *InterruptRunner)

// InterruptRunner is a cmdy.Runner that cancels the cmdy.Context passed to your command
// when an Interrupt signal is received by your app.
//
// If you send another interrupt while InterruptRunner is waiting for your command to
// exit, it aborts your command immediately and does not wait for shutdown to complete.
//
// NOTE: This API is experimental.
type InterruptRunner struct {
	*cmdy.Runner
	timeout   time.Duration
	onAbort   error
	onTimeout error
}

// InterruptTimeout allows you to configure how long the Runner will wait
// before giving up waiting for your command to shut itself down.
func InterruptTimeout(t time.Duration) InterruptRunnerOption {
	return func(i *InterruptRunner) { i.timeout = t }
}

func InterruptAbortErr(err error) InterruptRunnerOption {
	return func(i *InterruptRunner) { i.onAbort = err }
}

func InterruptTimeoutErr(err error) InterruptRunnerOption {
	return func(i *InterruptRunner) { i.onTimeout = err }
}

func NewInterruptRunner(runner *cmdy.Runner, opts ...InterruptRunnerOption) *InterruptRunner {
	rn := &InterruptRunner{
		Runner:    runner,
		onAbort:   ErrInterruptAborted,
		onTimeout: ErrInterruptTimeout,
	}
	for _, o := range opts {
		o(rn)
	}
	return rn
}

// InterruptibleRun is the interruptible equivalent of cmdy.Run: it creates
// an InterrptRunner with all defaults set and runs the command produced
// by your cmdy.Builder.
//
// NOTE: This API is experimental.
func InterruptibleRun(ctx context.Context, args []string, b cmdy.Builder) (rerr error) {
	return NewInterruptRunner(cmdy.DefaultRunner()).Run(ctx, cmdy.ProgName(), args, b)
}

// Run the command created by builder. If the program receives an os.Interrupt,
// ctx will be cancelled. If your command does not handle the 'ctx.Done()' condition
// in time, Run will return an error.
//
// A goroutine will be leaked if your command never completes in response to the
// Interrupt.
func (r *InterruptRunner) Run(ctx context.Context, name string, args []string, builder cmdy.Builder) (rerr error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan error, 1)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	go func() {
		done <- r.Runner.Run(ctx, name, args, builder)
	}()

	select {
	case <-sig:
		cancel()
	case err := <-done:
		return err
	}

	timeout := r.timeout
	if timeout <= 0 {
		timeout = DefaultInterruptTimeout
	}

	wait := time.After(timeout)

	select {
	case err := <-done:
		return err
	case <-sig:
		return ErrInterruptAborted
	case <-wait:
		return ErrInterruptTimeout
	}
}

func IsInterruptErr(err error) bool {
	return errors.Is(err, ErrInterruptAborted) || errors.Is(err, ErrInterruptTimeout)
}

var (
	ErrInterruptAborted cmdy.Error = &errInterruptAborted{}
	ErrInterruptTimeout cmdy.Error = &errInterruptTimeout{}
)

type errInterruptAborted struct{}

func (*errInterruptAborted) Error() string     { return "aborted!" }
func (*errInterruptAborted) Is(err error) bool { return err == ErrInterruptAborted }
func (*errInterruptAborted) Code() int         { return ExitInterrupt }

type errInterruptTimeout struct {
	inner error
}

func (*errInterruptTimeout) Timeout() bool     { return true }
func (*errInterruptTimeout) Error() string     { return "timeout waiting for shutdown!" }
func (*errInterruptTimeout) Is(err error) bool { return err == ErrInterruptTimeout }
func (*errInterruptTimeout) Code() int         { return ExitInterrupt }
func (e *errInterruptTimeout) Unwrap() error   { return e.inner }
