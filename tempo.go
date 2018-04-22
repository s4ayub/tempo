package tempo

type TimeSeries struct {
	Opened Block
	Closed []Block

	TimeHeader       uint64
	LatestTime       uint64
	LatestData       uint64
	LatestDataXor    uint64
	SecondLatestTime uint64
}

type Block struct {
	Length int
	Stream []byte
}
