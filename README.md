go-stdiolog
===========

Passively log the stdio for a process.

[![Build Status](https://travis-ci.org/dstokes/go-stdiolog.png)](https://travis-ci.org/dstokes/go-stdiolog)

usage
=====

Write stdio of child proccess to log files while preserving terminal stdio.

```shell
$ go-stdiolog -o stdout.log -e stderr.log -- foo -a bar
```

Pass stdin to child proccess.

```shell
$ echo ohai | go-stdiolog -o stdout.log -- cat
```

install
=======
```shell
$ go get github.com/dstokes/go-stdiolog
```

Make sure your `PATH` includes your `GOPATH` bin directory:

```shell
export PATH=$PATH:$GOPATH/bin
```
