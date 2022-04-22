package utils

import "sync/atomic"

func SetFlag(flag *uint32) bool {
	return atomic.CompareAndSwapUint32(flag, 0, 1)
}

func UnsetFlag(flag *uint32) bool {
	return atomic.CompareAndSwapUint32(flag, 1, 0)
}

func GetFlag(flag *uint32) bool {
	return atomic.LoadUint32(flag) == 1
}
