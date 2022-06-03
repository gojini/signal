package signal

import (
	"context"
	"os"
	"os/signal"
	"sync"
)

// Handler is a function that handles the signal. Handler can be registered
// with a signal router. It is invoked whenever as signal is fired.
type Handler func(os.Signal)

// StartFunc that will start the signal router, this method will not return
// until the router is stopped.
type StartFunc func() error

// StopFunc that will stop the signal router and clean up any resources used.
type StopFunc func(error)

// Router routes signals to registered handler. Router keeps track of signals
// that needs to be ignored or handled. If a handler is set for a signal it will
// invoke the handler when a signal is received.
type Router struct {
	signalCh   chan os.Signal
	signals    map[os.Signal]Handler
	ignSignals map[os.Signal]struct{}
	lock       *sync.RWMutex
}

// New returns a signal router, a start handler and stop handler for the router.
func New(ctx context.Context, cancel context.CancelFunc) (*Router, StartFunc, StopFunc) {
	channelSize := 10
	signalCh := make(chan os.Signal, channelSize)
	signals := make(map[os.Signal]Handler)
	lock := &sync.RWMutex{}

	startF := func() error {
		// This go routine dies with the server
		for {
			select {
			case <-ctx.Done():
				// Context got cancelled, exit
				return nil
			case sig := <-signalCh:
				func() {
					lock.RLock()
					defer lock.RUnlock()

					if h, ok := signals[sig]; ok {
						h(sig)
					}
				}()
			}
		}
	}

	stopF := func(error) {
		close(signalCh)
		cancel()
	}

	router := &Router{
		signalCh:   signalCh,
		signals:    signals,
		ignSignals: make(map[os.Signal]struct{}),
		lock:       lock,
	}

	return router, startF, stopF
}

func (s *Router) Handle(sig os.Signal, h Handler) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.signals[sig] = h
	signal.Notify(s.signalCh, sig)
	delete(s.ignSignals, sig)
}

func (s *Router) Reset(sig os.Signal) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.signals, sig)
	signal.Reset(sig)
}

func (s *Router) Ignore(sig os.Signal) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.signals, sig)
	signal.Ignore(sig)

	s.ignSignals[sig] = struct{}{}
}

func (s *Router) IsHandled(sig os.Signal) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	_, ok := s.signals[sig]

	return ok
}

func (s *Router) IsIgnored(sig os.Signal) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	_, ok := s.ignSignals[sig]

	return ok
}

// Fire a signal.
func (s *Router) Fire(sig os.Signal) {
	s.signalCh <- sig
}
