package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
)

type LogFile struct {
	Content *bytes.Buffer
}

func (l *LogFile) Write(b []byte) (int, error) {
	l.Content.Write(b)
	return len(b), nil
}

func (l *LogFile) Close() error {
	return nil
}

func TestFlagParsing(t *testing.T) {
	var stdin io.Reader
	stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{Stdin: stdin, Stdout: stdout, Stderr: stderr, Piped: false}
	args := strings.Split("go-stdiolog -v", " ")

	status := cli.Run(args)
	if status != ExitCodeOK {
		t.Errorf("expected %d to eq %d", status, ExitCodeOK)
	}

	expected := fmt.Sprintf("%s\n", Version)
	if !strings.Contains(stdout.String(), expected) {
		t.Errorf("expected %q to eq %q", stdout.String(), expected)
	}
}

func TestOutputLogging(t *testing.T) {
	var (
		stdin    io.Reader
		output   bytes.Buffer
		expected string = "ohai\n"
	)

	stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{Stdin: stdin, Stdout: stdout, Stderr: stderr}
	args := strings.Split("go-stdiolog -- echo ohai", " ")

	// overload output log
	outlog = &LogFile{&output}

	status := cli.Run(args)
	if status != ExitCodeOK {
		t.Errorf("expected %d to eq %d", status, ExitCodeOK)
	}

	// io pass-through
	if !(stdout.String() == expected) {
		t.Errorf("expected %q to eq %q", stdout.String(), expected)
	}

	// logfilh
	if !strings.Contains(output.String(), expected) {
		t.Errorf("expected %q to eq %q", stdout.String(), expected)
	}
}
