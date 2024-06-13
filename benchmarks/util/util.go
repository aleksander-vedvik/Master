package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	bench "github.com/aleksander-vedvik/benchmark/benchmark"
)

type mappingType map[int]string

var mapping mappingType = map[int]string{
	1: bench.PaxosBroadcastCall,
	2: bench.PaxosQuorumCall,
	3: bench.PaxosQuorumCallBroadcastOption,
	4: bench.PBFTWithGorums,
	5: bench.PBFTWithoutGorums,
}

func (m mappingType) String() string {
	ret := "\n"
	ret += "\t1: " + m[1] + "\n"
	ret += "\t2: " + m[2] + "\n"
	ret += "\t3: " + m[3] + "\n"
	ret += "\t4: " + m[4] + "\n"
	ret += "\t5: " + m[5] + "\n"
	return ret
}

func main() {
	benchTypesString := flag.String("bench", "0", "types of benchmark to create histograms for. Should be comma seperated list:"+mapping.String())
	numBuckets := flag.Int("num", 0, "number of buckets in histogram")
	folderPath := flag.String("path", "./csv", "folder path to csv files")
	throughput := flag.Int("t", 0, "throughput")
	flag.Parse()

	benchTypes := make([]string, 0)
	if *benchTypesString != "" {
		benchTypesSlice := strings.Split(*benchTypesString, ",")
		for _, benchTypeString := range benchTypesSlice {
			index, err := strconv.Atoi(benchTypeString)
			if err != nil {
				panic(err)
			}
			benchType, ok := mapping[index]
			if !ok {
				panic("invalid bench type")
			}
			benchTypes = append(benchTypes, benchType)
		}
	}
	if len(benchTypes) == 0 {
		calculatePerformance(*folderPath)
		return
	}
	if *throughput == 0 && *numBuckets == 0 {
		calculateTvsL(benchTypes, *folderPath)
		return
	}
	createHistogram(*numBuckets, benchTypes, *folderPath, *throughput)
}

func createHistogram(numbuckets int, benchTypes []string, path string, throughput int) {
	min, max := -1, -1
	allLatencies := make([][]int, 0)
	for _, name := range benchTypes {
		latencies, err := getLatencies(name, path, throughput)
		if err != nil {
			panic(err)
		}
		if latencies[0] < min || min == -1 {
			min = latencies[0]
		}
		if latencies[len(latencies)-1] > max || max == -1 {
			max = latencies[len(latencies)-1]
		}
		allLatencies = append(allLatencies, latencies)
	}
	bucketSize := (max - min) / numbuckets
	for i, latencies := range allLatencies {
		histogram := make(map[int]int, numbuckets)
		for _, latency := range latencies {
			index := (latency - min) / bucketSize
			if index < 0 {
				index = 0
			}
			if index >= numbuckets {
				index = numbuckets - 1
			}
			histogram[index*bucketSize]++
		}
		for bucket, number := range histogram {
			histogram[bucket] = 100 * number / len(latencies)
		}
		WriteHistogram(benchTypes[i], histogram, path, throughput)
	}
}

func getLatencies(name, path string, throughput int) ([]int, error) {
	allDurs := make([][]int, 0, 1000)
	for runNumber := 0; runNumber < 20; runNumber++ {
		filePath := fmt.Sprintf("%s/%s.T%v.R%v.durations.csv", path, name, throughput, runNumber)
		fmt.Println("reading file:", filePath)
		file, err := os.Open(filePath)
		if err != nil {
			break
		}
		defer file.Close()
		reader := csv.NewReader(file)
		records, err := reader.ReadAll()
		if err != nil {
			return nil, err
		}
		// only accepts sorted records. This is done by the benchmark.
		for i, record := range records {
			if i == 0 {
				// header
				continue
			}
			// record[1] = Latency (µs)
			latency, err := strconv.Atoi(record[1])
			if err != nil {
				return nil, err
			}
			if i >= len(allDurs) {
				allDurs = append(allDurs, []int{latency})
				continue
			}
			allDurs[i] = append(allDurs[i], latency)
		}
	}
	if len(allDurs) <= 0 {
		return nil, errors.New("no files read")
	}
	coalescedDurs := make([]int, len(allDurs))
	for i, durs := range allDurs {
		sum := 0
		for _, dur := range durs {
			sum += dur
		}
		mean := sum / len(durs)
		/*// sort slices to get median
		slices.Sort(durs)
		// precision is not super important. Round down to closest int.
		var median int
		medianIndex := len(durs) / 2 // integer division does floor operation
		if len(durs)%2 == 0 {
			median1 := durs[medianIndex-1]
			median2 := durs[medianIndex]
			median = (median1 + median2) / 2
		} else {
			median = durs[medianIndex]
		}
		coalescedDurs[i] = median*/
		coalescedDurs[i] = mean / 1000 // convert µs to ms
	}
	return coalescedDurs, nil
}

func WriteHistogram(name string, histogram map[int]int, folderPath string, throughput int) error {
	path := fmt.Sprintf("%s/%s.T%v.histogram.csv", folderPath, name, throughput)
	fmt.Println("writing histogram to:", path)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	w := csv.NewWriter(file)
	data := make([][]string, 1, len(histogram)+1)
	data[0] = []string{"Bucket", "Number"}
	for bucket, number := range histogram {
		data = append(data, []string{strconv.Itoa(bucket), strconv.Itoa(number)})
	}
	return w.WriteAll(data)
}

