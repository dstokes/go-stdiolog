package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"
)

// flags
var (
	outfile = flag.String("o", "", "logfile for stdout")
	errfile = flag.String("3", "", "logfile for stderr")
)

// hooks for testing
var (
	outlog io.Writer
	errlog io.Writer
)

func main() {
	flag.Parse()
	args := flag.Args()

	// open output files
	outlog, err := os.OpenFile(*outfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer outlog.Close()

	if *errfile == "" {
		errlog = outlog
	} else {
		errlog, err := os.OpenFile(*errfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer errlog.Close()
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmdout, _ := cmd.StdoutPipe()
	cmderr, _ := cmd.StderrPipe()

	go r(cmdout, &Logger{os.Stdout, outlog})
	go r(cmderr, &Logger{os.Stderr, errlog})

	cmd.Start()
	cmd.Wait()
}

type Logger struct {
	Stdio   io.Writer
	Logfile io.Writer
}

func (l *Logger) Write(p []byte) (n int, err error) {
	now := time.Now().Format(time.RFC3339)
	// don't block writing to stdio
	go fmt.Fprintf(l.Logfile, "%s %s", now, p)
	return fmt.Fprintf(l.Stdio, "%s", p)
}

func r(reader io.Reader, logger *Logger) {
	s := bufio.NewReader(reader)
	for {
		if l, _, err := s.ReadLine(); err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		} else {
			logger.Write(append(l, 0x0A))
		}
	}
}
