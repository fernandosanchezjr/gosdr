package utils

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"strings"
)

type Sps int

func (s Sps) String() string {
	return humanize.SI(float64(s), "sps")
}

func ParseSps(value string) (Sps, error) {
	var frequency, unit, err = humanize.ParseSI(value)
	if unit != "" && strings.ToLower(unit) != "sps" {
		return 0, fmt.Errorf("invalid unit: '%s'", unit)
	}
	return Sps(frequency), err
}
