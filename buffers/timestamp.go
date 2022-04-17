package buffers

import (
	"github.com/fernandosanchezjr/gosdr/utils"
	"strconv"
	"strings"
	"sync/atomic"
)

type Timestamp struct {
	digits []uint64
	last   int
}

func NewTimestamp() *Timestamp {
	return &Timestamp{
		digits: []uint64{0},
		last:   0,
	}
}

func (ts *Timestamp) String() string {
	strDigits := make([]string, len(ts.digits))
	for pos, value := range ts.digits {
		strDigits[pos] = strconv.FormatUint(value, 16)
	}
	return strings.Join(strDigits, ".")
}

func TimestampFromString(str string) (*Timestamp, error) {
	rawDigits := strings.Split(str, ".")
	digits := make([]uint64, len(rawDigits))
	for pos, raw := range rawDigits {
		if digit, err := strconv.ParseUint(raw, 16, 64); err == nil {
			digits[pos] = digit
		} else {
			return nil, err
		}
	}
	return &Timestamp{digits: digits, last: len(digits) - 1}, nil
}

func (ts *Timestamp) Child() *Timestamp {
	childDigits := make([]uint64, len(ts.digits)+1)
	copy(childDigits, ts.digits)
	return &Timestamp{digits: childDigits, last: ts.last + 1}
}

func (ts *Timestamp) Clone() *Timestamp {
	childDigits := make([]uint64, len(ts.digits))
	copy(childDigits, ts.digits)
	return &Timestamp{digits: childDigits, last: ts.last}
}

func (ts *Timestamp) Increment() {
	atomic.AddUint64(&ts.digits[ts.last], 1)
}

func (ts *Timestamp) Reset() {
	ts.digits[ts.last] = 0
}

func (ts *Timestamp) Copy(other *Timestamp) {
	if other.last < ts.last {
		other.digits = make([]uint64, len(ts.digits))
		other.last = ts.last
	}
	other.Reset()
	copy(other.digits, ts.digits)
}

func (ts *Timestamp) Less(other *Timestamp) bool {
	last := utils.MinInt(ts.last, other.last) + 1
	for i := 0; i < last; i++ {
		a, b := ts.digits[i], other.digits[i]
		if a < b {
			return true
		} else if a == b {
			continue
		} else {
			return false
		}
	}

	return ts.last <= other.last
}
