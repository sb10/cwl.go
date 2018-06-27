# cwl.go

[![GoDoc](https://godoc.org/github.com/sb10/cwl.go?status.svg)](https://godoc.org/github.com/sb10/cwl.go)
[![Go Report Card](https://goreportcard.com/badge/github.com/sb10/cwl.go)](https://goreportcard.com/report/github.com/sb10/cwl.go)
develop branch:
[![Build Status](https://travis-ci.org/sb10/cwl.go.svg?branch=develop)](https://travis-ci.org/sb10/cwl.go)
[![Coverage Status](https://coveralls.io/repos/github/sb10/cwl.go/badge.svg?branch=develop)](https://coveralls.io/github/sb10/cwl.go?branch=develop)

`cwl.go` is a parser of [CWL](https://github.com/common-workflow-language/common-workflow-language) files and their input parameters, for example [1st-tool.cwl](https://github.com/common-workflow-language/common-workflow-language/blob/master/v1.0/examples/1st-tool.cwl) and [echo-job.yml](https://github.com/common-workflow-language/common-workflow-language/blob/master/v1.0/examples/echo-job.yml).

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

	// or get command lines to run
	r, cmds, _ := cwl.Resolve("hello.cwl", "params.yaml")
	for _, cmd := range cmds {
		// arrange to only execute these in the correct order according to the
		// dependency tree; you could potentially Execute each command on a
		// different host, pulling in the return value of GetPriorOutputs() and
		// sending the output back over the wire.
		output, _ := cmd.Execute(r.GetPriorOutputs())
		r.SetOutput(cmd.UniqueID, output, cmd.Parameters)
	}
	fmt.Printf("%+v\n", r.Output())
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
