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
	ctx := context.Background()
	router := signal.New(ctx)
	waiter := sync.WaitGroup{}
	waiter.Add(1)

	router.Handle(syscall.SIGUSR1, func(sig os.Signal) {
		fmt.Println("signal called:", sig)
		waiter.Done()
	})

	go func() {
		if e := router.Start(); e != nil {
			panic(e)
		}
	}()

	// Simulate a signal
	router.Fire(syscall.SIGUSR1)

	// Wait for the signal to be handled
	waiter.Wait()

	// Stop the router
	router.Stop(nil)

	// Output: signal called: user defined signal 1
}
