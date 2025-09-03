package supervisor

import (
	"errors"
	"sync"

	"mp/internal/pkg/panicerr"
)

type Supervisor struct {
	group   sync.WaitGroup
	mutex   sync.Mutex
	stopped bool
	stops   []chan<- StopToken
	err     error
}

func New() *Supervisor {
	return &Supervisor{
		group:   sync.WaitGroup{},
		mutex:   sync.Mutex{},
		stopped: false,
		stops:   nil,
		err:     nil,
	}
}

func (s *Supervisor) Run(p Process) {
	w := func(stop <-chan StopToken) (err error) {
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
		return p(stop)
	}
	stop := s.addStop()
	s.group.Add(1)
	go func() {
		defer s.group.Done()
		s.setError(w(stop))
	}()
}

func (s *Supervisor) Wait() error {
	s.group.Wait()
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.err
}

func (s *Supervisor) Stop() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, s := range s.stops {
		close(s)
	}
	s.stops = nil
	s.stopped = true
}

func (s *Supervisor) addStop() <-chan StopToken {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	stop := make(chan StopToken)
	if s.stopped {
		close(stop)
	} else {
		s.stops = append(s.stops, stop)
	}
	return stop
}

func (s *Supervisor) setError(err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.err == nil {
		s.err = err
	}
}
