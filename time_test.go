package tempo

import (
	"encoding/binary"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeHeader(t *testing.T) {
	longForm := "Jan 2, 2006 at 3:04pm (UTC)"
	timeStamp, _ := time.Parse(longForm, "Feb 3, 2018 at 7:00pm (UTC)")
	ts := NewTimeSeries(timeStamp)

	expectedStartTime := timeStamp.Unix()
	assert.Equal(t, expectedStartTime, ts.TimeHeader)
}

func TestTimeEncodeFirstValue(t *testing.T) {
	startTime := time.Now()
	ts := NewTimeSeries(startTime)

	t1 := startTime.Add(time.Second * time.Duration(60))
	ts.timeEncode(t1.Unix())

	expected := []byte{60, 0}

	assert.Equal(t, expected, ts.Opened.Stream)
}

func TestTimeEncodeRange0(t *testing.T) {
	startTime := time.Now()
	ts := NewTimeSeries(startTime)

	t1 := startTime.Add(time.Second * time.Duration(60))
	ts.timeEncode(t1.Unix())
	ts.Opened.Length++

	t2 := t1.Add(time.Second * time.Duration(60))
	ts.timeEncode(t2.Unix())
	assert.Equal(t, int16(60), int16(binary.LittleEndian.Uint16(ts.Opened.Stream)))

	valueByte := ts.Opened.Stream[2]
	assert.Equal(t, int8(0), int8(valueByte))

	t3 := t2.Add(time.Second * time.Duration(60))
	ts.timeEncode(t3.Unix())
	valueByte = ts.Opened.Stream[3]
	assert.Equal(t, int8(0), int8(valueByte))
}

func TestTimeEncodeRange64(t *testing.T) {
	startTime := time.Now()
	ts := NewTimeSeries(startTime)

	t1 := startTime.Add(time.Second * time.Duration(62))
	ts.timeEncode(t1.Unix())
	ts.Opened.Length++

	t2 := t1.Add(time.Second * time.Duration(60))
	ts.timeEncode(t2.Unix())
	assert.Equal(t, int16(62), int16(binary.LittleEndian.Uint16(ts.Opened.Stream)))

	tagByte := ts.Opened.Stream[2]
	assert.Equal(t, range64, tagByte)

	valueByte := ts.Opened.Stream[3]
	assert.Equal(t, int8(-2), int8(valueByte))
}

func TestTimeEncodeRange256(t *testing.T) {
	startTime := time.Now()
	ts := NewTimeSeries(startTime)

	t1 := startTime.Add(time.Second * time.Duration(60))
	ts.timeEncode(t1.Unix())
	ts.Opened.Length++

	t2 := t1.Add(time.Second * time.Duration(200))
	ts.timeEncode(t2.Unix())

	assert.Equal(t, int16(60), int16(binary.LittleEndian.Uint16(ts.Opened.Stream)))
	assert.Equal(t, range256, ts.Opened.Stream[2])
	assert.Equal(t, int16(140), int16(binary.LittleEndian.Uint16(ts.Opened.Stream[3:])))
}

func TestTimeEncodeRange2048(t *testing.T) {
	startTime := time.Now()
	ts := NewTimeSeries(startTime)

	t1 := startTime.Add(time.Second * time.Duration(60))
	ts.timeEncode(t1.Unix())
	ts.Opened.Length++

	t2 := t1.Add(time.Second * time.Duration(2000))
	ts.timeEncode(t2.Unix())

	assert.Equal(t, int16(60), int16(binary.LittleEndian.Uint16(ts.Opened.Stream)))
	assert.Equal(t, range2048, ts.Opened.Stream[2])
	assert.Equal(t, int16(1940), int16(binary.LittleEndian.Uint16(ts.Opened.Stream[3:])))
}

func TestTimeEncodeNegativeRange64(t *testing.T) {
	startTime := time.Now()
	ts := NewTimeSeries(startTime)

	t1 := startTime.Add(time.Second * time.Duration(60))
	ts.timeEncode(t1.Unix())
	ts.Opened.Length++

	t2 := t1.Add(time.Second * time.Duration(40))
	ts.timeEncode(t2.Unix())

	assert.Equal(t, int16(60), int16(binary.LittleEndian.Uint16(ts.Opened.Stream)))
	assert.Equal(t, range64, ts.Opened.Stream[2])
	assert.Equal(t, int8(-20), int8(ts.Opened.Stream[3]))
}

func TestTimeEncodeNegativeRange256(t *testing.T) {
	startTime := time.Now()
	ts := NewTimeSeries(startTime)

	t1 := startTime.Add(time.Second * time.Duration(300))
	ts.timeEncode(t1.Unix())
	ts.Opened.Length++

	t2 := t1.Add(time.Second * time.Duration(100))
	ts.timeEncode(t2.Unix())

	assert.Equal(t, int16(300), int16(binary.LittleEndian.Uint16(ts.Opened.Stream)))
	assert.Equal(t, range256, ts.Opened.Stream[2])
	assert.Equal(t, int16(-200), int16(binary.LittleEndian.Uint16(ts.Opened.Stream[3:])))
}

func TestTimeEncodeNegativeRange2048(t *testing.T) {
	startTime := time.Now()
	ts := NewTimeSeries(startTime)

	t1 := startTime.Add(time.Second * time.Duration(600))
	ts.timeEncode(t1.Unix())
	ts.Opened.Length++

	t2 := t1.Add(time.Second * time.Duration(100))
	ts.timeEncode(t2.Unix())

	assert.Equal(t, int16(600), int16(binary.LittleEndian.Uint16(ts.Opened.Stream)))
	assert.Equal(t, range2048, ts.Opened.Stream[2])
	assert.Equal(t, int16(-500), int16(binary.LittleEndian.Uint16(ts.Opened.Stream[3:])))
}
