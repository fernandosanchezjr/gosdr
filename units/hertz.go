package units

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"strings"
)

type Hertz float64

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
	return float64(h)
}
