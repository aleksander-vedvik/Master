package bench

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
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
	path := fmt.Sprintf("./csv/stats/%s.%s.csv", name, benchname)
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
	totalDur := clientRecord.TotalDur.Seconds()
	if totalDur <= 0 {
		totalDur = 1
	}
	throughput := float64(clientRecord.Total) / totalDur
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
		fmt.Sprintf("%.2f", throughput),
	}
	data = append(data, clientRow)
	err = w.WriteAll(data)
	if err != nil {
		panic(err)
	}

	path2 := fmt.Sprintf("./csv/histogram/hist.%s.%s.csv", name, benchname)
	f, err := os.Create(path2)
	if err != nil {
		return err
	}
	defer file.Close()
	cw := csv.NewWriter(f)
	d := make([][]string, len(clientRecord.ReqDistribution)+1)
	d[0] = csvHeader2
	for i, r := range clientRecord.ReqDistribution {
		x := int(clientRecord.LatencyMin.Microseconds()) + i*clientRecord.BucketSize
		d[i+1] = []string{strconv.Itoa(x), strconv.Itoa(int(r))}
	}
	return cw.WriteAll(d)
}
