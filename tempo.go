package tempo

import "time"

type TimeSeries struct {
	Opened Block
	Closed []Block

	TimeHeader       int64
	LatestTime       int64
	LatestData       uint64
	LatestDataXor    uint64
	SecondLatestTime int64
}

type Block struct {
	Length int
	Stream []byte
}

// NewTimeSeries ....
func NewTimeSeries(start time.Time) *TimeSeries {
	ts := &TimeSeries{}

	// epoch ms
	ts.TimeHeader = start.Unix()

	return ts
}
