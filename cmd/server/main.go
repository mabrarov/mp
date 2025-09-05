package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mabrarov/mp/pkg/panicerr"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, signalStop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer signalStop()
	ctx, stop := context.WithCancel(ctx)

	var s errgroup.Group
	defer func() {
		stop()
		_ = s.Wait()
	}()

	for i := range 10 {
		id := i
		s.Go(wrapPanic(stop, func() error {
			return run(ctx, stop, id)
		}))
	}

	if err := s.Wait(); err != nil {
		fmt.Printf("Completed with error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Completed without error")
}

func wrapPanic(stop func(), f func() error) func() error {
	return func() (err error) {
		defer func() {
			if r := recover(); r != nil {
				pe := panicerr.New(r)
				if err == nil {
					err = pe
				} else {
					err = errors.Join(err, pe)
				}
				stop()
			}
		}()
		return f()
	}
}

func run(ctx context.Context, shutdown func(), id int) error {
	fmt.Printf("Started %d\n", id)
	mid := time.Tick(1 * time.Second)
	done := time.Tick(10 * time.Second)
	var err error
eventLoop:
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Stopped %d\n", id)
			break eventLoop
		case <-mid:
			n := rand.IntN(100)
			if n == 0 {
				err = fmt.Errorf("failed %d", id)
				fmt.Printf("Error %d\n", id)
				shutdown()
				break eventLoop
			}
			if n == 1 {
				panic(fmt.Sprintf("panic %d", id))
			}
		case <-done:
			fmt.Printf("Done %d\n", id)
			break eventLoop
		}
	}
	fmt.Printf("Completed %d\n", id)
	return err
}
