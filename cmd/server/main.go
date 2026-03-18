package main

import (
	"context"
	"fmt"
	"math/rand/v2"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, stop := context.WithCancel(context.Background())
	group, ctx := errgroup.WithContext(ctx)
	defer func() {
		stop()
		_ = group.Wait()
	}()

	signalChan := make(chan os.Signal, 2)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		signal.Stop(signalChan)
	}()
	go func() {
		<-signalChan
		fmt.Println("Shutdown signal received")
		stop()
	}()

	for i := range 10 {
		id := i
		group.Go(func() error {
			return run(ctx, id)
		})
	}

	if err := group.Wait(); err != nil {
		fmt.Printf("Completed with error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Completed without error")
}

func run(ctx context.Context, id int) error {
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
