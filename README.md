# tsmock

[![Go Report Card](https://goreportcard.com/badge/github.com/thorstenrie/tsmock)](https://goreportcard.com/report/github.com/thorstenrie/tsmock)
[![CodeFactor](https://www.codefactor.io/repository/github/thorstenrie/tsmock/badge)](https://www.codefactor.io/repository/github/thorstenrie/tsmock)
![OSS Lifecycle](https://img.shields.io/osslifecycle/thorstenrie/tsmock)

[![PkgGoDev](https://pkg.go.dev/badge/mod/github.com/thorstenrie/tsmock)](https://pkg.go.dev/mod/github.com/thorstenrie/tsmock)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/thorstenrie/tsmock)
![Libraries.io dependency status for GitHub repo](https://img.shields.io/librariesio/github/thorstenrie/tsmock)

![GitHub release (latest by date)](https://img.shields.io/github/v/release/thorstenrie/tsmock)
![GitHub last commit](https://img.shields.io/github/last-commit/thorstenrie/tsmock)
![GitHub commit activity](https://img.shields.io/github/commit-activity/m/thorstenrie/tsmock)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/thorstenrie/tsmock)
![GitHub Top Language](https://img.shields.io/github/languages/top/thorstenrie/tsmock)
![GitHub](https://img.shields.io/github/license/thorstenrie/tsmock)

The Go package tsmock provides an interface to test and mock Stdin based on files. The package enables testing of interactive console applications. It reads input from a file and
passes it to `os.Stdin`. It can be configured to set the visibility of the input and a delay in processing each line of the
input from the file. The mocked Stdin is executed in a Go routine and can be canceled with a context.

- **Simple**: Without configuration, just function calls
- **Easy to use**: Retrieve Stdin input from [os.Stdin](https://pkg.go.dev/os)
- **Tested**: Unit tests with a high code coverage
- **Dependencies**: Only depends on the [Go Standard Library](https://pkg.go.dev/std) as well as [tsfio](https://github.com/thorstenrie/tsfio) and [tserr](https://github.com/thorstenrie/tserr)

## Usage

The package is installed with 

````go
go get github.com/thorstenrie/tsmock
````

In the Go app, the package is imported with

````go
import "github.com/thorstenrie/tsmock"
````

## Mock Stdin

The global mocked Stdin is provided by `tsmock.Stdin`

```go
stdin := tsmock.Stdin
```

The variable of type `*os.File` for the input to Stdin is set with `Set`

```go
err := stdin.Set(f)
```

Visibility of the input and a delay of processing each line of the input can be configured with `Visibility` and `Delay`

```go
stdin.Visibility(false)
err := stdin.Delay(time.Milliseconds * 250)
```

The mocked stdin is executed with `Run`.

```go
err := stdin.Run(context.Background())
```

The `context` can be used to cancel the execution, for example with a timeout.

```go
ctx, cancel := context.WithTimeout(context.Background(), time.Second)
defer cancel()
err := stdin.Run(ctx)
```

The input can be retrieved with `os.Stdin`

```go
s := bufio.NewScanner(os.Stdin)
for s.Scan() {
  fmt.Println(s.Text())
}
```

## Example

```go
package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/thorstenrie/tsfio"
	"github.com/thorstenrie/tsmock"
)

var (
	filename = tsfio.Filename("test.txt")
	contents = "Aragorn\nGandalf\nGimli\nLegolas\nGollum\n"
)

func main() {
	tsfio.WriteStr(filename, contents)
	f, _ := tsfio.OpenFile(filename)
	stdin := tsmock.Stdin
	stdin.Set(f)
	stdin.Visibility(false)
	stdin.Delay(time.Millisecond * 250)
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*999)
	defer cancel()
	stdin.Run(ctx)
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		fmt.Println(s.Text())
	}
	if e := stdin.Err(); e != nil {
		fmt.Println(e)
	}
	f.Close()
	tsfio.RemoveFile(filename)
}
```
[Go Playground](https://go.dev/play/p/nfksVqPNaCj)

