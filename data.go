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
func (ts *TimeSeries) encodeData(data uint64) {
	if ts.LatestData == 0 {
		dataBytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(dataBytes, data)

		ts.Opened = append(ts.Opened, dataBytes...)
		ts.LatestData = data

		return
	}

	xor := data ^ ts.LatestData

	if xor == 0 {
		ts.Opened = append(ts.Opened, 0)

		ts.LatestData = data
		ts.LatestDataXor = xor

		return
	}

	lz := bits.LeadingZeros64(xor)
	tz := bits.TrailingZeros64(xor)

	bInfo := byte(similarBit) | byte(lz) // The header, tells us the similar bit, the control bit, and the LZ
	mxor := xor >> uint(tz)

	if lz == bits.LeadingZeros64(ts.LatestDataXor) && tz == bits.TrailingZeros64(ts.LatestDataXor) {
		ts.Opened = append(ts.Opened, bInfo)
		putUint64(&ts.Opened, mxor)

		ts.LatestDataXor = xor
		ts.LatestData = data

		return
	}

	bInfo |= byte(controlBit)
	mxorLength := byte(64 - lz - tz)

	ts.Opened = append(ts.Opened, bInfo)
	ts.Opened = append(ts.Opened, mxorLength)
	putUint64(&ts.Opened, mxor)

	ts.LatestDataXor = xor
	ts.LatestData = data
}

func putUint64(buf *[]byte, val uint64) {
	byteLength := 8 - (bits.LeadingZeros64(val) >> 3)
	b := make([]byte, byteLength)

	for i := 0; i < byteLength; i++ {
		b[i] = byte(val >> (uint(i) << 3))
	}

	*buf = append(*buf, b...)
}
