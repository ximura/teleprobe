package async

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Shutdown struct {
	done chan struct{}
}

func NewShutdown() *Shutdown {
	return &Shutdown{
		done: make(chan struct{}),
	}
}

// Run runs the shutdown service activity
func (s *Shutdown) Run(ctx context.Context) error {
	// listen for interrupt, terminated, and hangup signals and gracefully shutdown server
	log.Println("listening signals...")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	select {
	case sig := <-c:
		// Allow interrupt signals to be caught again in worse-case scenario
		// situations when the service hangs during a graceful shutdown.
		log.Println("graceful shutdown...")
		signal.Reset(os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
		return fmt.Errorf("signal: %s", sig.String())
	case <-ctx.Done():
		return ctx.Err()
	case <-s.done:
		return nil
	}
}

// Close closes the shutdown service activity
func (s *Shutdown) Close() error {
	close(s.done)
	return nil
}
