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
	bInfo := byte(similarBit) | byte(lz)

	mxor := xor << uint(lz)
	mxor = xor >> uint(lz+tz)
	mxorBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(mxorBytes, mxor)

	// The number of meaningful bits are equal to the previously stored xor
	if lz == bits.LeadingZeros64(ts.LatestDataXor) && tz == bits.TrailingZeros64(ts.LatestDataXor) {
		ts.Opened = append(ts.Opened, bInfo)
		ts.Opened = append(ts.Opened, mxorBytes...)

		ts.LatestDataXor = xor
		ts.LatestData = data

		return
	}

	bInfo |= byte(controlBit)

	mxorLength := byte(64 - lz - tz)

	ts.Opened = append(ts.Opened, bInfo)
	ts.Opened = append(ts.Opened, mxorLength)
	ts.Opened = append(ts.Opened, mxorBytes...)

	ts.LatestDataXor = xor
	ts.LatestData = data
}
