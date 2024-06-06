package bench

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

var csvHeader = []string{
	"Who",
	"NumMsgs",
	"DroppedMsgs",
	"Reqs",
	"ReqsSuccesful",
	"ReqsFailed",
	"LatencyAvg",
	"LatencyMin",
	"LatencyMax",
	"TotalDuration",
	"Throughput",
}

var csvHeader2 = []string{
	"Bucket",
	"Number",
}

func WriteToCsv(name, benchname string, records []Result, clientRecord ClientResult) error {
	path := fmt.Sprintf("./csv/stats/%s.csv", benchname)
	fmt.Println("writing to csv...")
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	w := csv.NewWriter(file)
	data := make([][]string, len(records)+1)
	data[0] = csvHeader
	for i, record := range records {
		row := []string{
			record.Id,
			strconv.Itoa(int(record.TotalNum)),
			strconv.Itoa(int(record.Dropped)),
			strconv.Itoa(int(record.FinishedReqsTotal)),
			strconv.Itoa(int(record.FinishedReqsSuccesful)),
			strconv.Itoa(int(record.FinishedReqsFailed)),
			strconv.Itoa(int(record.RoundTripLatencyAvg.Microseconds())),
			strconv.Itoa(int(record.RoundTripLatencyMin.Microseconds())),
			strconv.Itoa(int(record.RoundTripLatencyMax.Microseconds())),
			"",
			"",
		}
		data[i+1] = row
	}
	clientRow := []string{
		clientRecord.Id,
		"",
		"",
		strconv.Itoa(int(clientRecord.Total)),
		strconv.Itoa(int(clientRecord.Successes)),
		strconv.Itoa(int(clientRecord.Errs)),
		strconv.Itoa(int(clientRecord.LatencyAvg.Microseconds())),
		strconv.Itoa(int(clientRecord.LatencyMin.Microseconds())),
		strconv.Itoa(int(clientRecord.LatencyMax.Microseconds())),
		strconv.Itoa(int(clientRecord.TotalDur.Microseconds())),
		fmt.Sprintf("%.2f", clientRecord.Throughput),
	}
	data = append(data, clientRow)
	err = w.WriteAll(data)
	if err != nil {
		panic(err)
	}

	path2 := fmt.Sprintf("./csv/histogram/hist.%s.csv", benchname)
	f, err := os.Create(path2)
	if err != nil {
		return err
	}
	defer file.Close()
	cw := csv.NewWriter(f)
	d := make([][]string, len(clientRecord.ReqDistribution)+1)
	d[0] = csvHeader2
	for i, r := range clientRecord.ReqDistribution {
		x := int(clientRecord.LatencyMin.Milliseconds()) + i*clientRecord.BucketSize
		d[i+1] = []string{strconv.Itoa(x), strconv.Itoa(int(r))}
	}
	return cw.WriteAll(d)
}

func WriteThroughputVsLatency(name string, throughputVsLatency [][]string) error {
	path := fmt.Sprintf("./csv/%s.csv", name)
	fmt.Println("writing throughput vs latency file...")
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	w := csv.NewWriter(file)
	data := make([][]string, 1, len(throughputVsLatency)+1)
	data[0] = []string{"Throughput", "Latency (avg)", "Latency (med)"}
	data = append(data, throughputVsLatency...)
	return w.WriteAll(data)
}

func WriteDurations(name string, durations []time.Duration) error {
	path := fmt.Sprintf("./csv/%s.csv", name)
	fmt.Println("writing durations file...")
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	w := csv.NewWriter(file)
	data := make([][]string, 1, len(durations)+1)
	data[0] = []string{"Number", "Latency (µs)"}
	for number, latency := range durations {
		data = append(data, []string{strconv.Itoa(number), strconv.Itoa(int(latency.Microseconds()))})
	}
	return w.WriteAll(data)
}

func WritePerformance(name string, results []string) error {
	path := fmt.Sprintf("./csv/%s.csv", name)
	fmt.Println("writing performance file...")
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	w := csv.NewWriter(file)
	data := make([][]string, 2)
	data[0] = []string{"Reqs/client", "Mean (µs)", "Median (µs)", "Std. dev.", "Min (µs)", "Max (µs)"}
	data[1] = results
	return w.WriteAll(data)
}