func calculateTvsL(benchTypes []string, path string) {
	for i, name := range benchTypes {
		latencies, throughputs, err := getTLLatencies(name, path)
		if err != nil {
			panic(err)
		}
		for j, latency := range latencies {
			fmt.Println(benchTypes[i], throughputs[j], latency)
		}
		WriteTvsL(benchTypes[i], latencies, throughputs, path)
	}
}

func getTLLatencies(name, path string) ([]int, []int, error) {
	throughputs := make([]int, 0)
	allDurs := make([][]int, 0, 20)
	for runNumber := 0; runNumber < 20; runNumber++ {
		filePath := fmt.Sprintf("%s/%s.R%v.csv", path, name, runNumber)
		fmt.Println("reading file:", filePath)
		file, err := os.Open(filePath)
		if err != nil {
			break
		}
		defer file.Close()
		reader := csv.NewReader(file)
		records, err := reader.ReadAll()
		if err != nil {
			return nil, nil, err
		}
		// only accepts sorted records. This is done by the benchmark.
		for i, record := range records {
			if i == 0 {
				// header
				continue
			}
			// record[0] = Throughput
			throughput, err := strconv.Atoi(record[0])
			if err != nil {
				return nil, nil, err
			}
			if i >= len(throughputs) {
				throughputs = append(throughputs, throughput)
			}
			// record[1] = Latency (ms)
			latency, err := strconv.Atoi(record[1])
			if err != nil {
				return nil, nil, err
			}
			if i > len(allDurs) {
				allDurs = append(allDurs, []int{latency})
				continue
			}
			allDurs[i-1] = append(allDurs[i-1], latency)
		}
	}
	if len(allDurs) <= 0 {
		return nil, nil, errors.New("no files read")
	}
	coalescedDurs := make([]int, len(allDurs))
	for i, durs := range allDurs {
		sum := 0
		for _, dur := range durs {
			sum += dur
		}
		mean := sum / len(durs)
		/*// sort slices to get median
		slices.Sort(durs)
		// precision is not super important. Round down to closest int.
		var median int
		medianIndex := len(durs) / 2 // integer division does floor operation
		if len(durs)%2 == 0 {
			median1 := durs[medianIndex-1]
			median2 := durs[medianIndex]
			median = (median1 + median2) / 2
		} else {
			median = durs[medianIndex]
		}
		coalescedDurs[i] = median*/
		coalescedDurs[i] = mean
	}
	return coalescedDurs, throughputs, nil
}

func WriteTvsL(name string, latencies []int, throughputs []int, folderPath string) error {
	path := fmt.Sprintf("%s/%s.TvsL.csv", folderPath, name)
	fmt.Println("writing throughput vs latency to:", path)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	w := csv.NewWriter(file)
	data := make([][]string, 1, len(latencies)+1)
	data[0] = []string{"Throughput", "Latency"}
	for i := range latencies {
		data = append(data, []string{strconv.Itoa(throughputs[i]), strconv.Itoa(latencies[i])})
	}
	return w.WriteAll(data)
}

func calculatePerformance(path string) {
	benchTypes := []string{
		"Paxos.BroadcastCall.S3",
		"Paxos.QuorumCallBroadcastOption.S3",
		"Paxos.QuorumCall.S3",
		"PBFT.With.Gorums.S4",
		"PBFT.Without.Gorums.S4",
	}
	for _, name := range benchTypes {
		err := getPerformanceLatencies(name, path)
		if err != nil {
			panic(err)
		}
	}
}

func getPerformanceLatencies(name, path string) error {
	avg := 0.0
	median := 0.0
	stdv := 0.0
	min := 0.0
	max := 0.0
	runNumber := 0
	for ; runNumber < 20; runNumber++ {
		filePath := fmt.Sprintf("%s/%s.C10.R1000.Performance.%v.csv", path, name, runNumber)
		fmt.Println("reading file:", filePath)
		file, err := os.Open(filePath)
		if err != nil {
			break
		}
		defer file.Close()
		reader := csv.NewReader(file)
		records, err := reader.ReadAll()
		if err != nil {
			return err
		}
		record := records[1]
		avgLatency, err := strconv.Atoi(record[1])
		if err != nil {
			return err
		}
		avg += float64(avgLatency)
		medianLatency, err := strconv.Atoi(record[2])
		if err != nil {
			return err
		}
		median += float64(medianLatency)
		stdvLatency, err := strconv.Atoi(record[3])
		if err != nil {
			return err
		}
		stdv += float64(stdvLatency)
		minLatency, err := strconv.Atoi(record[4])
		if err != nil {
			return err
		}
		min += float64(minLatency)
		maxLatency, err := strconv.Atoi(record[5])
		if err != nil {
			return err
		}
		max += float64(maxLatency)
	}
	avg /= float64(1000 * runNumber)
	median /= float64(1000 * runNumber)
	stdv /= float64(1000 * runNumber)
	max /= float64(1000 * runNumber)
	min /= float64(1000 * runNumber)
	res := "\n\n"
	res += "Benchmark " + name + "\n"
	res += fmt.Sprintf("\t- Mean: %.2f\n", avg)
	res += fmt.Sprintf("\t- Median: %.2f\n", median)
	res += fmt.Sprintf("\t- Stdv: %.2f\n", stdv)
	res += fmt.Sprintf("\t- Min: %.2f\n", min)
	res += fmt.Sprintf("\t- Max: %.2f\n", max)
	res += "\n"
	fmt.Println(res)
	return nil
}
