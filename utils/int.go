package utils

import "math/bits"

func MinInt(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func MaxInt(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

// IsPow2 returns true if N is a perfect power of 2 (1, 2, 4, 8, ...) and false otherwise.
// Algorithm from: https://graphics.stanford.edu/~seander/bithacks.html#DetermineIfPowerOf2
func IsPow2(N int) bool {
	if N == 0 {
		return false
	}
	return (uint64(N) & uint64(N-1)) == 0
}

// NextPow2 returns the smallest power of 2 >= N.
func NextPow2(N int) int {
	if N == 0 {
		return 1
	}
	return 1 << uint64(bits.Len64(uint64(N-1)))
}
