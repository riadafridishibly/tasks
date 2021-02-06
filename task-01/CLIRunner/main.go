package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/gocarina/gocsv"
	"github.com/riadafridishibly/tasks/task-01/record"
)

type myFileHandler struct {
	sync.Mutex
	f *os.File
}

func (mf *myFileHandler) Write(data []byte) (int, error) {
	mf.Lock()
	defer mf.Unlock()

	return mf.f.Write(data)
}

func spwanProcess(wg *sync.WaitGroup, mf *myFileHandler, cmd string, args ...string) {
	r, w := io.Pipe()
	shCmd := exec.Command(cmd, args...)
	shCmd.Stdout = w

	teeReader := io.TeeReader(r, mf)

	go func() {
		io.Copy(os.Stdout, teeReader)
	}()

	if err := shCmd.Start(); err != nil {
		log.Fatal(err)
	}

	if err := shCmd.Wait(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Argument Required\n")
		os.Exit(1)
	}
	args := os.Args[1]
	args = strings.ReplaceAll(args, `\n`, "\n")
	// "Run,Title,Message 1,Message 2,Stream Delay,Run Times\n2,First Streamer,First Message,Second Msg,2,10\n2,Second Streamer,First Message,Second Msg,2,10"

	var cliRunners []record.CliRunnerRecord
	gocsv.UnmarshalString(args, &cliRunners)

	f, err := os.Create("sample.txt")
	if err != nil {
		log.Fatal(err)
	}
	mf := myFileHandler{
		f: f,
	}

	wg := sync.WaitGroup{}
	for _, clirunner := range cliRunners {
		shArgs := clirunner.CliStreamerRecordCSV()
		wg.Add(1)
		go spwanProcess(&wg, &mf, "../CLIStreamer/CLIStreamer", shArgs)
	}
	wg.Wait()
}
