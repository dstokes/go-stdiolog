package main

import "os"

func main() {
	stat, _ := os.Stdin.Stat()
	piped := (stat.Mode() & os.ModeCharDevice) == 0
	cli := &CLI{Stdin: os.Stdin, Stdout: os.Stdout, Stderr: os.Stderr, Piped: piped}
	os.Exit(cli.Run(os.Args))
}
