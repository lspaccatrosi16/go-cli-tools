package set

var exists = struct{}{}

type mSet[T comparable] map[T]struct{}

type Set[T comparable] struct {
	m mSet[T]
}

func (s *Set[T]) Add(vals ...T) {
	for _, val := range vals {
		s.m[val] = exists
	}
}

func (s *Set[T]) Contains(k T) bool {
	_, ok := s.m[k]
	return ok
}

func (s *Set[T]) Remove(k T) {
	delete(s.m, k)
}

func (s *Set[T]) GetIterator() func() (T, bool) {
	var valArr = []T{}
	for v := range s.m {
		valArr = append(valArr, v)
	}
	i := 0

	return func() (T, bool) {
		if i > len(valArr)-1 {
			var blank T
			return blank, false
		}
		val := valArr[i]
		i++
		return val, true
	}
}

func NewSet[T comparable]() *Set[T] {
	s := &Set[T]{}
	s.m = make(mSet[T])
	return s
}
