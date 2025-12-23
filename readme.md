# rmarzullo
A Go implementation of **Marzullo's algorithm** for computing a
quorum-consistent time interval across multiple sources.
This algorithm is commonly used in distributed systems for
clock synchronization and fault tolerance.
## Features
- Closed-interval semantics
- Quorum-based overlap selection
- Nanosecond-precision time
- Deterministic and safe
- No global state
## Usage
```go
m := rmarzullo.Marzullo{

MinOverlap: 3,
ToleranceNs: 10_000_000, // 10ms
}
interval, err := m.SmallestInterval([]rmarzullo.IntervalWithSource{
{SourceID: 1, Interval: rmarzullo.Interval{LowNs: -5, HighNs: 5}},
{SourceID: 2, Interval: rmarzullo.Interval{LowNs: -3, HighNs: 7}},
{SourceID: 3, Interval: rmarzullo.Interval{LowNs: -4, HighNs: 6}},
})
if err != nil {
panic(err)
}
```
## Credits
This implementation is inspired by TigerBeetle’s use of Marzullo’s algorithm
in distributed databases to sync servers clocks.