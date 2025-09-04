package supervisor

import (
	"context"
	"sync"

	"github.com/mabrarov/mp/pkg/panicerr"
)

type Supervisor struct {
	ctx    context.Context
	cancel context.CancelFunc
	group  sync.WaitGroup
	mutex  sync.Mutex
	err    error
}

func New(ctx context.Context) *Supervisor {
	ctx, cancel := context.WithCancel(ctx)
	return &Supervisor{
		ctx:    ctx,
		cancel: cancel,
		group:  sync.WaitGroup{},
		mutex:  sync.Mutex{},
		err:    nil,
	}
}

func (s *Supervisor) Go(p Process) {
	s.group.Add(1)
	go func() {
		defer s.group.Done()
		defer func() {
			if r := recover(); r != nil {
				s.setError(panicerr.New(r))
				s.Stop()
			}
		}()
		s.setError(p(s.ctx))
	}()
}

func (s *Supervisor) Wait() error {
	s.group.Wait()
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.err
}

func (s *Supervisor) Stop() {
	s.cancel()
}

func (s *Supervisor) setError(err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.err == nil {
		s.err = err
	}
}
