package main

import "sync"

type Pair struct {
	P Player
	B Boulder
}

// PairStack the stack of pairs
type PairStack struct {
	pairs []Pair
	lock  sync.RWMutex
}

// New creates a new ItemStack
func (s *PairStack) New() *PairStack {
	s.pairs = []Pair{}
	return s
}

// IsEmpty returns true for empty stack
func (s *PairStack) IsEmpty() bool {
	return len(s.pairs) == 0
}

// Push adds a Pair to the top of the stack
func (s *PairStack) Push(p Pair) {
	s.lock.Lock()
	s.pairs = append(s.pairs, p)
	s.lock.Unlock()
}

// Pop removes a Pair from the top of the stack
func (s *PairStack) Pop() *Pair {
	if s.IsEmpty() {
		return nil
	}
	s.lock.Lock()
	pair := s.pairs[len(s.pairs)-1]
	s.pairs = s.pairs[0 : len(s.pairs)-1]
	s.lock.Unlock()
	return &pair
}
