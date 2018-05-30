# cwl.go

[![GoDoc](https://godoc.org/github.com/sb10/cwl.go?status.svg)](https://godoc.org/github.com/sb10/cwl.go)
[![Go Report Card](https://goreportcard.com/badge/github.com/sb10/cwl.go)](https://goreportcard.com/report/github.com/sb10/cwl.go)
develop branch:
[![Build Status](https://travis-ci.org/sb10/cwl.go.svg?branch=develop)](https://travis-ci.org/sb10/cwl.go)
[![Coverage Status](https://coveralls.io/repos/github/sb10/cwl.go/badge.svg?branch=develop)](https://coveralls.io/github/sb10/cwl.go?branch=develop)

`cwl.go` is a parser of [CWL](https://github.com/common-workflow-language/common-workflow-language) files and their input parameters, for example [1st-tool.yaml](https://github.com/common-workflow-language/common-workflow-language/blob/master/v1.0/examples/1st-tool.cwl) and [echo-job.yml](https://github.com/common-workflow-language/common-workflow-language/blob/master/v1.0/examples/echo-job.yml).

This is a fork of github.com/otiai10/cwl.go that, amongst other things, adds a
"Resolve" method to turn CWL+params in to concrete command lines to run, with
dependency and resource usage information.

# Example

```go
package main

import (
	"fmt"
	"os"

	cwl "github.com/sb10/cwl.go"
)

func main() {
	cwlFile, _ := os.Open("hello.cwl")
	paramsFile, _ := os.Open("params.yaml")

	// basic parsing
	doc := cwl.NewCWL()
	doc.Decode(cwlFile)
	fmt.Printf("%+v\n", doc)

	// or get concrete command lines to run
	cmds, _ := cwl.Resolve(cwlFile, paramsFile)
	fmt.Printf("%s\n", cmds)
}
```

# Tests

Decode() tests run against the test cwl files in the cwl subdirectory for the
latest release version of the CWL spec. To run all tests, just:

```sh
go test
```

To run tests against a single file:

```sh
go test -cwl scatter-valuefrom-wf3.cwl
```
