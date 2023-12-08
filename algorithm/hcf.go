package algorithm

import "math"

func Hcf(nums ...int) int {
	var r int
	var highest int = math.MaxInt

	for i := 0; (i + 1) < len(nums); i++ {
		a := nums[0]
		b := nums[i+1]
		for {
			if a%b == 0 {
				break
			}
			r = a % b
			a = b
			b = r
		}

		if b < highest {
			highest = b
		}
	}

	return highest

}
