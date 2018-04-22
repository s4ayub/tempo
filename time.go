package tempo

import (
	"encoding/binary"
)

const (
	range0     = 0
	range64    = 1 << 7
	range256   = 1 << 6
	range2048  = 1 << 5
	rangeOther = 1 << 4
)

// Maps time ranges to the number of bits used to store the time
var timeBytes = map[byte]int{
	range0:     0,
	range64:    1,
	range256:   2,
	range2048:  2,
	rangeOther: 4,
}

func (ts *TimeSeries) timeEncode(time uint64) {
	if ts.Opened.Length == 0 {
		dod := ts.TimeHeader - time
		dodBytes := make([]byte, 2)
		binary.LittleEndian.PutUint16(dodBytes, uint16(dod))

		ts.Opened.Stream = append(ts.Opened.Stream, dodBytes...)
		ts.SecondLatestTime = ts.LatestTime
		ts.LatestTime = time

		return
	}

	// No overflow because the lowest time denomination is ms, and we flush ~2hr
	dod := int64(time-ts.LatestTime) - int64(ts.LatestTime-ts.SecondLatestTime)

	if dod == 0 {
		ts.Opened.Stream = append(ts.Opened.Stream, 0)
	} else if (-63 <= dod) && (dod <= 64) {
		b := make([]byte, 1+timeBytes[range0])
		bs := b[1:]
		b[0] = range0
		putInt64(&bs, dod, timeBytes[range0])

		ts.Opened.Stream = append(ts.Opened.Stream, b...)
	} else if (-255 <= dod) && (dod <= 256) {
		b := make([]byte, 1+timeBytes[range64])
		bs := b[1:]
		b[0] = range64
		putInt64(&bs, dod, timeBytes[range64])

		ts.Opened.Stream = append(ts.Opened.Stream, b...)
	} else if (-2047 <= dod) && (dod <= 2048) {
		b := make([]byte, 1+timeBytes[range2048])
		bs := b[1:]
		b[0] = range2048
		putInt64(&bs, dod, timeBytes[range2048])

		ts.Opened.Stream = append(ts.Opened.Stream, b...)
	} else {
		b := make([]byte, 1+timeBytes[rangeOther])
		bs := b[1:]
		b[0] = rangeOther
		putInt64(&bs, dod, timeBytes[rangeOther])

		ts.Opened.Stream = append(ts.Opened.Stream, b...)
	}

	ts.SecondLatestTime = ts.LatestTime
	ts.LatestTime = time
}
