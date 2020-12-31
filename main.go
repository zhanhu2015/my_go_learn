package main

import "sync"

type Stats struct {
	mutex sync.Mutex

	counters map[string]int
}

func (s *Stats) Snapshot() map[string]int {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.counters
}

func (s *Stats) Add(name string, num int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.counters[name] = num
}
