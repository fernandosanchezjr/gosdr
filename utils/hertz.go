package utils

import (
	"github.com/dustin/go-humanize"
)

type Hertz int

func (h Hertz) String() string {
	return humanize.SI(float64(h), "Hz")
}
