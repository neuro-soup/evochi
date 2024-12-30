package stack

import (
	"fmt"
	"sync"
)

type Stack[T any] struct {
	mu  *sync.RWMutex
	s   []T
	top int
}

func New[T any](values ...T) *Stack[T] {
	return &Stack[T]{
		mu:  new(sync.RWMutex),
		s:   values,
		top: len(values),
	}
}

func (s *Stack[T]) String() string {
	return fmt.Sprintf("stack%v", fmt.Sprint(s.s[:s.top]))
}

func (s *Stack[T]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.top
}

func (s *Stack[T]) Push(v T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.push(v)
}

func (s *Stack[T]) push(v T) {
	if s.top == len(s.s) {
		s.s = append(s.s, v)
	} else {
		s.s[s.top] = v
	}

	s.top++
}

func (s *Stack[T]) PushAll(values ...T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, v := range values {
		s.push(v)
	}
}

func (s *Stack[T]) Pop() T {
	if s.top == 0 {
		panic("stack is empty")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.top--
	return s.s[s.top]
}

func (s *Stack[T]) Peek() T {
	if s.top == 0 {
		panic("stack is empty")
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.s[s.top-1]
}

func (s *Stack[T]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.top = 0
}
