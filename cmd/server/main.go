package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mp/pkg/supervisor"
)

func main() {
	s := supervisor.New()
	defer func() {
		s.Stop()
		_ = s.Wait()
	}()

	userStop := make(chan os.Signal, 1)
	signal.Notify(userStop, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sg := <-userStop
		fmt.Printf("Got signal %v, exiting...\n", sg)
		s.Stop()
	}()

	for i := range 10 {
		num := i
		s.Run(func(stop <-chan supervisor.StopToken) error {
			return run(stop, s.Stop, num)
		})
	}

	if err := s.Wait(); err != nil {
		fmt.Printf("Completed with error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Completed without error")
}

func run(stop <-chan supervisor.StopToken, shutdown func(), num int) error {
	fmt.Printf("Started %d\n", num)
	mid := time.Tick(1 * time.Second)
	done := time.Tick(10 * time.Second)
	var err error
eventLoop:
	for {
		select {
		case <-stop:
			fmt.Printf("Stopped %d\n", num)
			break eventLoop
		case <-mid:
			if rand.IntN(100) == 0 {
				err = fmt.Errorf("failed %d", num)
				fmt.Printf("Error %d\n", num)
				shutdown()
				break eventLoop
			}
		case <-done:
			fmt.Printf("Done %d\n", num)
			break eventLoop
		}
	}
	fmt.Printf("Completed %d\n", num)
	return err
}
