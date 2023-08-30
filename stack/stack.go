package stack

type Stack[T any] struct {
	stack []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{}
}

func (s *Stack[T]) Pop() (val T, ok bool) {
	if len(s.stack) == 0 {
		return *new(T), false
	}
	end := len(s.stack) - 1
	v := s.stack[end]
	s.stack = s.stack[:end]
	return v, true
}

func (s *Stack[T]) Push(val T) {
	s.stack = append(s.stack, val)
}

func (s *Stack[T]) Reverse() {
	end := len(s.stack)
	newStack := make([]T, end)
	for i, v := range s.stack {
		newStack[end-i-1] = v
	}
	s.stack = newStack
}

func (s *Stack[T]) Inner() []T {
	return s.stack
}
