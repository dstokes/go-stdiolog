go-stdiolog
===========

Passively log the stdio for a process.

usage
=====

```shell
$ go-stdiolog -o stdout.log -e stderr.log -- foo -a bar
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
