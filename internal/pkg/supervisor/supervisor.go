package supervisor

import (
	"sync"

	"golang.org/x/sync/errgroup"
)

type Supervisor struct {
	mutex     sync.Mutex
	stop      []chan<- StopToken
	processes errgroup.Group
}

func New() *Supervisor {
	return &Supervisor{
		mutex:     sync.Mutex{},
		stop:      nil,
		processes: errgroup.Group{},
	}
}

func (s *Supervisor) Run(p Process) {
	stop := s.addStop()
	s.processes.Go(func() error {
		return p(stop)
	})
}

func (s *Supervisor) Wait() error {
	return s.processes.Wait()
}

func (s *Supervisor) Stop() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, s := range s.stop {
		close(s)
	}
	s.stop = nil
}

func (s *Supervisor) addStop() <-chan StopToken {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	stop := make(chan StopToken)
	s.stop = append(s.stop, stop)
	return stop
}
