package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"time"
)

const (
	ExitCodeOK             int = 0
	ExitCodeError          int = 1
	ExitCodeFlagParseError     = 10 + iota
	ExitCodeFileOpenError
)

const HelpText string = `Usage: stdiolog [options] -- <command> [options]

Options:
  -h  Print this message and exit
  -o  The stdout log
  -e  The stderr log (defaults to -o)
`

var (
	outlog, errlog io.WriteCloser
	wg             sync.WaitGroup
)

type CLI struct {
	Piped          bool
	Stdin          io.Reader
	Stdout, Stderr io.Writer
}

// invoke the cli with args
func (cli *CLI) Run(args []string) int {
	var err error

	// parse args string
	flags := flag.NewFlagSet("cFlags", flag.ContinueOnError)
	flags.SetOutput(cli.Stdout)

	help := flags.Bool("h", false, "print help and quit")
	version := flags.Bool("v", false, "print version and exit")
	outfile := flags.String("o", "", "logfile for stdout")
	errfile := flags.String("e", "", "logfile for stderr")

	if err = flags.Parse(args[1:]); err != nil {
		return ExitCodeFlagParseError
	}

	if *version {
		fmt.Fprintf(cli.Stdout, "%s\n", Version)
		return ExitCodeOK
	}

	if *help {
		fmt.Fprintf(cli.Stderr, HelpText)
		return ExitCodeOK
	}

	childArgs := flags.Args()

	// open output files
	if outlog == nil {
		outlog, err = os.OpenFile(*outfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			fmt.Fprintf(cli.Stderr, err.Error())
			return ExitCodeFileOpenError
		}
		defer outlog.Close()
	}

	if *errfile == "" {
		errlog = outlog
	} else {
		errlog, err = os.OpenFile(*errfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			fmt.Fprintf(cli.Stderr, err.Error())
			return ExitCodeFileOpenError
		}
		defer errlog.Close()
	}

	// setup child process
	cmd := exec.Command(childArgs[0], childArgs[1:]...)
	cmdout, _ := cmd.StdoutPipe()
	cmderr, _ := cmd.StderrPipe()

	// check for incoming data on stdin
	if cli.Piped {
		wg.Add(1)
		stdin, _ := cmd.StdinPipe()
		go func() {
			// copy stdin to child process then close stdin
			if _, err := io.Copy(stdin, cli.Stdin); err == nil {
				stdin.Close()
				wg.Done()
			} else {
				panic(err)
			}
		}()
	}

	wg.Add(2)
	go ReadWriteLine(cmdout, &Logger{cli.Stdout, outlog})
	go ReadWriteLine(cmderr, &Logger{cli.Stderr, errlog})

	cmd.Start()
	wg.Wait()

	return ExitCodeOK
}

type Logger struct {
	Stdio   io.Writer
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

func ReadWriteLine(reader io.ReadCloser, writer io.Writer) {
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
