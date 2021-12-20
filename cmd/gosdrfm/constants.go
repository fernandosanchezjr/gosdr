package main

const (
	defaultSampleRate = 24000
	defaultBufLen     = 16384
	maximumOversample = 16
	maximumBufLen     = (maximumOversample * defaultBufLen)
	minimumRate       = 1000000
)
