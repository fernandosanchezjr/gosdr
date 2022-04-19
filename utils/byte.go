package utils

func ConvertByte(u byte) float32 {
	return (float32(u) - 127.4) / 128
}
