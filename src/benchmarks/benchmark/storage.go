package bench

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

var csvHeader = []string{
	"Who",
	"TotalNum",
	"GoroutinesStarted",
	"GoroutinesStopped",
	"FinishedReqsTotal",
	"FinishedReqsSuccesful",
	"FinishedReqsFailed",
	"Processed",
	"Dropped",
	"Invalid",
	"AlreadyProcessed",
	"RoundTripLatencyAvg",
	"RoundTripLatencyMin",
	"RoundTripLatencyMax",
	"ReqLatencyAvg",
	"ReqLatencyMin",
	"ReqLatencyMax",
}

func WriteToCsv(path string, records []Result) error {
	fmt.Println("writing to csv...")
	//records := []Employee{
	//{"E01", 25},
	//{"E02", 26},
	//{"E03", 24},
	//{"E04", 26},
	//}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	w := csv.NewWriter(file)
	//defer w.Flush() // Using Write
	//for _, record := range records {
	//row := []string{record.ID, strconv.Itoa(record.Age)}
	//if err := w.Write(row); err != nil {
	//log.Fatalln("error writing record to file", err)
	//}
	//}

	// Using WriteAll
	data := make([][]string, len(records)+1)
	data[0] = csvHeader
	for i, record := range records {
		row := []string{
			record.Id,
			strconv.Itoa(int(record.TotalNum)),
			strconv.Itoa(int(record.GoroutinesStarted)),
			strconv.Itoa(int(record.GoroutinesStopped)),
			strconv.Itoa(int(record.FinishedReqsTotal)),
			strconv.Itoa(int(record.FinishedReqsSuccesful)),
			strconv.Itoa(int(record.FinishedReqsFailed)),
			strconv.Itoa(int(record.Processed)),
			strconv.Itoa(int(record.Dropped)),
			strconv.Itoa(int(record.Invalid)),
			strconv.Itoa(int(record.AlreadyProcessed)),
			strconv.Itoa(int(record.RoundTripLatencyAvg)),
			strconv.Itoa(int(record.RoundTripLatencyMin)),
			strconv.Itoa(int(record.RoundTripLatencyMax)),
			strconv.Itoa(int(record.ReqLatencyAvg)),
			strconv.Itoa(int(record.ReqLatencyMin)),
			strconv.Itoa(int(record.ReqLatencyMax)),
		}
		data[i+1] = row
	}
	return w.WriteAll(data)
}
