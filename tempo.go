package tempo

type TimeSeries struct {
	Opened []byte
	Closed [][]byte

	LatestTime       []byte
	LatestValue      []byte
	SecondLatestTime []byte
}