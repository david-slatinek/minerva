package random

import "math/rand/v2"

func Int(min, max int) int {
	return rand.IntN(max+1-min) + min
}
