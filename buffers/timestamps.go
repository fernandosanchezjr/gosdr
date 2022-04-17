package buffers

type Timestamps []*Timestamp

func (ts Timestamps) Len() int {
	return len(ts)
}

func (ts Timestamps) Swap(i, j int) {
	ts[i], ts[j] = ts[j], ts[i]
}

func (ts Timestamps) Less(i, j int) bool {
	return ts[i].Less(ts[j])
}
