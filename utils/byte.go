package utils

func ConvertByte[F float32 | float64](u byte) F {
	return (F(u) - 127.4) / 128
}
