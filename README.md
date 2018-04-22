# tempo
Go implementation of facebook's gorilla: an in-memory write-through cache for time-series data

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
