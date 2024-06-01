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
	6: bench.PBFTNoOrder,
	7: bench.PBFTOrder,
	8: bench.Simple,
}

func (m mappingType) String() string {
	ret := "\n"
	ret += "\t0: " + m[0] + "\n"
	ret += "\t1: " + m[1] + "\n"
	ret += "\t2: " + m[2] + "\n"
	ret += "\t3: " + m[3] + "\n"
	ret += "\t4: " + m[4] + "\n"
	ret += "\t5: " + m[5] + "\n"
	ret += "\t6: " + m[6] + "\n"
	ret += "\t7: " + m[7] + "\n"
	return ret
}

func main() {
	benchTypesString := flag.String("bench", "0", "types of benchmark to create histograms for. Should be comma seperated list:"+mapping.String())
	numBuckets := flag.Int("num", 0, "number of buckets in histogram")
	folderPath := flag.String("path", "./csv", "folder path to csv files")
	flag.Parse()

	benchTypesSlice := strings.Split(*benchTypesString, ",")
	benchTypes := make([]string, 0)
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
	createHistogram(*numBuckets, benchTypes, *folderPath)
}

func createHistogram(numbuckets int, benchTypes []string, path string) {
	min, max := -1, -1
	allLatencies := make([][]int, 0)
	for _, name := range benchTypes {
		latencies, err := getLatencies(name, path)
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
		WriteHistogram(benchTypes[i], histogram, path)
	}
}

func getLatencies(name string, path string) ([]int, error) {
	allDurs := make([][]int, 0, 1000)
	for runNumber := 0; runNumber < 20; runNumber++ {
		filePath := fmt.Sprintf("%s/%s.%v.durations.csv", path, name, runNumber)
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
			if i > len(allDurs) {
				allDurs = append(allDurs, []int{latency})
				continue
			}
			allDurs[i] = append(allDurs[i], latency)
		}
	}
	if len(allDurs) <= 0 {
		return nil, errors.New("no files read")
	}
	avgDurs := make([]int, len(allDurs))
	for i, durs := range allDurs {
		sum := 0
		for _, dur := range durs {
			sum += dur
		}
		// precision is not super important. Round down to closest int.
		avgDurs[i] = sum / len(durs)
	}
	return avgDurs, nil
}

func WriteHistogram(name string, histogram map[int]int, folderPath string) error {
	path := fmt.Sprintf("%s/%s.histogram.csv", folderPath, name)
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
