package main

import (
	"context"
	"fmt"
	"math/rand/v2"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mabrarov/mp/pkg/supervisor"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	s := supervisor.New(ctx)
	defer func() {
		s.Stop()
		_ = s.Wait()
	}()

	for i := range 10 {
		num := i
		s.Go(func(ctx context.Context) error {
			completed := false
			defer func() {
				if !completed {
					s.Stop()
				}
			}()
			err := run(ctx, s.Stop, num)
			completed = true
			return err
		})
	}

	if err := s.Wait(); err != nil {
		fmt.Printf("Completed with error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Completed without error")
}

func run(ctx context.Context, shutdown func(), num int) error {
	fmt.Printf("Started %d\n", num)
	mid := time.Tick(1 * time.Second)
	done := time.Tick(10 * time.Second)
	var err error
eventLoop:
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Stopped %d\n", num)
			break eventLoop
		case <-mid:
			n := rand.IntN(100)
			if n == 0 {
				err = fmt.Errorf("failed %d", num)
				fmt.Printf("Error %d\n", num)
				shutdown()
				break eventLoop
			}
			if n == 1 {
				panic(fmt.Sprintf("panic %d", num))
			}
		case <-done:
			fmt.Printf("Done %d\n", num)
			break eventLoop
		}
	}
	fmt.Printf("Completed %d\n", num)
	return err
}
