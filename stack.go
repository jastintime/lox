package main

type Stack[T any] struct {
	list []T
}

func (s *Stack[T]) Push(element T) {
	s.list = append(s.list, element)
}

func (s *Stack[T]) Pop() T {
	last := s.list[len(s.list)-1]
	s.list = s.list[:len(s.list)-1]
	return last
}

func (s Stack[T]) IsEmpty() bool {
	return len(s.list) == 0
}

func (s Stack[T]) Peek() T {
	return s.list[len(s.list)-1]
}

func (s Stack[T]) Get(i int) T {
	return s.list[i]
}

func (s Stack[T]) Size() int {
	return len(s.list)
}
