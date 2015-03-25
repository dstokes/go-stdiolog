package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"
)

var (
	// flags
	outfile = flag.String("o", "", "logfile for stdout")
	errfile = flag.String("e", "", "logfile for stderr")

	outlog io.WriteCloser
	errlog io.WriteCloser
	wg     sync.WaitGroup
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
		errlog, err = os.OpenFile(*errfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer errlog.Close()
	}

	// setup child process
	cmd := exec.Command(args[0], args[1:]...)
	cmdout, _ := cmd.StdoutPipe()
	cmderr, _ := cmd.StderrPipe()

	// check for incoming data on stdin
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		wg.Add(1)
		stdin, _ := cmd.StdinPipe()
		go func() {
			// copy stdin to child process then close stdin
			if _, err := io.Copy(stdin, os.Stdin); err == nil {
				stdin.Close()
				wg.Done()
			} else {
				panic(err)
			}
		}()
	}

	wg.Add(2)
	go ReadWriteLine(cmdout, &Logger{os.Stdout, outlog})
	go ReadWriteLine(cmderr, &Logger{os.Stderr, errlog})

	cmd.Start()
	wg.Wait()
}

type Logger struct {
	Stdio   io.WriteCloser
	Logfile io.WriteCloser
}

func (l *Logger) Write(p []byte) (n int, err error) {
	now := time.Now().Format(time.RFC3339)
	// don't block writing to stdio
	go fmt.Fprintf(l.Logfile, "%s %s", now, p)
	return fmt.Fprintf(l.Stdio, "%s", p)
}

func (l *Logger) Close() error {
	return l.Logfile.Close()
}

func ReadWriteLine(reader io.ReadCloser, writer io.WriteCloser) {
	s := bufio.NewReader(reader)
	for {
		if l, _, err := s.ReadLine(); err == io.EOF {
			if len(l) > 0 {
				writer.Write(l)
			}
			break
		} else if err != nil {
			panic(err)
		} else {
			writer.Write(append(l, 0x0A))
		}
	}
	wg.Done()
}
