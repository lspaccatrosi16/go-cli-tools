package algorithm

import "fmt"

type orderfn = func(lastTerm int, backwards bool) int

type Sequence struct {
	evaluated map[int]int
	f         func(int, bool) int
	highestN  int
	highestV  int
	lowestN   int
	lowestV   int
}

func (s *Sequence) getForwards(i int) int {
	if val, has := s.evaluated[i]; has {
		return val
	} else {
		for j := s.highestN; j < i; j++ {
			val := s.f(s.highestV, false)
			s.highestV = val
			s.evaluated[j] = val
		}
		s.highestN = i
		return s.highestV
	}
}

func (s *Sequence) getBackwards(i int) int {
	if val, has := s.evaluated[i]; has {
		return val
	} else {
		for j := s.lowestN; j > i; j-- {
			val := s.f(s.lowestV, true)
			s.lowestV = val
			s.evaluated[j] = val
		}
		s.lowestN = i
		return s.lowestV
	}
}

func (s *Sequence) Get(i int) int {
	if i >= 1 {
		return s.getForwards(i)
	} else {
		return s.getBackwards(i)
	}
}

func seqSolveOrder(nums ...int) orderfn {
	if len(nums) < 2 {
		return func(int, bool) int {
			return 0
		}
	}

	firstTerm := nums[0]

	firstDiff := nums[1] - nums[0]
	lastDiff := nums[len(nums)-1] - nums[len(nums)-2]
	isSame := true

	for i := 0; i < len(nums); i++ {
		if nums[i] != firstTerm {
			isSame = false
			break
		}
	}

	if isSame {
		return func(_ int, backwards bool) int {
			return firstTerm
		}
	} else {
		differences := []int{}

		for i := 1; i < len(nums); i++ {
			differences = append(differences, nums[i]-nums[i-1])
		}
		diffFn := seqSolveOrder(differences...)

		fwDiff := lastDiff
		bwDiff := firstDiff

		return func(lastTerm int, backwards bool) int {
			fmt.Printf("enter %d %t\n", lastTerm, backwards)
			var diff, val int
			if backwards {
				diff = diffFn(bwDiff, backwards)
				bwDiff = diff
				val = lastTerm - diff
			} else {
				diff = diffFn(fwDiff, backwards)
				fwDiff = diff
				val = lastTerm + diff
			}

			fmt.Printf("eval %d, %d, %d %t\n", lastTerm, diff, val, backwards)
			return val
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
		lowestN:   1,
		lowestV:   nums[0],
	}

	return seq

}
