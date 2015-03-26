package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestFlagParsing(t *testing.T) {
	var stdin *os.File
	stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{Stdin: stdin, Stdout: stdout, Stderr: stderr}
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
