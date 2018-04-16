package tempo

import (
	"sync"
)

type TimeSeries struct {
	mutex  sync.Mutex
	Opened []byte
	Closed [][]byte

	LatestTime       uint64
	LatestData       uint64
	LatestDataXor    uint64
	SecondLatestTime uint64
}
