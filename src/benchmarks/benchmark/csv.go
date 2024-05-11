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
}

func WriteToCsv(path string, records []Result, clientRecord ClientResult) error {
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
	}
	data = append(data, clientRow)
	return w.WriteAll(data)
}
