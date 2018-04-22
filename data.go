package tempo

import (
	"encoding/binary"
	"math/bits"
)

// Data info header
const (
	similarBit uint64 = 1 << 7
	controlBit uint64 = 1 << 6
)

// EncodeData compresses and stores new data to a byte-stream of compressed data
func (ts *TimeSeries) dataEncode(data uint64) {
	if ts.Opened.Length == 0 {
		dataBytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(dataBytes, data)

		ts.Opened.Stream = append(ts.Opened.Stream, dataBytes...)
		ts.LatestData = data
		ts.Opened.Length++

		return
	}

	xor := data ^ ts.LatestData

	if xor == 0 {
		ts.Opened.Stream = append(ts.Opened.Stream, 0)

		ts.LatestData = data
		ts.LatestDataXor = xor
		ts.Opened.Length++

		return
	}

	lz := bits.LeadingZeros64(xor)
	tz := bits.TrailingZeros64(xor)

	bInfo := byte(similarBit) | byte(lz) // The header, tells us the similar bit, the control bit, and the LZ
	mxor := xor >> uint(tz)

	if lz == bits.LeadingZeros64(ts.LatestDataXor) && tz == bits.TrailingZeros64(ts.LatestDataXor) {
		ts.Opened.Stream = append(ts.Opened.Stream, bInfo)
		putUint64(&ts.Opened.Stream, mxor)

		ts.LatestDataXor = xor
		ts.LatestData = data
		ts.Opened.Length++

		return
	}

	bInfo |= byte(controlBit)
	mxorLength := byte(64 - lz - tz)

	ts.Opened.Stream = append(ts.Opened.Stream, bInfo)
	ts.Opened.Stream = append(ts.Opened.Stream, mxorLength)
	putUint64(&ts.Opened.Stream, mxor)

	ts.LatestDataXor = xor
	ts.LatestData = data
	ts.Opened.Length++
}

func putUint64(buf *[]byte, val uint64) {
	byteLength := 8 - (bits.LeadingZeros64(val) >> 3)
	b := make([]byte, byteLength)

	// Parse val byte-by-byte
	for i := 0; i < byteLength; i++ {
		b[i] = byte(val >> (uint(i) << 3))
	}

	*buf = append(*buf, b...)
}

func putInt64(buf *[]byte, val int64, byteLength int) {
	for i := 0; i < byteLength; i++ {
		*buf = append(*buf, byte(val>>(uint(i)<<3)))
	}
}
