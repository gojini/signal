package signal_test

import (
	"context"
	"fmt"
	"os"
	"sync"
	"syscall"

	"gojini.dev/signal"
)

func ExampleRouter() {
	ctx, cancel := context.WithCancel(context.Background())
	router, startF, stopF := signal.New(ctx, cancel)
	waiter := sync.WaitGroup{}
	waiter.Add(1)

	router.Handle(syscall.SIGUSR1, func(sig os.Signal) {
		fmt.Println("signal called:", sig)
		waiter.Done()
	})

	go func() {
		if e := startF(); e != nil {
			panic(e)
		}
	}()

	// Simulate a signal
	router.Fire(syscall.SIGUSR1)

	// Wait for the signal to be handled
	waiter.Wait()

	// Stop the router
	stopF(nil)

	// Output: signal called: user defined signal 1
}
