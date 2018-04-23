# tempo
Go implementation of Facebook's GorillaDB: an in-memory write-through cache for time-series data

```Go
import (
  "time"
  "github.com/s4ayub/tempo"
)

func main() {
  // A time series is the underlying data structure behind the time series map.
  // A user of the package will never interact with the TimeSeries struct, but the
  // example code highlights our implementation of the delta-of-delta encoding for time and XOR encoding for values.

  startTime := time.Now()
  ts := NewTimeSeries(startTime)

  t1 := startTime.Add(time.Second * time.Duration(60)).Unix() // store in seconds
  ts.timeEncode(t1)
  ts.dataEncode(2)

  t2 := t1.Add(time.Second * time.Duration(60)).Unix()
  ts.timeEncode(t2)
  ts.dataEncode(3)
  
  fmt.Println(ts)
}
```

# To do:
- [x] Lay out the TimeSeries struct to encase timestamps and data in a single byte stream
- [x] Implement and test XOR encoding for data values
- [x] Implement and test delta-of-delta encoding for timestamps
- [ ] Write a decoder to return decoded blocks of time series data
- [ ] Build the TimeSeriesMap structure which will abstract away the TimeSeries struct

# Changes to the encoding scheme:
Tempo stores its timestamps and data values in a single byte-stream (for now).
The encoding schemes used by GorillaDB takes advantage of bit-level granularity.

Here, we discuss how we changed the algorithms to be byte-friendly, and the ramifications for doing so.

Timeseries encoding scheme:
---------------------------
1) The tag bits are always stored in 8 bits instead of 4.
    - The compression ratio for timestamps is reduced

2) The value bits associated with each tag is changed:
    - [0]           -> 0 bits (the same as before)
    - [-63, 64]     -> 8 bits
    - [-255, 256]   -> 16 bits
    - [-2047, 2048] -> 16 bits
    - [Otherwise]   -> 32 bits (the same as before)

    - Again, the compression ratio drops as unused bits are added.
    - [-255, 256] for example, requires only 9 bits but 16 bits are reserved.

3) Example of an encoded timestamp entry
    - {Tag bits}{Value}
        - 'Tag bits' is the tag defining the number of bits the time value will use
        - 'Value'    is the actual delta-of-delta time value
        - Each { } item is 1 byte, but '...' represents a variable amount of bytes

Data encoding scheme:
---------------------
1) The dissimilar bit, the control bit, and the leading zeroes are stored in the first byte.
    - [xxxx xxxx]
    - b[7] is the dissimilar bit
    - b[6] is the control bit
    - b[5:0] is the number of leading zeroes of the meaningful xor
        - LZ is encoded in GorillaDB, but it is not encoded in Tempo
        - These 6 bits would have been here regardless of whether we used it to store LZ

2) The control bit represents whether the length of the meaningful xor is reused
    - In GorillaDB, this is an indicator for whether the LZ and the MXOR length is reusued
    - In Tempo, this is an indicator for whether the MXOR length is reused
        - The LZ is always in the first byte

3) Example of an encoded data entry
    - {Header}{MXOR length}{MXOR...}
        - 'Header'      is the byte containing the dissimilar bit, control bit, and LZ
        - 'MXOR length' is the number of bits used to store the MXOR
        - 'MXOR'        is the actual meaningful xor bits
        - Each { } item is 1 byte, but '...' represents a variable amount of bytes
