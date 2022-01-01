package units

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"strings"
)

type Hertz int

func (h Hertz) String() string {
	return humanize.SI(float64(h), "Hz")
}

func ParseHertz(value string) (Hertz, error) {
	var frequency, unit, err = humanize.ParseSI(value)
	if unit != "" && strings.ToLower(unit) != "hz" {
		return 0, fmt.Errorf("invalid unit: '%s'", unit)
	}
	return Hertz(frequency), err
}

func (h Hertz) Float() float64 {
	var value, _ = humanize.ComputeSI(float64(h))
	return value
}
