package bench

var benchTypes = map[string]struct {
	run func(benchmarkOption) ([]Result, error)
}{
	"Paxos": {
		run: func(bench benchmarkOption) ([]Result, error) {
			return runBenchmark(bench, PaxosBenchmark{})
		},
	},
	"PBFT": {
		run: func(bench benchmarkOption) ([]Result, error) {
			return runBenchmark(bench, PbftBenchmark{})
		},
	},
	"Simple": {
		run: func(bench benchmarkOption) ([]Result, error) {
			return runBenchmark(bench, SimpleBenchmark{})
		},
	},
}
