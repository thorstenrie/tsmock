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

The Go package tsmock provides an interface to test and mock Stdin based on files. It reads input from a file and
passes it to os.Stdin. It can be configured to set the visibility of the input and a delay in processing each line of the
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

## Example

