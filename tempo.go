package tempo

type TimeSeries struct {
	Opened []byte
	Closed [][]byte

	LatestTime       uint64
	LatestData       uint64
	LatestDataXor    uint64
	SecondLatestTime uint64
}
