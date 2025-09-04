package supervisor

import (
	"context"
	"errors"
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
		s.setError(runWithoutPanic(p, s.ctx))
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

func runWithoutPanic(p Process, ctx context.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			pe := panicerr.New(r)
			if err == nil {
				err = pe
			} else {
				err = errors.Join(err, pe)
			}
		}
	}()
	return p(ctx)
}

func (s *Supervisor) setError(err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.err == nil {
		s.err = err
	}
}
