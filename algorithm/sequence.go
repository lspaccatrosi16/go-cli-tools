package algorithm

type orderfn = func(lastTerm int) int

type Sequence struct {
	evaluated map[int]int
	f         func(int) int
	highestN  int
	highestV  int
}

func (s *Sequence) Get(i int) int {
	if val, has := s.evaluated[i]; has {
		return val
	} else {
		for j := s.highestN; j < i; j++ {
			val := s.f(s.highestV)
			s.highestV = val
			s.evaluated[j] = val
		}
		s.highestN = i
		return s.highestV
	}
}

func seqSolveOrder(nums ...int) orderfn {
	if len(nums) < 2 {
		return func(int) int {
			return 0
		}
	}

	firstTerm := nums[0]
	lastDiff := nums[len(nums)-1] - nums[len(nums)-2]
	isSame := true

	for i := 0; i < len(nums); i++ {
		if nums[i] != firstTerm {
			isSame = false
			break
		}
	}

	if isSame {
		return func(int) int {
			return firstTerm
		}
	} else {
		differences := []int{}

		for i := 1; i < len(nums); i++ {
			differences = append(differences, nums[i]-nums[i-1])
		}
		diffFn := seqSolveOrder(differences...)

		prevDiff := lastDiff

		return func(lastTerm int) int {
			diff := diffFn(prevDiff)

			prevDiff = diff
			return lastTerm + diff
		}
	}
}

func SolveSequence(nums ...int) *Sequence {
	f := seqSolveOrder(nums...)

	evalMap := map[int]int{}

	var lastN, lastV int

	for i, n := range nums {
		evalMap[i+1] = n
		lastN = i + 1
		lastV = n
	}

	seq := &Sequence{
		f:         f,
		evaluated: evalMap,
		highestN:  lastN,
		highestV:  lastV,
	}

	return seq

}
