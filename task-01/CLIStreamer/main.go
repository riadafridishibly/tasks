package main

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/gocarina/gocsv"
	"github.com/riadafridishibly/tasks/task-01/record"
)

func main() {
	if len(os.Args) < 2 {
		// TODO: Better error message
		fmt.Fprintf(os.Stderr, "Argument required\n")
		os.Exit(1)
	}

	// Argument Format
	// Title,            Message 1,     Message 2,  Stream Delay, Run Times
	// CLI Invoker Name, First Message, Second Msg, 2,            10

	args := os.Args[1]
	args = strings.ReplaceAll(args, `\n`, "\n")

	var cliStreamers []record.CliStreamerRecord
	gocsv.UnmarshalString(args, &cliStreamers)

	wg := sync.WaitGroup{}

	for _, clistreamer := range cliStreamers {
		wg.Add(1)
		clis := clistreamer
		go clis.Run(&wg)
	}

	wg.Wait()
}
