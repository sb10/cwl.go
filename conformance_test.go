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
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v2"
)

var conTestNum = flag.Int("ctest", 0, "run only this conformance test")

type conformanceTest struct {
	Tool   string
	Job    string
	Doc    string
	Output map[string]interface{} // []string | string | map[string]string | int | []int
}

type conformanceTests []conformanceTest

func TestConformance(t *testing.T) {
	// parse the conformance test file from the official spec repo (what cwltest
	// works on to test cwlref-runner)
	conformanceDir := filepath.Join("cwl", version)
	conformanceFilePath := filepath.Join(conformanceDir, "conformance_test_"+version+".yaml")
	conformanceFile, err := ioutil.ReadFile(conformanceFilePath)
	if err != nil {
		t.Fatalf("conformance file %s could not be read: %s", conformanceFilePath, err)
	}
	c := &conformanceTests{}
	err = yaml.Unmarshal(conformanceFile, c)
	if err != nil {
		t.Fatalf("conformance file %s could not be parsed: %s", conformanceFilePath, err)
	}

	assert := assert.New(t)

	// run each test specified there
	done := 0
	toDo := 24 // TODO: not yet fully compatible, working on conformance test by test, total 111
	for _, test := range *c {
		if *conTestNum != 0 {
			done++
			if done != *conTestNum {
				continue
			}
		}

		cwlPath := filepath.Join(conformanceDir, test.Tool)
		paramsPath := filepath.Join(conformanceDir, test.Job)

		tmpDir, err := ioutil.TempDir("", "cwlgo.tests")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tmpDir)
		outDir := filepath.Join(tmpDir, "output")
		tmpDirPrefix := filepath.Join(tmpDir, "tmp")
		tmpOutDirPrefix := filepath.Join(tmpDir, "tmpout")

		err = os.Mkdir(outDir, 0700)
		if err != nil {
			t.Fatal(err)
		}

		stagingDir := filepath.Join(tmpDir, "inputs")
		err = os.Mkdir(stagingDir, 0700)
		if err != nil {
			t.Fatal(err)
		}
		ifc := func(path string) string {
			staged := filepath.Join(stagingDir, filepath.Base(path))
			copyFile(path, staged)
			return staged
		}

		config := ResolveConfig{
			OutputDir:       outDir,
			TmpDirPrefix:    tmpDirPrefix,
			TmpOutDirPrefix: tmpOutDirPrefix,
			Cores:           2,
		}

		r, cmds, err := Resolve("myworkflow", cwlPath, paramsPath, config, ifc)
		if !assert.Nil(err, cwlPath+" failed to Resolve()") {
			break
		}

		assert.True(len(cmds) >= 1, test.Doc)

		// var output interface{}
		var erre error
		for _, cmd := range cmds {
			_, erre = cmd.Execute()
			if !assert.Nil(erre, test.Doc+" failed") {
				break
			}
		}
		output := r.Output()

		if erre == nil {
			// if we expect "Any" location, make the actual match
			for k, v := range test.Output {
				switch x := v.(type) {
				case map[interface{}]interface{}:
					if loc := x["location"]; loc == "Any" {
						if m, ok := output.(map[string]interface{}); ok {
							if n, exists := m[k]; exists {
								o := n.(map[interface{}]interface{})
								if _, exists := o["location"]; exists {
									o["location"] = "Any"
								}
							}
						}
					}
				}
			}

			assert.Equal(test.Output, output, test.Doc)
		}

		if *conTestNum != 0 {
			break
		}

		done++
		if done == toDo {
			break
		}
	}
}

// copyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, attempt to create a hard link
// between the two files. If that fails, copy the file contents from src to dst.
func copyFile(src, dst string) error {
	// from https://stackoverflow.com/a/21067803/675083
	sfi, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("copyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("copyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return err
		}
	}
	if err = os.Link(src, dst); err == nil {
		return err
	}
	return copyFileContents(src, dst)
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) error {
	// from https://stackoverflow.com/a/21067803/675083
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	err = out.Sync()
	return err
}
