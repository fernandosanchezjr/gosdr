package buffers

import (
	log "github.com/sirupsen/logrus"
	"sort"
	"strings"
	"testing"
)

func TestNewTimestamp(t *testing.T) {
	ts := NewTimestamp()
	cts := ts.Child()
	for i := 0; i < 0xf; i++ {
		for j := 0; j < 0xf; j++ {
			cts.Increment()
			log.WithField("timestamp", cts.String()).Info("Incremented child")
		}
		ts.Increment()
		ts.Copy(cts)
	}
}

func TestTimestamp_Less(t *testing.T) {
	testCases := [][2]string{
		{"0", "0.1"},
		{"0", "0.0"},
		{"0", "1"},
		{"0", "1.0"},
		{"1", "1.0"},
		{"1.0", "1.1"},
		{"0.0", "1"},
	}
	for pos, test := range testCases {
		a, b := test[0], test[1]
		tsA, _ := TimestampFromString(a)
		tsB, _ := TimestampFromString(b)
		if !tsA.Less(tsB) {
			t.Fatal("Failed case", pos+1, a, "<", b)
		}
	}
}

func TestTimestamps_Sort(t *testing.T) {
	rawStamps := []string{
		"0.1",
		"1.1",
		"0",
		"1.0",
		"1",
		"0.0",
	}
	expectedStamps := []string{
		"0",
		"0.0",
		"0.1",
		"1",
		"1.0",
		"1.1",
	}
	stamps := make(Timestamps, len(rawStamps))
	for pos, test := range rawStamps {
		s, _ := TimestampFromString(test)
		stamps[pos] = s
	}
	sort.Stable(stamps)
	for pos, s := range stamps {
		if strings.Compare(s.String(), expectedStamps[pos]) != 0 {
			t.Fatal("Failed sort at index", pos, s, "vs", expectedStamps[pos])
		}
	}
}
