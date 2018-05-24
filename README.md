# cwl.go

[![Build Status](https://travis-ci.org/sb10/cwl.go.svg?branch=develop)](https://travis-ci.org/sb10/cwl.go) [![GoDoc](https://godoc.org/github.com/sb10/cwl.go?status.svg)](https://godoc.org/github.com/sb10/cwl.go)

`cwl.go` is just a parser of CWL file and input files based on [CWL](https://github.com/common-workflow-language/common-workflow-language), for example [1st-tool.yaml](https://github.com/common-workflow-language/common-workflow-language/blob/master/v1.0/examples/1st-tool.cwl) and [echo-job.yml](https://github.com/common-workflow-language/common-workflow-language/blob/master/v1.0/examples/echo-job.yml).

This is a fork of github.com/otiai10/cwl.go that, amongst other things, adds a
"Resolve" method to turn CWL+params in to concrete command lines to run, with
dependency information.

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
	doc := cwl.NewCWL()
	doc.Decode(cwlFile)

	paramsFile, _ := os.Open("params.yaml")
	cmds, _ := doc.Resolve(paramsFile)
	fmt.Printf("cmds:\n%s\n", cmds)
}
```

# Tests

## Prerequisite

`xtest.sh` requires Go package `github.com/otiai10/mint`

To install it.

```
go get -u github.com/otiai10/mint
```

## Why xtest.sh and How to do test with it.

Because there are both array and dictionary in CWL specification, and as you know Golang can't keep order of map keys, the test fails sometimes by order problem. Therefore, [`./xtest.sh`](https://github.com/sb10/cwl.go/blob/master/xtest.sh) tries testing each case several times eagerly unless it passes.

For all cases,

```sh
./xtest.sh
```

For only 1 case which matches `_wf3`,

```sh
./xtest.sh _wf3
```

Or if you want to execute single test for just 1 time (NOT eagerly),

```sh
go test ./tests -run _wf3
```
