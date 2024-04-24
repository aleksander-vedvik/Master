package metrics

type Metrics struct {
	Name       string
	ServerAddr string
	//mallocsPerOp int
	//throughput   int
	//latencyTotal int
	//latencyPerOp int
	TotalOps   uint64
	TotalTime  int64
	Throughput float64
	LatencyAvg float64
	LatencyVar float64
	// microbenchmark?
	//AllocsPerOp uint64
	//MemPerOp    uint64
	//ServerStats []struct {
	//Allocs uint64
	//Memory uint64
	//}
	//
	MsgStats struct {
		TotalNum         uint64
		Processed        uint64
		Discarded        uint64
		RoundTripLatency struct {
			Avg    uint64
			Median uint64
			Min    uint64
			Max    uint64
		}
		ReqLatency struct {
			Avg    uint64
			Median uint64
			Min    uint64
			Max    uint64
		}
		ShardDistribution map[int]int
		// measures unique number of broadcastIDs processed simultaneounsly
		ConcurrencyDistribution struct {
			Avg    uint64
			Median uint64
			Min    uint64
			Max    uint64
		}
	}
}

func New() *Metrics {
	return &Metrics{}
}
