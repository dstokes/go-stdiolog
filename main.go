package main

import "os"

func main() {
	cli := &CLI{Stdin: os.Stdin, Stdout: os.Stdout, Stderr: os.Stderr}
	os.Exit(cli.Run(os.Args))
}
