package tempo

import (
	"encoding/binary"
)

const (
	range0     byte = 0
	range64    byte = 1 << 7
	range256   byte = 1 << 6
	range2048  byte = 1 << 5
	rangeOther byte = 1 << 4
)

// Maps time ranges to the number of bits used to store the time
var timeBytes = map[byte]int{
	range0:     0,
	range64:    1, // when decoding, cast time as int8
	range256:   2, // when decoding, cast time as int16
	range2048:  2, // when decoding, cast time as int16
	rangeOther: 4, // when decoding, cast time as int32
}

func (ts *TimeSeries) timeEncode(time int64) {
	if ts.Opened.Length == 0 {
		dod := time - ts.TimeHeader
		dodBytes := make([]byte, 2)
		binary.LittleEndian.PutUint16(dodBytes, uint16(dod))

		ts.Opened.Stream = append(ts.Opened.Stream, dodBytes...)
		ts.SecondLatestTime = ts.TimeHeader
		ts.LatestTime = time

		return
	}

	dod := (time - ts.LatestTime) - (ts.LatestTime - ts.SecondLatestTime)

	if dod == 0 {
		ts.Opened.Stream = append(ts.Opened.Stream, 0)
	} else if (-63 <= dod) && (dod <= 64) {
		tagAndValue := make([]byte, 1)
		tagAndValue[0] = range64
		putInt64(&tagAndValue, dod, timeBytes[range64])

		ts.Opened.Stream = append(ts.Opened.Stream, tagAndValue...)

	} else if (-255 <= dod) && (dod <= 256) {
		tagAndValue := make([]byte, 1)
		tagAndValue[0] = range256
		putInt64(&tagAndValue, dod, timeBytes[range256])

		ts.Opened.Stream = append(ts.Opened.Stream, tagAndValue...)

	} else if (-2047 <= dod) && (dod <= 2048) {
		tagAndValue := make([]byte, 1)
		tagAndValue[0] = range2048
		putInt64(&tagAndValue, dod, timeBytes[range2048])

		ts.Opened.Stream = append(ts.Opened.Stream, tagAndValue...)
	} else {
		tagAndValue := make([]byte, 1)
		tagAndValue[0] = rangeOther
		putInt64(&tagAndValue, dod, timeBytes[rangeOther])

		ts.Opened.Stream = append(ts.Opened.Stream, tagAndValue...)
	}

	ts.SecondLatestTime = ts.LatestTime
	ts.LatestTime = time
}
