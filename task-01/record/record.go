package record

import (
	"fmt"
	"sync"
	"time"

	"github.com/gocarina/gocsv"
)

// CliStreamerRecord stores streamer info
type CliStreamerRecord struct {
	Title       string `csv:"Title"`
	Message1    string `csv:"Message 1"`
	Message2    string `csv:"Message 2"`
	StreamDelay int    `csv:"Stream Delay"`
	RunTimes    int    `csv:"Run Times"`
}

// CliStreamerRecordCSV runner record -> streamer record csv
func (cliRunnerRecord CliRunnerRecord) CliStreamerRecordCSV() string {
	cliStreamerRecords := []CliStreamerRecord{cliRunnerRecord.CliStreamerRecord()}

	out, err := gocsv.MarshalString(cliStreamerRecords)

	if err != nil {
		panic(err)
	}

	return out
}

// Run method runs a record asynchronously
func (csr *CliStreamerRecord) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < csr.RunTimes; i++ {
		fmt.Printf("%s->%s\n", csr.Title, csr.Message1)
		time.Sleep(time.Duration(csr.StreamDelay) * time.Second)
		fmt.Printf("%s->%s\n", csr.Title, csr.Message2)
	}
}

// CliRunnerRecord stores runner info
type CliRunnerRecord struct {
	Run         string `csv:"Run"`
	Title       string `csv:"Title"`
	Message1    string `csv:"Message 1"`
	Message2    string `csv:"Message 2"`
	StreamDelay int    `csv:"Stream Delay"`
	RunTimes    int    `csv:"Run Times"`
}

// CliStreamerRecord CliRunnerRecord -> CliStreamerRecord
func (cliRunnerRecord CliRunnerRecord) CliStreamerRecord() CliStreamerRecord {
	return CliStreamerRecord{
		Title:       cliRunnerRecord.Title,
		Message1:    cliRunnerRecord.Message1,
		Message2:    cliRunnerRecord.Message2,
		StreamDelay: cliRunnerRecord.StreamDelay,
		RunTimes:    cliRunnerRecord.RunTimes,
	}
}

// CSV takes a pointer to array of CLIRunnerRecord and return serilized string
func CSV(cliRunners *[]CliRunnerRecord) string {
	out, err := gocsv.MarshalString(cliRunners)

	if err != nil {
		panic(err)
	}

	return out
}
