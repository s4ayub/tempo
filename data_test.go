package tempo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeFirstValue(t *testing.T) {
	ts := &TimeSeries{}
	ts.dataEncode(12)

	expected := []byte{12, 0, 0, 0, 0, 0, 0, 0}

	assert.Equal(t, expected, ts.Opened.Stream)
	assert.Equal(t, uint64(12), ts.LatestData)
	assert.Equal(t, uint64(0), ts.LatestDataXor)
}

func TestEncodeSameValues(t *testing.T) {
	ts := &TimeSeries{}
	ts.dataEncode(12)
	ts.dataEncode(12)

	expected := []byte{12, 0, 0, 0, 0, 0, 0, 0, 0}

	assert.Equal(t, expected, ts.Opened.Stream)
	assert.Equal(t, uint64(12), ts.LatestData)
	assert.Equal(t, uint64(0), ts.LatestDataXor)
}

func TestEncodeDifferentValues(t *testing.T) {
	ts := &TimeSeries{}
	ts.dataEncode(12)
	ts.dataEncode(24)

	expected := []byte{12, 0, 0, 0, 0, 0, 0, 0, 251, 3, 5}

	assert.Equal(t, expected, ts.Opened.Stream)
	assert.Equal(t, uint64(24), ts.LatestData)
	assert.Equal(t, uint64(0x14), ts.LatestDataXor)
}

func TestEncodeDifferentValuesSameMeaningfulXORLength(t *testing.T) {
	ts := &TimeSeries{}
	ts.dataEncode(12)
	ts.dataEncode(24)
	ts.dataEncode(12)

	expected := []byte{12, 0, 0, 0, 0, 0, 0, 0, 251, 3, 5, 187, 5}

	assert.Equal(t, expected, ts.Opened.Stream)
	assert.Equal(t, uint64(12), ts.LatestData)
	assert.Equal(t, uint64(0x14), ts.LatestDataXor)
}
