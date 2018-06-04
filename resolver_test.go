// This file is part of cwl.go.
// Author: Sendu Bala <sb10@sanger.ac.uk>.
//
// Copyright © 2018 Genome Research Limited
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package cwl

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testDataDir = "testdata"

var update = flag.Bool("update", false, "update .golden files")

func TestResolver(t *testing.T) {
	files, err := ioutil.ReadDir(testDataDir)
	if err != nil {
		t.Fatal(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)

	// run Resolve() on every cwl file in our testdata folder, and check that it
	// it produces the correct commands
	for _, file := range files {
		name := file.Name()
		cwlPath := filepath.Join(testDataDir, name)

		if !strings.HasSuffix(cwlPath, ".cwl") {
			continue
		}

		paramsPath := strings.Replace(cwlPath, ".cwl", ".yaml", 1)
		if _, err := os.Stat(paramsPath); err != nil && os.IsNotExist(err) {
			paramsPath = ""
		}

		actual, err := Resolve(cwlPath, paramsPath, ResolveConfig{}, DefaultInputFileCallback)
		if !assert.Nil(err, name+" failed to Resolve()") {
			continue
		}

		// we need to strip current working dir from each cmd's Cwd, so we don't
		// store the author's directory in the golden files, or expect author's
		// directory
		for _, cmd := range actual {
			if strings.HasPrefix(cmd.Cwd, wd) {
				cmd.Cwd = strings.Replace(cmd.Cwd, wd, ".", 1)
			}
		}

		goldenPath := cwlPath + ".golden"
		if *update {
			writeCommands(t, goldenPath, actual)
		}

		expected := readCommands(t, goldenPath)
		assert.Equal(expected, actual, name+" did not Resolve() correctly")
	}
}

func openFile(t *testing.T, path string) *os.File {
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	return f
}

func writeCommands(t *testing.T, path string, cmds Commands) {
	j, err := json.Marshal(cmds)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile(path, j, 0644)
	if err != nil {
		t.Fatal(err)
	}
}

func readCommands(t *testing.T, path string) Commands {
	f := openFile(t, path)
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	var cmds Commands
	err = json.Unmarshal(b, &cmds)
	if err != nil {
		t.Fatal(err)
	}
	return cmds
}
