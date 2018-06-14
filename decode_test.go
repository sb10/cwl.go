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
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	a "github.com/stretchr/testify/assert"
)

const version = "v1.0"

var cwlPath = flag.String("cwl", "", "run tests on only this cwl file")

type expectationTester func(assert *a.Assertions, root *Root)

func TestDecode(t *testing.T) {
	testOfficialDir := filepath.Join("cwl", version, version)
	files, err := ioutil.ReadDir(testOfficialDir)
	if err != nil {
		t.Fatal(err)
	}

	if len(files) < 10 {
		t.Fatal("cwl directory is a git submodule that must be checked out")
	}

	// in addition to the common decode tests, we want to assert particular
	// things depending on the CWL file in question
	expectations := map[string]expectationTester{
		"basename-fields-test.cwl":                 basenameFieldsTest,
		"binding-test.cwl":                         bindingTest,
		"bwa-mem-tool.cwl":                         bwaMemToolTest,
		"cat1-testcli.cwl":                         cat1TestcliTest,
		"cat3-nodocker.cwl":                        cat3NodockerTest,
		"cat3-tool-mediumcut.cwl":                  cat3ToolMediumcutTest,
		"cat3-tool-shortcut.cwl":                   cat3ToolShortcutTest,
		"cat3-tool.cwl":                            cat3ToolTest,
		"cat4-tool.cwl":                            cat4ToolTest,
		"cat5-tool.cwl":                            cat5ToolTest,
		"conflict-wf.cwl":                          conflictWfTest,
		"count-lines10-wf.cwl":                     countLines10WfTest,
		"count-lines11-wf.cwl":                     countLines11WfTest,
		"count-lines12-wf.cwl":                     countLines12WfTest,
		"count-lines1-wf.cwl":                      countLines1WfTest,
		"count-lines2-wf.cwl":                      countLines2WfTest,
		"count-lines3-wf.cwl":                      countLines3WfTest,
		"count-lines4-wf.cwl":                      countLines4WfTest,
		"count-lines5-wf.cwl":                      countLines5WfTest,
		"count-lines6-wf.cwl":                      countLines6WfTest,
		"count-lines7-wf.cwl":                      countLines7WfTest,
		"count-lines8-wf.cwl":                      countLines8WfTest,
		"count-lines9-wf.cwl":                      countLines9WfTest,
		"default_path.cwl":                         defaultPathTest,
		"dir2.cwl":                                 dir2Test,
		"dir3.cwl":                                 dir3Test,
		"dir4.cwl":                                 dir4Test,
		"dir5.cwl":                                 dir5Test,
		"dir6.cwl":                                 dir6Test,
		"dir7.cwl":                                 dir7Test,
		"dir.cwl":                                  dirTest,
		"docker-array-secondaryfiles.cwl":          dockerArraySecondaryfilesTest,
		"docker-output-dir.cwl":                    dockerOutputDirTest,
		"dynresreq.cwl":                            dynresreqTest,
		"echo-file-tool.cwl":                       echoFileToolTest,
		"echo-tool.cwl":                            echoToolTest,
		"env-tool1.cwl":                            envTool1Test,
		"env-tool2.cwl":                            envTool2Test,
		"envvar2.cwl":                              envvar2Test,
		"envvar.cwl":                               envvarTest,
		"env-wf1.cwl":                              envWf1Test,
		"env-wf2.cwl":                              envWf2Test,
		"env-wf3.cwl":                              envWf3Test,
		"file-literal-ex.cwl":                      fileLiteralExTest,
		"formattest2.cwl":                          formattest2Test,
		"formattest3.cwl":                          formattest3Test,
		"formattest.cwl":                           formattestTest,
		"glob-expr-list.cwl":                       globExprListTest,
		"imported-hint.cwl":                        importedHintTest,
		"initialworkdirrequirement-docker-out.cwl": initialworkdirrequirementDockerOutTest,
		"initialwork-path.cwl":                     initialworkPathTest,
		"inline-js.cwl":                            inlineJsTest,
		"js-expr-req-wf.cwl":                       jsExprReqWfTest,
		"metadata.cwl":                             metadataTest,
		"nameroot.cwl":                             namerootTest,
		"nested-array.cwl":                         nestedArrayTest,
		"null-defined.cwl":                         nullDefinedTest,
		"null-expression1-tool.cwl":                nullExpression1ToolTest,
		"null-expression2-tool.cwl":                nullExpression2ToolTest,
		"optional-output.cwl":                      optionalOutputTest,
		"params2.cwl":                              params2Test,
		"params.cwl":                               paramsTest,
		"parseInt-tool.cwl":                        parseIntToolTest,
		"record-output.cwl":                        recordOutputTest,
		"recursive-input-directory.cwl":            recursiveInputDirectoryTest,
		"rename.cwl":                               renameTest,
		"revsort-packed.cwl":                       revsortPackedTest,
		"revsort.cwl":                              revsortTest,
		"revtool.cwl":                              revtoolTest,
		"scatter-valueFrom-tool.cwl":               scatterValueFromToolTest,
		"scatter-valuefrom-wf1.cwl":                scatterValuefromWf1Test,
		"scatter-valuefrom-wf2.cwl":                scatterValuefromWf2Test,
		"scatter-valuefrom-wf3.cwl":                scatterValuefromWf3Test,
		"scatter-valuefrom-wf4.cwl":                scatterValuefromWf4Test,
		"scatter-valuefrom-wf5.cwl":                scatterValuefromWf5Test,
		"scatter-valuefrom-wf6.cwl":                scatterValuefromWf6Test,
		"scatter-wf1.cwl":                          scatterWf1Test,
		"scatter-wf2.cwl":                          scatterWf2Test,
		"scatter-wf3.cwl":                          scatterWf3Test,
		"scatter-wf4.cwl":                          scatterWf4Test,
		"schemadef-tool.cwl":                       schemadefToolTest,
		"schemadef-wf.cwl":                         schemadefWfTest,
		"search.cwl":                               searchTest,
		"shellchar2.cwl":                           shellchar2Test,
		"shellchar.cwl":                            shellcharTest,
		"shelltest.cwl":                            shelltestTest,
		"sorttool.cwl":                             sorttoolTest,
		"stagefile.cwl":                            stagefileTest,
		"stderr-mediumcut.cwl":                     stderrMediumcutTest,
		"stderr-shortcut.cwl":                      stderrShortcutTest,
		"stderr.cwl":                               stderrTest,
		"step-valuefrom2-wf.cwl":                   stepValuefrom2WfTest,
		"step-valuefrom3-wf.cwl":                   stepValuefrom3WfTest,
		"step-valuefrom-wf.cwl":                    stepValuefromWfTest,
		"sum-wf.cwl":                               sumWfTest,
		"template-tool.cwl":                        templateToolTest,
		"test-cwl-out2.cwl":                        testCwlOut2Test,
		"test-cwl-out.cwl":                         testCwlOutTest,
		"tmap-tool.cwl":                            tmapToolTest,
		"wc2-tool.cwl":                             wc2ToolTest,
		"wc3-tool.cwl":                             wc3ToolTest,
		"wc4-tool.cwl":                             wc4ToolTest,
		"wc-tool.cwl":                              wcToolTest,
		"writable-dir.cwl":                         writableDirTest,
	}

	// test Decode() against every cwl file in the official spec repo, in
	// parallel
	for _, file := range files {
		name := file.Name()
		path := filepath.Join(testOfficialDir, name)

		if !strings.HasSuffix(path, ".cwl") {
			continue
		}

		if *cwlPath != "" {
			if !strings.Contains(name, *cwlPath) {
				continue
			}
		}

		t.Run(name, func(pt *testing.T) {
			pt.Parallel()

			assert := a.New(pt)

			f := openFile(pt, path)
			defer f.Close()

			root := NewCWL()
			err := root.Decode(f)
			assert.Nil(err)
			assert.Equal(version, root.Version)

			if tester, exists := expectations[name]; exists {
				root.Sort()
				tester(assert, root)
			} else {
				pt.Fatal("missing detailed tests for " + name)
			}
		})
	}
}

func openFile(t *testing.T, path string) *os.File {
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	return f
}

func lineNumber() string {
	_, _, line, _ := runtime.Caller(1)
	return fmt.Sprintf("line %d", line)
}

func basenameFieldsTest(assert *a.Assertions, root *Root) {
	assert.Equal("Workflow", root.Class)
	assert.Equal("StepInputExpressionRequirement", root.Requirements[0].Class)

	assert.Equal(1, len(root.Inputs))
	assert.Equal("tool", root.Inputs[0].ID)
	assert.Equal("File", root.Inputs[0].Types[0].Type)

	assert.Equal(2, len(root.Outputs), lineNumber())
	assert.Equal("extFile", root.Outputs[0].ID)
	assert.Equal("File", root.Outputs[0].Types[0].Type)
	assert.Equal("ext/out", root.Outputs[0].Source[0])
	assert.Equal("rootFile", root.Outputs[1].ID)
	assert.Equal("File", root.Outputs[1].Types[0].Type)
	assert.Equal("root/out", root.Outputs[1].Source[0])

	assert.Equal(2, len(root.Steps))
	for _, st := range root.Steps {
		assert.Equal("echo-file-tool.cwl", st.Run.Value)
		for _, in := range st.In {
			switch in.ID {
			case "in":
				assert.Equal(fmt.Sprintf("$(inputs.tool.name%s)", st.ID), in.ValueFrom)
			case "tool":
				assert.Equal("", in.ValueFrom)
			}
		}
		assert.Equal("out", st.Out[0].ID)
	}
}

func bindingTest(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())

	assert.Equal("DockerRequirement", root.Hints[0].Class, lineNumber())
	assert.Equal("python:2-slim", root.Hints[0].DockerPull, lineNumber())

	assert.Equal("#args.py", root.Inputs[0].ID, lineNumber())
	assert.Equal("File", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal(reflect.Map, root.Inputs[0].Default.Kind, lineNumber())
	assert.Equal(-1, root.Inputs[0].Binding.Position, lineNumber())
	assert.Equal("reference", root.Inputs[1].ID, lineNumber())
	assert.Equal("File", root.Inputs[1].Types[0].Type, lineNumber())
	assert.Equal(2, root.Inputs[1].Binding.Position, lineNumber())
	assert.Equal("reads", root.Inputs[2].ID, lineNumber())
	assert.Equal("array", root.Inputs[2].Types[0].Type, lineNumber())
	assert.Equal("File", root.Inputs[2].Types[0].Items[0].Type, lineNumber())
	assert.Equal("-YYY", root.Inputs[2].Types[0].Binding.Prefix, lineNumber())
	assert.Equal(3, root.Inputs[2].Binding.Position, lineNumber())
	assert.Equal("-XXX", root.Inputs[2].Binding.Prefix, lineNumber())

	assert.Equal("args", root.Outputs[0].ID, lineNumber())
	assert.Equal("string[]", root.Outputs[0].Types[0].Type, lineNumber())
}

func bwaMemToolTest(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.IsType(Hints{}, root.Hints)
	assert.Equal("ResourceRequirement", root.Hints[0].Class, lineNumber())
	assert.Equal(2, root.Hints[0].CoresMin, lineNumber())

	assert.Equal(5, len(root.Inputs), lineNumber())
	assert.IsType(Inputs{}, root.Inputs)

	assert.Equal("args.py", root.Inputs[0].ID, lineNumber())
	assert.Equal("File", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal(-1, root.Inputs[0].Binding.Position, lineNumber())

	assert.Equal("min_std_max_min", root.Inputs[1].ID, lineNumber())
	assert.Equal("array", root.Inputs[1].Types[0].Type, lineNumber())
	assert.Equal("int", root.Inputs[1].Types[0].Items[0].Type, lineNumber())
	assert.Equal(1, root.Inputs[1].Binding.Position, lineNumber())
	assert.Equal(",", root.Inputs[1].Binding.Separator, lineNumber())

	assert.Equal("minimum_seed_length", root.Inputs[2].ID, lineNumber())
	assert.Equal("int", root.Inputs[2].Types[0].Type, lineNumber())
	assert.Equal(1, root.Inputs[2].Binding.Position, lineNumber())
	assert.Equal("-m", root.Inputs[2].Binding.Prefix, lineNumber())

	assert.Equal("reference", root.Inputs[3].ID, lineNumber())
	assert.Equal("File", root.Inputs[3].Types[0].Type, lineNumber())
	assert.Equal(2, root.Inputs[3].Binding.Position, lineNumber())

	assert.Equal("reads", root.Inputs[4].ID, lineNumber())
	assert.Equal("array", root.Inputs[4].Types[0].Type, lineNumber())
	assert.Equal("File", root.Inputs[4].Types[0].Items[0].Type, lineNumber())
	assert.Equal(3, root.Inputs[4].Binding.Position, lineNumber())

	assert.Equal("sam", root.Outputs[0].ID, lineNumber())
	assert.Equal("null", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("File", root.Outputs[0].Types[1].Type, lineNumber())
	assert.Equal([]string{"output.sam"}, root.Outputs[0].Binding.Glob, lineNumber())
	assert.Equal("args", root.Outputs[1].ID, lineNumber())
	assert.Equal("array", root.Outputs[1].Types[0].Type, lineNumber())
	assert.Equal("string", root.Outputs[1].Types[0].Items[0].Type, lineNumber())
}

func cat1TestcliTest(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal("Print the contents of a file to stdout using 'cat' running in a docker container.", root.Doc, lineNumber())

	assert.Equal("DockerRequirement", root.Hints[0].Class, lineNumber())
	assert.Equal("python:2-slim", root.Hints[0].DockerPull, lineNumber())

	assert.Equal("args.py", root.Inputs[0].ID, lineNumber())
	assert.Equal("File", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal(reflect.Map, root.Inputs[0].Default.Kind, lineNumber())
	assert.Equal(-1, root.Inputs[0].Binding.Position, lineNumber())
	assert.Equal("numbering", root.Inputs[1].ID, lineNumber())
	assert.Equal("null", root.Inputs[1].Types[0].Type, lineNumber())
	assert.Equal("boolean", root.Inputs[1].Types[1].Type, lineNumber())
	assert.Equal("file1", root.Inputs[2].ID, lineNumber())
	assert.Equal("File", root.Inputs[2].Types[0].Type, lineNumber())
	assert.Equal(1, root.Inputs[2].Binding.Position, lineNumber())

	assert.Equal("args", root.Outputs[0].ID, lineNumber())
	assert.Equal("string[]", root.Outputs[0].Types[0].Type, lineNumber())

	assert.Equal("python", root.BaseCommands[0], lineNumber())
	assert.Equal("cat", root.Arguments[0].Value, lineNumber())
}

func cat3NodockerTest(assert *a.Assertions, root *Root) {
	assert.Equal("Print the contents of a file to stdout using 'cat'.", root.Doc, lineNumber())
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal(1, len(root.BaseCommands), lineNumber())
	assert.Equal("cat", root.BaseCommands[0], lineNumber())
	assert.Equal("output.txt", root.Stdout, lineNumber())
	assert.Equal(1, len(root.Inputs), lineNumber())
	assert.Equal("file1", root.Inputs[0].ID, lineNumber())
	assert.Equal("File", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("Input File", root.Inputs[0].Label, lineNumber())
	assert.Equal("The file that will be copied using 'cat'", root.Inputs[0].Doc, lineNumber())
	assert.Equal(1, root.Inputs[0].Binding.Position, lineNumber())
}

func cat3ToolMediumcutTest(assert *a.Assertions, root *Root) {
	assert.Equal("Print the contents of a file to stdout using 'cat' running in a docker container.", root.Doc, lineNumber())
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal(1, len(root.BaseCommands), lineNumber())
	assert.Equal("cat", root.BaseCommands[0], lineNumber())
	assert.Equal("cat-out", root.Stdout, lineNumber())
	assert.IsType(Hints{}, root.Hints)
	assert.Equal("DockerRequirement", root.Hints[0].Class, lineNumber())
	assert.Equal("debian:wheezy", root.Hints[0].DockerPull, lineNumber())
	assert.Equal(1, len(root.Inputs), lineNumber())
	assert.Equal("file1", root.Inputs[0].ID, lineNumber())
	assert.Equal("File", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("Input File", root.Inputs[0].Label, lineNumber())
	assert.Equal("The file that will be copied using 'cat'", root.Inputs[0].Doc, lineNumber())
	assert.Equal(1, root.Inputs[0].Binding.Position, lineNumber())
}

func cat3ToolShortcutTest(assert *a.Assertions, root *Root) {
	assert.Equal("Print the contents of a file to stdout using 'cat' running in a docker container.", root.Doc, lineNumber())
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal(1, len(root.BaseCommands), lineNumber())
	assert.Equal("cat", root.BaseCommands[0], lineNumber())
	assert.IsType(Hints{}, root.Hints)
	assert.Equal("DockerRequirement", root.Hints[0].Class, lineNumber())
	assert.Equal("debian:wheezy", root.Hints[0].DockerPull, lineNumber())
	assert.Equal(1, len(root.Inputs), lineNumber())
	assert.Equal("file1", root.Inputs[0].ID, lineNumber())
	assert.Equal("File", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("Input File", root.Inputs[0].Label, lineNumber())
	assert.Equal("The file that will be copied using 'cat'", root.Inputs[0].Doc, lineNumber())
	assert.Equal(1, root.Inputs[0].Binding.Position, lineNumber())
}

func cat3ToolTest(assert *a.Assertions, root *Root) {
	assert.Equal("Print the contents of a file to stdout using 'cat' running in a docker container.", root.Doc, lineNumber())
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal(1, len(root.BaseCommands), lineNumber())
	assert.Equal("cat", root.BaseCommands[0], lineNumber())
	assert.Equal("output.txt", root.Stdout, lineNumber())
	assert.IsType(Hints{}, root.Hints)
	assert.Equal("DockerRequirement", root.Hints[0].Class, lineNumber())
	assert.Equal("debian:wheezy", root.Hints[0].DockerPull, lineNumber())
	assert.Equal(1, len(root.Inputs), lineNumber())
	assert.Equal("file1", root.Inputs[0].ID, lineNumber())
	assert.Equal("File", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("Input File", root.Inputs[0].Label, lineNumber())
	assert.Equal("The file that will be copied using 'cat'", root.Inputs[0].Doc, lineNumber())
	assert.Equal(1, root.Inputs[0].Binding.Position, lineNumber())
}

func cat4ToolTest(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal(1, len(root.Inputs), lineNumber())
	assert.Equal("file1", root.Inputs[0].ID, lineNumber())
	assert.Equal(1, len(root.Outputs), lineNumber())
	assert.Equal("output_txt", root.Outputs[0].ID, lineNumber())
	assert.Equal("File", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal([]string{"output.txt"}, root.Outputs[0].Binding.Glob, lineNumber())
	assert.Equal(1, len(root.BaseCommands), lineNumber())
	assert.Equal("cat", root.BaseCommands[0], lineNumber())
	assert.Equal("output.txt", root.Stdout, lineNumber())
	assert.Equal("$(inputs.file1.path)", root.Stdin, lineNumber())
}

func cat5ToolTest(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal("Print the contents of a file to stdout using 'cat' running in a docker container.", root.Doc, lineNumber())
	assert.Equal(2, len(root.Hints), lineNumber())
	assert.IsType(Hints{}, root.Hints)
	assert.Equal("DockerRequirement", root.Hints[0].Class, lineNumber())
	assert.Equal("debian:wheezy", root.Hints[0].DockerPull, lineNumber())
	assert.Equal("ex:BlibberBlubberFakeRequirement", root.Hints[1].Class, lineNumber())
	assert.Equal("fraggleFroogle", root.Hints[1].FakeField, lineNumber())
	assert.Equal(1, len(root.Inputs), lineNumber())
	assert.Equal("file1", root.Inputs[0].ID, lineNumber())
	assert.Equal("File", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("Input File", root.Inputs[0].Label, lineNumber())
	assert.Equal("The file that will be copied using 'cat'", root.Inputs[0].Doc, lineNumber())
	assert.Equal(1, root.Inputs[0].Binding.Position, lineNumber())
	assert.Equal(1, len(root.Outputs), lineNumber())
	assert.Equal("output_file", root.Outputs[0].ID, lineNumber())
	assert.Equal("File", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal([]string{"output.txt"}, root.Outputs[0].Binding.Glob, lineNumber())
	assert.Equal(1, len(root.BaseCommands), lineNumber())
	assert.Equal("cat", root.BaseCommands[0], lineNumber())
	assert.Equal("output.txt", root.Stdout, lineNumber())
	// $namespaces
	assert.Equal(1, len(root.Namespaces), lineNumber())
	assert.Equal("http://example.com/", root.Namespaces[0]["ex"], lineNumber())
}

func conflictWfTest(assert *a.Assertions, root *Root) {
	assert.Equal("echo", root.Graphs[0].ID, lineNumber())
	assert.Equal("CommandLineTool", root.Graphs[0].Class, lineNumber())
	assert.Equal("text", root.Graphs[0].Inputs[0].ID, lineNumber())
	assert.Equal("string", root.Graphs[0].Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("fileout", root.Graphs[0].Outputs[0].ID, lineNumber())
	assert.Equal("File", root.Graphs[0].Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("out.txt", root.Graphs[0].Outputs[0].Binding.Glob[0], lineNumber())
	assert.Equal("echo", root.Graphs[0].BaseCommands[0], lineNumber())
	assert.Equal("out.txt", root.Graphs[0].Stdout, lineNumber())

	assert.Equal("cat", root.Graphs[1].ID, lineNumber())
	assert.Equal("CommandLineTool", root.Graphs[1].Class, lineNumber())
	assert.Equal("file1", root.Graphs[1].Inputs[0].ID, lineNumber())
	assert.Equal("File", root.Graphs[1].Inputs[0].Types[0].Type, lineNumber())
	assert.Equal(1, root.Graphs[1].Inputs[0].Binding.Position, lineNumber())
	assert.Equal("file2", root.Graphs[1].Inputs[1].ID, lineNumber())
	assert.Equal("File", root.Graphs[1].Inputs[1].Types[0].Type, lineNumber())
	assert.Equal(2, root.Graphs[1].Inputs[1].Binding.Position, lineNumber())
	assert.Equal("fileout", root.Graphs[1].Outputs[0].ID, lineNumber())
	assert.Equal("File", root.Graphs[1].Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("out.txt", root.Graphs[1].Outputs[0].Binding.Glob[0], lineNumber())
	assert.Equal("cat", root.Graphs[1].BaseCommands[0], lineNumber())
	assert.Equal("out.txt", root.Graphs[1].Stdout, lineNumber())

	assert.Equal("collision", root.Graphs[2].ID, lineNumber())
	assert.Equal("Workflow", root.Graphs[2].Class, lineNumber())
	assert.Equal("input_1", root.Graphs[2].Inputs[0].ID, lineNumber())
	assert.Equal("string", root.Graphs[2].Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("input_2", root.Graphs[2].Inputs[1].ID, lineNumber())
	assert.Equal("string", root.Graphs[2].Inputs[1].Types[0].Type, lineNumber())
	assert.Equal("fileout", root.Graphs[2].Outputs[0].ID, lineNumber())
	assert.Equal("File", root.Graphs[2].Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("cat_step/fileout", root.Graphs[2].Outputs[0].Source[0], lineNumber())
	assert.Equal("cat_step", root.Graphs[2].Steps[0].ID, lineNumber())
	assert.Equal("#cat", root.Graphs[2].Steps[0].Run.Value, lineNumber())
	assert.Equal("file1", root.Graphs[2].Steps[0].In[0].ID, lineNumber())
	assert.Equal("echo_1/fileout", root.Graphs[2].Steps[0].In[0].Source[0], lineNumber())
	assert.Equal("file2", root.Graphs[2].Steps[0].In[1].ID, lineNumber())
	assert.Equal("echo_2/fileout", root.Graphs[2].Steps[0].In[1].Source[0], lineNumber())
	assert.Equal("fileout", root.Graphs[2].Steps[0].Out[0].ID, lineNumber())
	assert.Equal("echo_1", root.Graphs[2].Steps[1].ID, lineNumber())
	assert.Equal("#echo", root.Graphs[2].Steps[1].Run.Value, lineNumber())
	assert.Equal("text", root.Graphs[2].Steps[1].In[0].ID, lineNumber())
	assert.Equal("input_1", root.Graphs[2].Steps[1].In[0].Source[0], lineNumber())
	assert.Equal("fileout", root.Graphs[2].Steps[1].Out[0].ID, lineNumber())
	assert.Equal("echo_2", root.Graphs[2].Steps[2].ID, lineNumber())
	assert.Equal("#echo", root.Graphs[2].Steps[2].Run.Value, lineNumber())
	assert.Equal("text", root.Graphs[2].Steps[2].In[0].ID, lineNumber())
	assert.Equal("input_2", root.Graphs[2].Steps[2].In[0].Source[0], lineNumber())
	assert.Equal("fileout", root.Graphs[2].Steps[2].Out[0].ID, lineNumber())
}

func countLines10WfTest(assert *a.Assertions, root *Root) {
	assert.Equal("Workflow", root.Class, lineNumber())

	assert.Equal("file1", root.Inputs[0].ID, lineNumber())
	assert.Equal("File", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("count_output", root.Outputs[0].ID, lineNumber())
	assert.Equal("int", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal([]string{"step1/count_output"}, root.Outputs[0].Source, lineNumber())

	assert.Equal("SubworkflowFeatureRequirement", root.Requirements[0].Class, lineNumber())

	assert.Equal("step1", root.Steps[0].ID, lineNumber())
	assert.Equal("file1", root.Steps[0].In[0].ID, lineNumber())
	assert.Equal("file1", root.Steps[0].In[0].Source[0], lineNumber())
	assert.Equal("count_output", root.Steps[0].Out[0].ID, lineNumber())

	assert.Equal("Workflow", root.Steps[0].Run.Workflow.Class, lineNumber())
	assert.Equal("file1", root.Steps[0].Run.Workflow.Inputs[0].ID, lineNumber())
	assert.Equal("File", root.Steps[0].Run.Workflow.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("count_output", root.Steps[0].Run.Workflow.Outputs[0].ID, lineNumber())
	assert.Equal("int", root.Steps[0].Run.Workflow.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal([]string{"step2/output"}, root.Steps[0].Run.Workflow.Outputs[0].Source, lineNumber())
	// Recursive steps
	assert.Equal("step1", root.Steps[0].Run.Workflow.Steps[0].ID, lineNumber())
	assert.Equal("wc-tool.cwl", root.Steps[0].Run.Workflow.Steps[0].Run.Value, lineNumber())
	assert.Equal("file1", root.Steps[0].Run.Workflow.Steps[0].In[0].ID, lineNumber())
	assert.Equal("file1", root.Steps[0].Run.Workflow.Steps[0].In[0].Source[0], lineNumber())
	assert.Equal("output", root.Steps[0].Run.Workflow.Steps[0].Out[0].ID, lineNumber())
	assert.Equal("step2", root.Steps[0].Run.Workflow.Steps[1].ID, lineNumber())
	assert.Equal("parseInt-tool.cwl", root.Steps[0].Run.Workflow.Steps[1].Run.Value, lineNumber())
	assert.Equal("file1", root.Steps[0].Run.Workflow.Steps[1].In[0].ID, lineNumber())
	assert.Equal("step1/output", root.Steps[0].Run.Workflow.Steps[1].In[0].Source[0], lineNumber())
	assert.Equal("output", root.Steps[0].Run.Workflow.Steps[1].Out[0].ID, lineNumber())
}

func countLines11WfTest(assert *a.Assertions, root *Root) {
	assert.Equal("Workflow", root.Class, lineNumber())

	assert.Equal("file1", root.Inputs[0].ID, lineNumber())
	assert.Equal("File?", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("count_output", root.Outputs[0].ID, lineNumber())
	assert.Equal("int", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal([]string{"step2/output"}, root.Outputs[0].Source, lineNumber())

	assert.Equal("step1", root.Steps[0].ID, lineNumber())
	assert.Equal("wc-tool.cwl", root.Steps[0].Run.Value, lineNumber())
	assert.Equal("file1", root.Steps[0].In[0].ID, lineNumber())
	assert.Equal(reflect.Map, root.Steps[0].In[0].Default.Kind, lineNumber())
	assert.Equal("output", root.Steps[0].Out[0].ID, lineNumber())

	assert.Equal("step2", root.Steps[1].ID, lineNumber())
	assert.Equal("parseInt-tool.cwl", root.Steps[1].Run.Value, lineNumber())
	assert.Equal("file1", root.Steps[1].In[0].ID, lineNumber())
	assert.Equal("step1/output", root.Steps[1].In[0].Source[0], lineNumber())
	assert.Equal("output", root.Steps[1].Out[0].ID, lineNumber())
}

func countLines12WfTest(assert *a.Assertions, root *Root) {
	assert.Equal("Workflow", root.Class, lineNumber())

	assert.Equal("file1", root.Inputs[0].ID, lineNumber())
	assert.Equal("array", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("File", root.Inputs[0].Types[0].Items[0].Type, lineNumber())
	assert.Equal("file2", root.Inputs[1].ID, lineNumber())
	assert.Equal("array", root.Inputs[1].Types[0].Type, lineNumber())
	assert.Equal("File", root.Inputs[1].Types[0].Items[0].Type, lineNumber())
	assert.Equal("count_output", root.Outputs[0].ID, lineNumber())
	assert.Equal("int", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal([]string{"step1/output"}, root.Outputs[0].Source, lineNumber())

	assert.Equal("MultipleInputFeatureRequirement", root.Requirements[0].Class, lineNumber())

	assert.Equal("step1", root.Steps[0].ID, lineNumber())
	assert.Equal("wc3-tool.cwl", root.Steps[0].Run.Value, lineNumber())
	assert.Equal("file1", root.Steps[0].In[0].ID, lineNumber())
	assert.Equal("file1", root.Steps[0].In[0].Source[0], lineNumber())
	assert.Equal("file2", root.Steps[0].In[0].Source[1], lineNumber())
	assert.Equal("merge_flattened", root.Steps[0].In[0].LinkMerge, lineNumber())
	assert.Equal("output", root.Steps[0].Out[0].ID, lineNumber())
}

func countLines1WfTest(assert *a.Assertions, root *Root) {
	assert.Equal("Workflow", root.Class, lineNumber())

	assert.Equal("file1", root.Inputs[0].ID, lineNumber())
	assert.Equal("File", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("count_output", root.Outputs[0].ID, lineNumber())
	assert.Equal("int", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal([]string{"step2/output"}, root.Outputs[0].Source, lineNumber())

	assert.Equal("step1", root.Steps[0].ID, lineNumber())
	assert.Equal("wc-tool.cwl", root.Steps[0].Run.Value, lineNumber())
	assert.Equal("file1", root.Steps[0].In[0].ID, lineNumber())
	assert.Equal([]string{"file1"}, root.Steps[0].In[0].Source, lineNumber())
	assert.Equal("output", root.Steps[0].Out[0].ID, lineNumber())
	assert.Equal("step2", root.Steps[1].ID, lineNumber())
	assert.Equal("parseInt-tool.cwl", root.Steps[1].Run.Value, lineNumber())
	assert.Equal("file1", root.Steps[1].In[0].ID, lineNumber())
	assert.Equal([]string{"step1/output"}, root.Steps[1].In[0].Source, lineNumber())
	assert.Equal("output", root.Steps[1].Out[0].ID, lineNumber())
}

func countLines2WfTest(assert *a.Assertions, root *Root) {
	assert.Equal("Workflow", root.Class, lineNumber())

	assert.Equal("InlineJavascriptRequirement", root.Requirements[0].Class, lineNumber())
	assert.Equal("file1", root.Inputs[0].ID, lineNumber())
	assert.Equal("File", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("count_output", root.Outputs[0].ID, lineNumber())
	assert.Equal("int", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal([]string{"step2/parseInt_output"}, root.Outputs[0].Source, lineNumber())

	assert.Equal("wc_file1", root.Steps[0].In[0].ID, lineNumber())
	assert.Equal("file1", root.Steps[0].In[0].Source[0], lineNumber())
	assert.Equal("wc_output", root.Steps[0].Out[0].ID, lineNumber())
	assert.IsType(Run{}, root.Steps[0].Run)
	assert.Equal("wc", root.Steps[0].Run.Workflow.ID, lineNumber())
	assert.Equal("CommandLineTool", root.Steps[0].Run.Workflow.Class, lineNumber())
	assert.Equal("wc_file1", root.Steps[0].Run.Workflow.Inputs[0].ID, lineNumber())
	assert.Equal("File", root.Steps[0].Run.Workflow.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("wc_output", root.Steps[0].Run.Workflow.Outputs[0].ID, lineNumber())
	assert.Equal("File", root.Steps[0].Run.Workflow.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal([]string{"output.txt"}, root.Steps[0].Run.Workflow.Outputs[0].Binding.Glob, lineNumber())
	assert.Equal("output.txt", root.Steps[0].Run.Workflow.Stdout, lineNumber())
	assert.Equal("wc", root.Steps[0].Run.Workflow.BaseCommands[0], lineNumber())
	assert.Equal("parseInt_file1", root.Steps[1].In[0].ID, lineNumber())
	assert.Equal("step1/wc_output", root.Steps[1].In[0].Source[0], lineNumber())
	assert.Equal("parseInt_output", root.Steps[1].Out[0].ID, lineNumber())
	assert.Equal("ExpressionTool", root.Steps[1].Run.Workflow.Class, lineNumber())
	assert.Equal("parseInt_file1", root.Steps[1].Run.Workflow.Inputs[0].ID, lineNumber())
	assert.Equal("File", root.Steps[1].Run.Workflow.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal(true, root.Steps[1].Run.Workflow.Inputs[0].Binding.LoadContents, lineNumber())
	assert.Equal("parseInt_output", root.Steps[1].Run.Workflow.Outputs[0].ID, lineNumber())
	assert.Equal("int", root.Steps[1].Run.Workflow.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("${return {'parseInt_output': parseInt(inputs.parseInt_file1.contents)};}\n", root.Steps[1].Run.Workflow.Expression, lineNumber())
}

func countLines3WfTest(assert *a.Assertions, root *Root) {
	assert.Equal("Workflow", root.Class, lineNumber())

	assert.Equal("file1", root.Inputs[0].ID, lineNumber())
	assert.Equal("File[]", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("count_output", root.Outputs[0].ID, lineNumber())
	assert.Equal("int[]", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal([]string{"step1/output"}, root.Outputs[0].Source, lineNumber())

	assert.Equal("ScatterFeatureRequirement", root.Requirements[0].Class, lineNumber())
	assert.Equal("step1", root.Steps[0].ID, lineNumber())
	assert.Equal("wc2-tool.cwl", root.Steps[0].Run.Value, lineNumber())
	assert.Equal("file1", root.Steps[0].Scatter[0], lineNumber())
	assert.Equal("file1", root.Steps[0].In[0].ID, lineNumber())
	assert.Equal("file1", root.Steps[0].In[0].Source[0], lineNumber())
	assert.Equal("output", root.Steps[0].Out[0].ID, lineNumber())
}

func countLines4WfTest(assert *a.Assertions, root *Root) {
	assert.Equal("Workflow", root.Class, lineNumber())

	assert.Equal("file1", root.Inputs[0].ID, lineNumber())
	assert.Equal("File", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("file2", root.Inputs[1].ID, lineNumber())
	assert.Equal("File", root.Inputs[1].Types[0].Type, lineNumber())
	assert.Equal("count_output", root.Outputs[0].ID, lineNumber())
	assert.Equal("int[]", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal([]string{"step1/output"}, root.Outputs[0].Source, lineNumber())

	assert.Equal("ScatterFeatureRequirement", root.Requirements[0].Class, lineNumber())
	assert.Equal("MultipleInputFeatureRequirement", root.Requirements[1].Class, lineNumber())
	assert.Equal("step1", root.Steps[0].ID, lineNumber())
	assert.Equal("wc2-tool.cwl", root.Steps[0].Run.Value, lineNumber())
	assert.Equal("file1", root.Steps[0].Scatter[0], lineNumber())
	assert.Equal("file1", root.Steps[0].In[0].ID, lineNumber())
	assert.Equal("file1", root.Steps[0].In[0].Source[0], lineNumber())
	assert.Equal("file2", root.Steps[0].In[0].Source[1], lineNumber())
	assert.Equal("output", root.Steps[0].Out[0].ID, lineNumber())
}

func countLines5WfTest(assert *a.Assertions, root *Root) {
	assert.Equal("Workflow", root.Class, lineNumber())

	assert.Equal("file1", root.Inputs[0].ID, lineNumber())
	assert.Equal("File", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal(reflect.Map, root.Inputs[0].Default.Kind, lineNumber())
	assert.Equal("count_output", root.Outputs[0].ID, lineNumber())
	assert.Equal("int", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal([]string{"step1/output"}, root.Outputs[0].Source, lineNumber())

	assert.Equal("step1", root.Steps[0].ID, lineNumber())
	assert.Equal("wc2-tool.cwl", root.Steps[0].Run.Value, lineNumber())
	assert.Equal("file1", root.Steps[0].In[0].ID, lineNumber())
	assert.Equal("file1", root.Steps[0].In[0].Source[0], lineNumber())
	assert.Equal("output", root.Steps[0].Out[0].ID, lineNumber())
}

func countLines6WfTest(assert *a.Assertions, root *Root) {
	assert.Equal("Workflow", root.Class, lineNumber())

	assert.Equal("file1", root.Inputs[0].ID, lineNumber())
	assert.Equal("File[]", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("file2", root.Inputs[1].ID, lineNumber())
	assert.Equal("File[]", root.Inputs[1].Types[0].Type, lineNumber())
	assert.Equal("count_output", root.Outputs[0].ID, lineNumber())
	assert.Equal("int[]", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal([]string{"step1/output"}, root.Outputs[0].Source, lineNumber())

	assert.Equal("ScatterFeatureRequirement", root.Requirements[0].Class, lineNumber())
	assert.Equal("MultipleInputFeatureRequirement", root.Requirements[1].Class, lineNumber())

	assert.Equal("step1", root.Steps[0].ID, lineNumber())
	assert.Equal("wc3-tool.cwl", root.Steps[0].Run.Value, lineNumber())
	assert.Equal("file1", root.Steps[0].Scatter[0], lineNumber())
	assert.Equal("file1", root.Steps[0].In[0].ID, lineNumber())
	assert.Equal("file1", root.Steps[0].In[0].Source[0], lineNumber())
	assert.Equal("file2", root.Steps[0].In[0].Source[1], lineNumber())
	assert.Equal("merge_nested", root.Steps[0].In[0].LinkMerge, lineNumber())
	assert.Equal("output", root.Steps[0].Out[0].ID, lineNumber())
}

func countLines7WfTest(assert *a.Assertions, root *Root) {
	assert.Equal("Workflow", root.Class, lineNumber())

	assert.Equal("file1", root.Inputs[0].ID, lineNumber())
	assert.Equal("File[]", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("file2", root.Inputs[1].ID, lineNumber())
	assert.Equal("File[]", root.Inputs[1].Types[0].Type, lineNumber())
	assert.Equal("count_output", root.Outputs[0].ID, lineNumber())
	assert.Equal("int", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal([]string{"step1/output"}, root.Outputs[0].Source, lineNumber())

	assert.Equal("MultipleInputFeatureRequirement", root.Requirements[0].Class, lineNumber())

	assert.Equal("step1", root.Steps[0].ID, lineNumber())
	assert.Equal("wc3-tool.cwl", root.Steps[0].Run.Value, lineNumber())
	assert.Equal("file1", root.Steps[0].In[0].ID, lineNumber())
	assert.Equal("file1", root.Steps[0].In[0].Source[0], lineNumber())
	assert.Equal("file2", root.Steps[0].In[0].Source[1], lineNumber())
	assert.Equal("merge_flattened", root.Steps[0].In[0].LinkMerge, lineNumber())
	assert.Equal("output", root.Steps[0].Out[0].ID, lineNumber())
}

func countLines8WfTest(assert *a.Assertions, root *Root) {
	assert.Equal("Workflow", root.Class, lineNumber())

	assert.Equal("file1", root.Inputs[0].ID, lineNumber())
	assert.Equal("File", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("count_output", root.Outputs[0].ID, lineNumber())
	assert.Equal("int", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal([]string{"step1/count_output"}, root.Outputs[0].Source, lineNumber())

	assert.Equal("SubworkflowFeatureRequirement", root.Requirements[0].Class, lineNumber())

	assert.Equal("step1", root.Steps[0].ID, lineNumber())
	assert.Equal("count-lines1-wf.cwl", root.Steps[0].Run.Value, lineNumber())
	assert.Equal("file1", root.Steps[0].In[0].ID, lineNumber())
	assert.Equal("file1", root.Steps[0].In[0].Source[0], lineNumber())
	assert.Equal("count_output", root.Steps[0].Out[0].ID, lineNumber())
}

func countLines9WfTest(assert *a.Assertions, root *Root) {
	assert.Equal("Workflow", root.Class, lineNumber())

	assert.Equal(0, len(root.Inputs), lineNumber())

	assert.Equal("count_output", root.Outputs[0].ID, lineNumber())
	assert.Equal("int", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal([]string{"step2/output"}, root.Outputs[0].Source, lineNumber())

	assert.Equal("step1", root.Steps[0].ID, lineNumber())
	assert.Equal("wc-tool.cwl", root.Steps[0].Run.Value, lineNumber())
	assert.Equal("file1", root.Steps[0].In[0].ID, lineNumber())
	assert.Equal(reflect.Map, root.Steps[0].In[0].Default.Kind, lineNumber())
	assert.Equal("output", root.Steps[0].Out[0].ID, lineNumber())

	assert.Equal("step2", root.Steps[1].ID, lineNumber())
	assert.Equal("parseInt-tool.cwl", root.Steps[1].Run.Value, lineNumber())
	assert.Equal("file1", root.Steps[1].In[0].ID, lineNumber())
	assert.Equal("step1/output", root.Steps[1].In[0].Source[0], lineNumber())
	assert.Equal("output", root.Steps[1].Out[0].ID, lineNumber())
}

func defaultPathTest(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal(1, len(root.Inputs), lineNumber())
	assert.Equal("file1", root.Inputs[0].ID, lineNumber())
	assert.Equal("File", root.Inputs[0].Types[0].Type, lineNumber())
	// TODO support default: section
	// TODO support outputs: []
	assert.Equal(2, len(root.Arguments), lineNumber())
	assert.Equal("cat", root.Arguments[0].Value, lineNumber())
	assert.Equal("$(inputs.file1.path)", root.Arguments[1].Value, lineNumber())
}

func dir2Test(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())

	assert.Equal("DockerRequirement", root.Hints[0].Class, lineNumber())
	assert.Equal("debian:8", root.Hints[0].DockerPull, lineNumber())
	assert.Equal("ShellCommandRequirement", root.Hints[1].Class, lineNumber())
	assert.Equal("indir", root.Inputs[0].ID, lineNumber())
	assert.Equal("Directory", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("outlist", root.Outputs[0].ID, lineNumber())
	assert.Equal("File", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal([]string{"output.txt"}, root.Outputs[0].Binding.Glob, lineNumber())

	assert.Equal("cd", root.Arguments[0].Value, lineNumber())
	assert.Equal("$(inputs.indir.path)", root.Arguments[1].Value, lineNumber())
	assert.Equal(false, root.Arguments[2].Binding.ShellQuote, lineNumber())
	assert.Equal("&&", root.Arguments[2].Binding.ValueFrom.Key(), lineNumber())
	assert.Equal("find", root.Arguments[3].Value, lineNumber())
	assert.Equal(".", root.Arguments[4].Value, lineNumber())
	assert.Equal(false, root.Arguments[5].Binding.ShellQuote, lineNumber())
	assert.Equal("|", root.Arguments[5].Binding.ValueFrom.Key(), lineNumber())
	assert.Equal("sort", root.Arguments[6].Value, lineNumber())
	assert.Equal("output.txt", root.Stdout, lineNumber())
}

func dir3Test(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal("tar", root.BaseCommands[0], lineNumber())
	assert.Equal("xvf", root.BaseCommands[1], lineNumber())
	assert.Equal("inf", root.Inputs[0].ID, lineNumber())
	assert.Equal("File", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal(1, root.Inputs[0].Binding.Position, lineNumber())
	assert.Equal("outdir", root.Outputs[0].ID, lineNumber())
	assert.Equal("Directory", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal([]string{"."}, root.Outputs[0].Binding.Glob, lineNumber())
}

func dir4Test(assert *a.Assertions, root *Root) {
	assert.Equal("ShellCommandRequirement", root.Requirements[0].Class, lineNumber())
	assert.Equal("inf", root.Inputs[0].ID, lineNumber())
	assert.Equal("File", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("outlist", root.Outputs[0].ID, lineNumber())
	assert.Equal("File", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal([]string{"output.txt"}, root.Outputs[0].Binding.Glob, lineNumber())

	assert.Equal("cd", root.Arguments[0].Value, lineNumber())
	assert.Equal("$(inputs.inf.dirname)/xtestdir", root.Arguments[1].Value, lineNumber())
	assert.Equal(false, root.Arguments[2].Binding.ShellQuote, lineNumber())
	assert.Equal("&&", root.Arguments[2].Binding.ValueFrom.Key(), lineNumber())
	assert.Equal("find", root.Arguments[3].Value, lineNumber())
	assert.Equal(".", root.Arguments[4].Value, lineNumber())
	assert.Equal(false, root.Arguments[5].Binding.ShellQuote, lineNumber())
	assert.Equal("|", root.Arguments[5].Binding.ValueFrom.Key(), lineNumber())
	assert.Equal("sort", root.Arguments[6].Value, lineNumber())
	assert.Equal("output.txt", root.Stdout, lineNumber())
}

func dir5Test(assert *a.Assertions, root *Root) {
	assert.Equal("ShellCommandRequirement", root.Requirements[0].Class, lineNumber())
	assert.Equal("InitialWorkDirRequirement", root.Requirements[1].Class, lineNumber())
	assert.Equal("$(inputs.indir.listing)", root.Requirements[1].Listing[0].Location, lineNumber())
	assert.Equal("indir", root.Inputs[0].ID, lineNumber())
	assert.Equal("Directory", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("outlist", root.Outputs[0].ID, lineNumber())
	assert.Equal("File", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal([]string{"output.txt"}, root.Outputs[0].Binding.Glob, lineNumber())

	assert.Equal("find", root.Arguments[0].Value, lineNumber())
	assert.Equal("-L", root.Arguments[1].Value, lineNumber())
	assert.Equal(".", root.Arguments[2].Value, lineNumber())
	assert.Equal(false, root.Arguments[3].Binding.ShellQuote, lineNumber())
	assert.Equal("|", root.Arguments[3].Binding.ValueFrom.Key(), lineNumber())
	assert.Equal("sort", root.Arguments[4].Value, lineNumber())
	assert.Equal("output.txt", root.Stdout, lineNumber())
}

func dir6Test(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())

	assert.Equal("ShellCommandRequirement", root.Requirements[0].Class, lineNumber())
	assert.Equal("indir", root.Inputs[0].ID, lineNumber())
	assert.Equal("Directory", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal(-1, root.Inputs[0].Binding.Position, lineNumber())
	assert.Equal("cd", root.Inputs[0].Binding.Prefix, lineNumber())
	assert.Equal("outlist", root.Outputs[0].ID, lineNumber())
	assert.Equal("File", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal([]string{"output.txt"}, root.Outputs[0].Binding.Glob, lineNumber())

	assert.Equal(false, root.Arguments[0].Binding.ShellQuote, lineNumber())
	assert.Equal("&&", root.Arguments[0].Binding.ValueFrom.Key(), lineNumber())
	assert.Equal("find", root.Arguments[1].Value, lineNumber())
	assert.Equal(".", root.Arguments[2].Value, lineNumber())
	assert.Equal(false, root.Arguments[3].Binding.ShellQuote, lineNumber())
	assert.Equal("|", root.Arguments[3].Binding.ValueFrom.Key(), lineNumber())
	assert.Equal("sort", root.Arguments[4].Value, lineNumber())
	assert.Equal("output.txt", root.Stdout, lineNumber())
}

func dir7Test(assert *a.Assertions, root *Root) {
	assert.Equal("ExpressionTool", root.Class, lineNumber())

	assert.Equal("InlineJavascriptRequirement", root.Requirements[0].Class, lineNumber())

	assert.Equal("files", root.Inputs[0].ID, lineNumber())
	assert.Equal("File[]", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("dir", root.Outputs[0].ID, lineNumber())
	assert.Equal("Directory", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal(`${
return {"dir": {"class": "Directory", "basename": "a_directory", "listing": inputs.files}};
}`, root.Expression)
}

func dirTest(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())

	assert.Equal("ShellCommandRequirement", root.Requirements[0].Class, lineNumber())

	assert.Equal("Directory", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("outlist", root.Outputs[0].ID, lineNumber())
	assert.Equal("File", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal([]string{"output.txt"}, root.Outputs[0].Binding.Glob, lineNumber())

	assert.Equal("cd", root.Arguments[0].Value, lineNumber())
	assert.Equal("$(inputs.indir.path)", root.Arguments[1].Value, lineNumber())
	assert.Equal(false, root.Arguments[2].Binding.ShellQuote, lineNumber())
	assert.Equal("&&", root.Arguments[2].Binding.ValueFrom.Key(), lineNumber())
	assert.Equal("find", root.Arguments[3].Value, lineNumber())
	assert.Equal(".", root.Arguments[4].Value, lineNumber())
	assert.Equal(false, root.Arguments[5].Binding.ShellQuote, lineNumber())
	assert.Equal("|", root.Arguments[5].Binding.ValueFrom.Key(), lineNumber())
	assert.Equal("sort", root.Arguments[6].Value, lineNumber())
	assert.Equal("output.txt", root.Stdout, lineNumber())
}

func dockerArraySecondaryfilesTest(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal(3, len(root.Requirements), lineNumber())
	assert.Equal("DockerRequirement", root.Requirements[0].Class, lineNumber())
	assert.Equal("debian:8", root.Requirements[0].DockerPull, lineNumber())
	assert.Equal("InlineJavascriptRequirement", root.Requirements[1].Class, lineNumber())
	assert.Equal("ShellCommandRequirement", root.Requirements[2].Class, lineNumber())
	assert.Equal(1, len(root.Inputs), lineNumber())
	assert.Equal("fasta_path", root.Inputs[0].ID, lineNumber())
	assert.Equal("array", root.Inputs[0].Types[0].Type, lineNumber())
	// TODO items: File
	assert.Equal(".fai", root.Inputs[0].SecondaryFiles[0].Entry, lineNumber())
	assert.Equal("bai_list", root.Outputs[0].ID, lineNumber())
	assert.Equal("File", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("fai.list", root.Outputs[0].Binding.Glob[0], lineNumber())
	// TODO: Fix "Alias.Key()"
	assert.Equal(`{ var fai_list = ""; for (var i = 0; i < inputs.fasta_path.length; i ++) { fai_list += " cat " + inputs.fasta_path[i].path +".fai" + " >> fai.list && " } return fai_list.slice(0,-3) }`, root.Arguments[0].Binding.ValueFrom.Key(), lineNumber())
	assert.Equal(1, root.Arguments[0].Binding.Position, lineNumber())
	assert.Equal(false, root.Arguments[0].Binding.ShellQuote, lineNumber())
	assert.Equal(0, len(root.BaseCommands), lineNumber())
}

func dockerOutputDirTest(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal("DockerRequirement", root.Requirements[0].Class, lineNumber())
	assert.Equal("debian:8", root.Requirements[0].DockerPull, lineNumber())
	assert.Equal("/other", root.Requirements[0].DockerOutputDirectory, lineNumber())
	assert.Equal(0, len(root.Inputs), lineNumber())
	assert.Equal(1, len(root.Outputs), lineNumber())
	assert.Equal("thing", root.Outputs[0].ID, lineNumber())
	assert.Equal("File", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("thing", root.Outputs[0].Binding.Glob[0], lineNumber())
	assert.Equal(2, len(root.BaseCommands), lineNumber())
	assert.Equal("touch", root.BaseCommands[0], lineNumber())
	assert.Equal("/other/thing", root.BaseCommands[1], lineNumber())
}

func dynresreqTest(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal("ResourceRequirement", root.Requirements[0].Class, lineNumber())
	// TODO check CoresMin and CoresMax
	//assert.Equal("$(inputs.special_file.size)", root.Requirements[0].CoreMin, lineNumber())
	assert.Equal("special_file", root.Inputs[0].ID, lineNumber())
	assert.Equal("File", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("output", root.Outputs[0].ID, lineNumber())
	assert.Equal("stdout", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("echo", root.BaseCommands[0], lineNumber())
	assert.Equal("cores.txt", root.Stdout, lineNumber())
	assert.Equal("$(runtime.cores)", root.Arguments[0].Value, lineNumber())
}

func echoFileToolTest(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal("echo", root.BaseCommands[0], lineNumber())
	assert.Equal("in", root.Inputs[0].ID, lineNumber())
	assert.Equal("string", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal(1, root.Inputs[0].Binding.Position, lineNumber())
	// TODO check specification for this test ID and Type
	assert.Equal("out", root.Outputs[0].ID, lineNumber())
	assert.Equal("stdout", root.Outputs[0].Types[0].Type, lineNumber())
}

func echoToolTest(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal("in", root.Inputs[0].ID, lineNumber())
	assert.Equal("Any", root.Inputs[0].Types[0].Type, lineNumber())
	// TODO check specification for this test ID and Type
	assert.Equal("out", root.Outputs[0].ID, lineNumber())
	assert.Equal("string", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("out.txt", root.Outputs[0].Binding.Glob[0], lineNumber())
	assert.Equal(true, root.Outputs[0].Binding.LoadContents, lineNumber())
	assert.Equal("echo", root.BaseCommands[0], lineNumber())
	assert.Equal("out.txt", root.Stdout, lineNumber())
}

func envTool1Test(assert *a.Assertions, root *Root) {
	assert.Equal(3, len(root.BaseCommands), lineNumber())
	assert.Equal("/bin/bash", root.BaseCommands[0], lineNumber())
	assert.Equal("-c", root.BaseCommands[1], lineNumber())
	assert.Equal("echo $TEST_ENV", root.BaseCommands[2], lineNumber())
	assert.Equal(1, len(root.Inputs), lineNumber())
	// TODO ignore "in: string'
	assert.Equal(1, len(root.Outputs), lineNumber())
	assert.Equal("out", root.Outputs[0].ID, lineNumber())
	assert.Equal("File", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal([]string{"out"}, root.Outputs[0].Binding.Glob, lineNumber())
}

func envTool2Test(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal(1, len(root.Hints), lineNumber())
	assert.IsType(Hints{}, root.Hints)
	assert.Equal("EnvVarRequirement", root.Hints[0].Class, lineNumber())
	assert.Equal("TEST_ENV", root.Hints[0].Envs[0].Name, lineNumber())
	assert.Equal("$(inputs.in)", root.Hints[0].Envs[0].Value, lineNumber())
	assert.Equal(1, len(root.Inputs), lineNumber())
	// TODO in: string
	assert.Equal(1, len(root.Outputs), lineNumber())
	assert.Equal(3, len(root.BaseCommands), lineNumber())
	assert.Equal("/bin/bash", root.BaseCommands[0], lineNumber())
	assert.Equal("-c", root.BaseCommands[1], lineNumber())
	assert.Equal("echo $TEST_ENV", root.BaseCommands[2], lineNumber())
}

func envvar2Test(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal(0, len(root.Inputs), lineNumber())
	assert.Equal(0, len(root.Outputs), lineNumber())
	assert.Equal(1, len(root.Requirements), lineNumber())
	assert.Equal("ShellCommandRequirement", root.Requirements[0].Class, lineNumber())
	assert.Equal("DockerRequirement", root.Hints[0].Class, lineNumber())
	assert.Equal("debian:8", root.Hints[0].DockerPull, lineNumber())
	assert.Equal(12, len(root.Arguments), lineNumber())
	assert.Equal("echo", root.Arguments[0].Value, lineNumber())
	assert.Equal("\"HOME=$HOME\"", root.Arguments[1].Binding.ValueFrom.Key(), lineNumber())
	assert.Equal(false, root.Arguments[1].Binding.ShellQuote, lineNumber())
	assert.Equal("\"TMPDIR=$TMPDIR\"", root.Arguments[2].Binding.ValueFrom.Key(), lineNumber())
	assert.Equal(false, root.Arguments[2].Binding.ShellQuote, lineNumber())
	assert.Equal("&&", root.Arguments[3].Binding.ValueFrom.Key(), lineNumber())
	assert.Equal(false, root.Arguments[3].Binding.ShellQuote, lineNumber())
	assert.Equal("test", root.Arguments[4].Value, lineNumber())
	assert.Equal("\"$HOME\"", root.Arguments[5].Binding.ValueFrom.Key(), lineNumber())
	assert.Equal(false, root.Arguments[5].Binding.ShellQuote, lineNumber())
	assert.Equal("=", root.Arguments[6].Value, lineNumber())
	assert.Equal("$(runtime.outdir)", root.Arguments[7].Value, lineNumber())
	assert.Equal("-a", root.Arguments[8].Value, lineNumber())
	assert.Equal("\"$TMPDIR\"", root.Arguments[9].Binding.ValueFrom.Key(), lineNumber())
	assert.Equal(false, root.Arguments[9].Binding.ShellQuote, lineNumber())
	assert.Equal("=", root.Arguments[10].Value, lineNumber())
	assert.Equal("$(runtime.tmpdir)", root.Arguments[11].Value, lineNumber())
}

func envvarTest(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal(0, len(root.Inputs), lineNumber())
	assert.Equal(0, len(root.Outputs), lineNumber())
	assert.Equal(1, len(root.Requirements), lineNumber())
	assert.Equal("ShellCommandRequirement", root.Requirements[0].Class, lineNumber())
	assert.Equal(12, len(root.Arguments), lineNumber())
	assert.Equal("echo", root.Arguments[0].Value, lineNumber())
	assert.Equal("\"HOME=$HOME\"", root.Arguments[1].Binding.ValueFrom.Key(), lineNumber())
	assert.Equal(false, root.Arguments[1].Binding.ShellQuote, lineNumber())
	assert.Equal("\"TMPDIR=$TMPDIR\"", root.Arguments[2].Binding.ValueFrom.Key(), lineNumber())
	assert.Equal(false, root.Arguments[2].Binding.ShellQuote, lineNumber())
	assert.Equal("&&", root.Arguments[3].Binding.ValueFrom.Key(), lineNumber())
	assert.Equal(false, root.Arguments[3].Binding.ShellQuote, lineNumber())
	assert.Equal("test", root.Arguments[4].Value, lineNumber())
	assert.Equal("\"$HOME\"", root.Arguments[5].Binding.ValueFrom.Key(), lineNumber())
	assert.Equal(false, root.Arguments[5].Binding.ShellQuote, lineNumber())
	assert.Equal("=", root.Arguments[6].Value, lineNumber())
	assert.Equal("$(runtime.outdir)", root.Arguments[7].Value, lineNumber())
	assert.Equal("-a", root.Arguments[8].Value, lineNumber())
	assert.Equal("\"$TMPDIR\"", root.Arguments[9].Binding.ValueFrom.Key(), lineNumber())
	assert.Equal(false, root.Arguments[9].Binding.ShellQuote, lineNumber())
	assert.Equal("=", root.Arguments[10].Value, lineNumber())
	assert.Equal("$(runtime.tmpdir)", root.Arguments[11].Value, lineNumber())
}

func envWf1Test(assert *a.Assertions, root *Root) {
	assert.Equal("Workflow", root.Class, lineNumber())
	assert.Equal(1, len(root.Inputs), lineNumber())
	assert.Equal("string", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal(1, len(root.Outputs), lineNumber())
	assert.Equal("out", root.Outputs[0].ID, lineNumber())
	assert.Equal("File", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal([]string{"step1/out"}, root.Outputs[0].Source, lineNumber())
	assert.Equal(1, len(root.Requirements), lineNumber())
	assert.Equal("EnvVarRequirement", root.Requirements[0].Class, lineNumber())
	assert.Equal("TEST_ENV", root.Requirements[0].EnvDef[0].Name, lineNumber())
	assert.Equal(`override`, root.Requirements[0].EnvDef[0].Value, lineNumber())
}

func envWf2Test(assert *a.Assertions, root *Root) {
	// this cwl is almost identical to the one for envWf1
	envWf1Test(assert, root)
}

func envWf3Test(assert *a.Assertions, root *Root) {
	assert.Equal("Workflow", root.Class, lineNumber())
	assert.Equal(1, len(root.Inputs), lineNumber())
	assert.Equal("string", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal(1, len(root.Outputs), lineNumber())
	assert.Equal("out", root.Outputs[0].ID, lineNumber())
	assert.Equal("File", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal([]string{"step1/out"}, root.Outputs[0].Source, lineNumber())
	assert.Equal(0, len(root.Requirements), lineNumber())
	assert.Equal(1, len(root.Steps), lineNumber())
	assert.Equal("EnvVarRequirement", root.Steps[0].Requirements[0].Class, lineNumber())
	assert.Equal("TEST_ENV", root.Steps[0].Requirements[0].EnvDef[0].Name, lineNumber())
	assert.Equal(`override`, root.Steps[0].Requirements[0].EnvDef[0].Value, lineNumber())
}

func fileLiteralExTest(assert *a.Assertions, root *Root) {
	assert.Equal("ExpressionTool", root.Class, lineNumber())
	assert.Equal(0, len(root.Inputs), lineNumber())
	// TODO check specification for this test ID and Type
	assert.Equal("lit", root.Outputs[0].ID, lineNumber())
	assert.Equal("File", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal(`${
return {"lit": {"class": "File", "basename": "a_file", "contents": "Hello file literal."}};
}`, root.Expression)
}

func formattest2Test(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal(1, len(root.Inputs), lineNumber())
	assert.Equal("input", root.Inputs[0].ID, lineNumber())
	assert.Equal("File", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("edam:format_2330", root.Inputs[0].Format, lineNumber())
	assert.Equal(0, root.Inputs[0].Binding.Position, lineNumber())
	assert.Equal(1, len(root.Outputs), lineNumber())
	assert.Equal("output", root.Outputs[0].ID, lineNumber())
	assert.Equal("File", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal([]string{"output.txt"}, root.Outputs[0].Binding.Glob, lineNumber())
	assert.Equal("$(inputs.input.format)", root.Outputs[0].Format, lineNumber())
	assert.Equal(1, len(root.BaseCommands), lineNumber())
	assert.Equal("rev", root.BaseCommands[0], lineNumber())
	assert.Equal("output.txt", root.Stdout, lineNumber())
}

func formattest3Test(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	// $namespaces
	assert.Equal(2, len(root.Namespaces), lineNumber())
	assert.Equal("http://edamontology.org/", root.Namespaces[0]["edam"], lineNumber())
	assert.Equal("http://galaxyproject.org/formats/", root.Namespaces[1]["gx"], lineNumber())
	// $schemas
	assert.Equal(2, len(root.Schemas), lineNumber())
	assert.Equal("EDAM.owl", root.Schemas[0], lineNumber())
	assert.Equal("gx_edam.ttl", root.Schemas[1], lineNumber())
	assert.Equal("Reverse each line using the `rev` command", root.Doc, lineNumber())
	assert.Equal("DockerRequirement", root.Hints[0].Class, lineNumber())
	assert.Equal("debian:wheezy", root.Hints[0].DockerPull, lineNumber())
	assert.Equal(1, len(root.Inputs), lineNumber())
	assert.Equal("input", root.Inputs[0].ID, lineNumber())
	assert.Equal("File", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("gx:fasta", root.Inputs[0].Format, lineNumber())
	assert.Equal(0, root.Inputs[0].Binding.Position, lineNumber())
	assert.Equal(1, len(root.Outputs), lineNumber())
	assert.Equal("output", root.Outputs[0].ID, lineNumber())
	assert.Equal("File", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal([]string{"output.txt"}, root.Outputs[0].Binding.Glob, lineNumber())
	assert.Equal("$(inputs.input.format)", root.Outputs[0].Format, lineNumber())
	assert.Equal(1, len(root.BaseCommands), lineNumber())
	assert.Equal("rev", root.BaseCommands[0], lineNumber())
	assert.Equal("output.txt", root.Stdout, lineNumber())
}

func formattestTest(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal(1, len(root.Inputs), lineNumber())
	assert.Equal("input", root.Inputs[0].ID, lineNumber())
	assert.Equal("File", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("edam:format_2330", root.Inputs[0].Format, lineNumber())
	assert.Equal(0, root.Inputs[0].Binding.Position, lineNumber())
	assert.Equal(1, len(root.Outputs), lineNumber())
	assert.Equal("output", root.Outputs[0].ID, lineNumber())
	assert.Equal("File", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal([]string{"output.txt"}, root.Outputs[0].Binding.Glob, lineNumber())
	assert.Equal("edam:format_2330", root.Outputs[0].Format, lineNumber())
	assert.Equal(1, len(root.BaseCommands), lineNumber())
	assert.Equal("rev", root.BaseCommands[0], lineNumber())
	assert.Equal("output.txt", root.Stdout, lineNumber())
}

func globExprListTest(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal(1, len(root.Inputs), lineNumber())
	assert.Equal("ids", root.Inputs[0].ID, lineNumber())
	assert.Equal("string[]", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal(1, root.Inputs[0].Binding.Position, lineNumber())
	assert.Equal(1, len(root.Outputs), lineNumber())
	assert.Equal("files", root.Outputs[0].ID, lineNumber())
	assert.Equal("File[]", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal([]string{"$(inputs.ids)"}, root.Outputs[0].Binding.Glob, lineNumber())
	assert.Equal(1, len(root.BaseCommands), lineNumber())
	assert.Equal("touch", root.BaseCommands[0], lineNumber())
}

func importedHintTest(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal(0, len(root.Inputs), lineNumber())
	assert.Equal(1, len(root.Outputs), lineNumber())
	assert.Equal("out", root.Outputs[0].ID, lineNumber())
	// TODO test out: stdout 's stdout
	assert.IsType(Hints{}, root.Hints)
	assert.Equal("envvar.yml", root.Hints[0].Import, lineNumber())
	assert.Equal(3, len(root.BaseCommands), lineNumber())
	assert.Equal("/bin/bash", root.BaseCommands[0], lineNumber())
	assert.Equal("-c", root.BaseCommands[1], lineNumber())
	assert.Equal("echo $TEST_ENV", root.BaseCommands[2], lineNumber())
	assert.Equal("out", root.Stdout, lineNumber())
}

func initialworkdirrequirementDockerOutTest(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal(1, len(root.Inputs), lineNumber())
	assert.Equal("INPUT", root.Inputs[0].ID, lineNumber())
	assert.Equal("File", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal(1, len(root.Outputs), lineNumber())
	assert.Equal("OUTPUT", root.Outputs[0].ID, lineNumber())
	assert.Equal("File", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal([]string{"$(inputs.INPUT.basename)"}, root.Outputs[0].Binding.Glob, lineNumber())
	assert.Equal(".fai", root.Outputs[0].SecondaryFiles[0].Entry, lineNumber())
	// TODO outputs
	assert.Equal(2, len(root.Requirements), lineNumber())
	assert.Equal("DockerRequirement", root.Requirements[0].Class, lineNumber())
	assert.Equal("debian:8", root.Requirements[0].DockerPull, lineNumber())
	assert.Equal("InitialWorkDirRequirement", root.Requirements[1].Class, lineNumber())
	assert.Equal("$(inputs.INPUT)", root.Requirements[1].Listing[0].Location, lineNumber())
	// TODO: fix "Alias.Key()"
	assert.Equal("inputs.INPUT.basename).fai", root.Arguments[0].Binding.ValueFrom.Key(), lineNumber())
	// TODO test against "position" but currently just put 0 is failed
	assert.Equal(1, len(root.BaseCommands), lineNumber())
	assert.Equal("touch", root.BaseCommands[0], lineNumber())
}

func initialworkPathTest(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal(1, len(root.Inputs), lineNumber())
	assert.Equal("file1", root.Inputs[0].ID, lineNumber())
	assert.Equal("File", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal(0, len(root.Outputs), lineNumber())

	assert.Equal("InitialWorkDirRequirement", root.Requirements[0].Class, lineNumber())
	assert.Equal("bob.txt", root.Requirements[0].Listing[0].EntryName, lineNumber())
	assert.Equal(`$(inputs.file1)`, root.Requirements[0].Listing[0].Entry, lineNumber())
	assert.Equal("ShellCommandRequirement", root.Requirements[1].Class, lineNumber())

	assert.Equal(false, root.Arguments[0].Binding.ShellQuote, lineNumber())
	assert.Equal(`test "$(inputs.file1.path)" = "$(runtime.outdir)/bob.txt"
`, root.Arguments[0].Binding.ValueFrom.Key(), lineNumber())
	// TODO write basecommand
}

func inlineJsTest(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	// TODO test BaseCommand because this file has two baseCommand fields
	//fmt.Println(root.BaseCommands)
	//assert.Equal(0, len(root.BaseCommands), lineNumber())
	//assert.Equal("touch", root.BaseCommands[0], lineNumber())
	assert.Equal(1, len(root.Requirements), lineNumber())
	assert.Equal("InlineJavascriptRequirement", root.Requirements[0].Class, lineNumber())
	assert.Equal(1, len(root.Requirements), lineNumber())
	assert.Equal("DockerRequirement", root.Hints[0].Class, lineNumber())
	assert.Equal("python:2-slim", root.Hints[0].DockerPull, lineNumber())
	assert.Equal(1, len(root.Inputs), lineNumber())
	assert.Equal("args.py", root.Inputs[0].ID, lineNumber())
	assert.Equal("File", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal(reflect.Map, root.Inputs[0].Default.Kind, lineNumber())
	assert.Equal(-1, root.Inputs[0].Binding.Position, lineNumber())
	assert.Equal(1, len(root.Outputs), lineNumber())
	assert.Equal("args", root.Outputs[0].ID, lineNumber())
	assert.Equal("array", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("string", root.Outputs[0].Types[0].Items[0].Type, lineNumber())
	assert.Equal(3, len(root.Arguments), lineNumber())
	assert.Equal("-A", root.Arguments[0].Binding.Prefix, lineNumber())
	// {{{ TODO: Fix "Alias.Key()"
	assert.Equal("1+1", root.Arguments[0].Binding.ValueFrom.Key(), lineNumber())
	assert.Equal("-B", root.Arguments[1].Binding.Prefix, lineNumber())
	assert.Equal(`"/foo/bar/baz".split('/').slice(-1)[0]`, root.Arguments[1].Binding.ValueFrom.Key(), lineNumber())
	assert.Equal("-C", root.Arguments[2].Binding.Prefix, lineNumber())
	assert.Equal(`{
  var r = [];
  for (var i = 10; i >= 1; i--) {
    r.push(i);
  }
  return r;
}
`, root.Arguments[2].Binding.ValueFrom.Key())
}

func jsExprReqWfTest(assert *a.Assertions, root *Root) {
	assert.Equal(2, len(root.Graphs), lineNumber())
	// 0
	assert.Equal("tool", root.Graphs[0].ID, lineNumber())
	assert.Equal("CommandLineTool", root.Graphs[0].Class, lineNumber())
	assert.Equal("InlineJavascriptRequirement", root.Graphs[0].Requirements[0].Class, lineNumber())
	assert.Equal("function foo() { return 2; }", root.Graphs[0].Requirements[0].ExpressionLib[0].Value, lineNumber())
	assert.Equal(0, len(root.Graphs[0].Inputs), lineNumber())
	assert.Equal("echo", root.Graphs[0].Arguments[0].Value, lineNumber())
	assert.Equal("whatever.txt", root.Graphs[0].Stdout, lineNumber())
	assert.Equal(1, len(root.Graphs[0].Outputs), lineNumber())
	assert.Equal("out", root.Graphs[0].Outputs[0].ID, lineNumber())
	assert.Equal("stdout", root.Graphs[0].Outputs[0].Types[0].Type, lineNumber())
	// 1
	assert.Equal("wf", root.Graphs[1].ID, lineNumber())
	assert.Equal("Workflow", root.Graphs[1].Class, lineNumber())
	assert.Equal("InlineJavascriptRequirement", root.Graphs[1].Requirements[0].Class, lineNumber())
	assert.Equal("function bar() { return 1; }", root.Graphs[1].Requirements[0].ExpressionLib[0].Value, lineNumber())
	assert.Equal(0, len(root.Graphs[1].Inputs), lineNumber())
	assert.Equal("out", root.Graphs[1].Outputs[0].ID, lineNumber())
	assert.Equal("File", root.Graphs[1].Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("tool/out", root.Graphs[1].Outputs[0].Source[0], lineNumber())
	assert.Equal("tool", root.Graphs[1].Steps[0].ID, lineNumber())
	// assert.Equal("#tool", root.Graphs[1].Steps[0].Run.Workflow.ID, lineNumber())
	assert.Equal(0, len(root.Graphs[1].Steps[0].In), lineNumber())
	// TODO check empty In
	assert.Equal(1, len(root.Graphs[1].Steps[0].Out), lineNumber())
	assert.Equal("out", root.Graphs[1].Steps[0].Out[0].ID, lineNumber())
}

func metadataTest(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal("Print the contents of a file to stdout using 'cat' running in a docker container.", root.Doc, lineNumber())
	assert.Equal(1, len(root.Hints), lineNumber())
	assert.IsType(Hints{}, root.Hints)
	assert.Equal("DockerRequirement", root.Hints[0].Class, lineNumber())
	assert.Equal("debian:wheezy", root.Hints[0].DockerPull, lineNumber())

	assert.Equal(2, len(root.Inputs), lineNumber())
	assert.Equal("numbering", root.Inputs[0].ID, lineNumber())
	assert.Equal("file1", root.Inputs[1].ID, lineNumber())
	assert.Equal("File", root.Inputs[1].Types[0].Type, lineNumber())
	assert.Equal(1, root.Inputs[1].Binding.Position, lineNumber())
	assert.Equal(0, len(root.Outputs), lineNumber())
	assert.Equal(1, len(root.BaseCommands), lineNumber())
	assert.Equal("cat", root.BaseCommands[0], lineNumber())
	// $namespaces
	assert.Equal(2, len(root.Namespaces), lineNumber())
	assert.Equal("http://purl.org/dc/terms/", root.Namespaces[0]["dct"], lineNumber())
	assert.Equal("http://xmlns.com/foaf/0.1/", root.Namespaces[1]["foaf"], lineNumber())
	// $schemas
	assert.Equal(2, len(root.Schemas), lineNumber())
	assert.Equal("foaf.rdf", root.Schemas[0], lineNumber())
	assert.Equal("dcterms.rdf", root.Schemas[1], lineNumber())
}

func namerootTest(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal(1, len(root.Inputs), lineNumber())
	assert.Equal("file1", root.Inputs[0].ID, lineNumber())
	assert.Equal("File", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal(1, len(root.Outputs), lineNumber())
	assert.Equal("b", root.Outputs[0].ID, lineNumber())
	assert.Equal(0, len(root.BaseCommands), lineNumber())
	assert.Equal(4, len(root.Arguments), lineNumber())
	assert.Equal("echo", root.Arguments[0].Value, lineNumber())
	assert.Equal("$(inputs.file1.basename)", root.Arguments[1].Value, lineNumber())
	assert.Equal("$(inputs.file1.nameroot)", root.Arguments[2].Value, lineNumber())
	assert.Equal("$(inputs.file1.nameext)", root.Arguments[3].Value, lineNumber())
}

func nestedArrayTest(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal("echo", root.BaseCommands[0], lineNumber())
	assert.Equal(1, len(root.Inputs), lineNumber())
	assert.Equal("letters", root.Inputs[0].ID, lineNumber())
	assert.Equal("array", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("array", root.Inputs[0].Types[0].Items[0].Type, lineNumber())
	assert.Equal("string", root.Inputs[0].Types[0].Items[0].Items[0].Type, lineNumber())
	assert.Equal(1, root.Inputs[0].Binding.Position, lineNumber())
	assert.Equal("echo.txt", root.Stdout, lineNumber())
	assert.Equal("echo", root.Outputs[0].ID, lineNumber())
	assert.Equal("stdout", root.Outputs[0].Types[0].Type, lineNumber())
}

func nullDefinedTest(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal(1, len(root.Requirements), lineNumber())
	assert.Equal("InlineJavascriptRequirement", root.Requirements[0].Class, lineNumber())
	assert.Equal(1, len(root.Inputs), lineNumber())
	assert.Equal("file1", root.Inputs[0].ID, lineNumber())
	assert.Equal("File?", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("out", root.Outputs[0].ID, lineNumber())
	assert.Equal("string", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("out.txt", root.Outputs[0].Binding.Glob[0], lineNumber())
	assert.Equal(false, root.Outputs[0].Binding.Contents, lineNumber())
	assert.Equal("$(self[0].contents)", root.Outputs[0].Binding.Eval, lineNumber())
	assert.Equal("out.txt", root.Stdout, lineNumber())
	assert.Equal(2, len(root.Arguments), lineNumber())
	assert.Equal("echo", root.Arguments[0].Value, lineNumber())
	assert.Equal(`$(inputs.file1 === null ? "t" : "f")`, root.Arguments[1].Value, lineNumber())
}

func nullExpression1ToolTest(assert *a.Assertions, root *Root) {
	assert.Equal("ExpressionTool", root.Class, lineNumber())
	assert.Equal(1, len(root.Requirements), lineNumber())
	assert.Equal("InlineJavascriptRequirement", root.Requirements[0].Class, lineNumber())
	assert.Equal(1, len(root.Inputs), lineNumber())
	assert.Equal("i1", root.Inputs[0].ID, lineNumber())
	assert.Equal("Any", root.Inputs[0].Types[0].Type, lineNumber())
	// TODO input default
	//assert.Equal("File", root.Inputs[0].Default.Class, lineNumber())
	//fmt.Println(t, root.Inputs[0].Default)
	assert.Equal("output", root.Outputs[0].ID, lineNumber())
	assert.Equal("int", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal(`$({'output': (inputs.i1 == 'the-default' ? 1 : 2)})`, root.Expression, lineNumber())
}

func nullExpression2ToolTest(assert *a.Assertions, root *Root) {
	assert.Equal("ExpressionTool", root.Class, lineNumber())
	assert.Equal(1, len(root.Requirements), lineNumber())
	assert.Equal("InlineJavascriptRequirement", root.Requirements[0].Class, lineNumber())
	assert.Equal(1, len(root.Inputs), lineNumber())
	assert.Equal("i1", root.Inputs[0].ID, lineNumber())
	assert.Equal("Any", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("output", root.Outputs[0].ID, lineNumber())
	assert.Equal("int", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal(`$({'output': (inputs.i1 == 'the-default' ? 1 : 2)})`, root.Expression, lineNumber())
}

func optionalOutputTest(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal("Print the contents of a file to stdout using 'cat' running in a docker container.", root.Doc, lineNumber())
	assert.Equal("DockerRequirement", root.Hints[0].Class, lineNumber())
	assert.Equal("debian:wheezy", root.Hints[0].DockerPull, lineNumber())
	assert.Equal("file1", root.Inputs[0].ID, lineNumber())
	assert.Equal("File", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("Input File", root.Inputs[0].Label, lineNumber())
	assert.Equal("The file that will be copied using 'cat'", root.Inputs[0].Doc, lineNumber())
	assert.Equal(1, root.Inputs[0].Binding.Position, lineNumber())

	assert.Equal("optional_file", root.Outputs[0].ID, lineNumber())
	assert.Equal("File?", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("bumble.txt", root.Outputs[0].Binding.Glob[0], lineNumber())
	assert.Equal("output_file", root.Outputs[1].ID, lineNumber())
	assert.Equal("File", root.Outputs[1].Types[0].Type, lineNumber())
	assert.Equal("output.txt", root.Outputs[1].Binding.Glob[0], lineNumber())
	assert.Equal(".idx", root.Outputs[1].SecondaryFiles[0].Entry, lineNumber())

	assert.Equal("cat", root.BaseCommands[0], lineNumber())
	assert.Equal("output.txt", root.Stdout, lineNumber())
}

func params2Test(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal("bar", root.Inputs[0].ID, lineNumber())
	assert.Equal("Any", root.Inputs[0].Types[0].Type, lineNumber())
	// TODO check specification for this test ID and Type
	assert.Equal("$import", root.Outputs[0].ID, lineNumber())
	assert.Equal("params_inc.yml", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("true", root.BaseCommands[0], lineNumber())
}

func paramsTest(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal("bar", root.Inputs[0].ID, lineNumber())
	assert.Equal("Any", root.Inputs[0].Types[0].Type, lineNumber())
	// TODO check specification for this test ID and Type
	assert.Equal("$import", root.Outputs[0].ID, lineNumber())
	assert.Equal("params_inc.yml", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("true", root.BaseCommands[0], lineNumber())
}

func parseIntToolTest(assert *a.Assertions, root *Root) {
	assert.Equal("ExpressionTool", root.Class, lineNumber())
	assert.Equal(1, len(root.Requirements), lineNumber())
	assert.IsType(Requirements{}, root.Requirements)
	assert.Equal("InlineJavascriptRequirement", root.Requirements[0].Class, lineNumber())
	assert.Equal(1, len(root.Inputs), lineNumber())
	assert.Equal("file1", root.Inputs[0].ID, lineNumber())
	assert.Equal("File", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal(true, root.Inputs[0].Binding.LoadContents, lineNumber())
	assert.Equal(1, len(root.Outputs), lineNumber())
	assert.Equal("output", root.Outputs[0].ID, lineNumber())
	assert.Equal("int", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("$({'output': parseInt(inputs.file1.contents)})", root.Expression, lineNumber())
}

func recordOutputTest(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal("irec", root.Inputs[0].ID, lineNumber())
	assert.Equal("ShellCommandRequirement", root.Requirements[0].Class, lineNumber())
	assert.Equal("record", root.Inputs[0].Types[0].Type, lineNumber())
	// TODO type allows name ?
	assert.Equal("ifoo", root.Inputs[0].Types[0].Fields[0].Name, lineNumber())
	assert.Equal("File", root.Inputs[0].Types[0].Fields[0].Types[0].Type, lineNumber())
	assert.Equal(2, root.Inputs[0].Types[0].Fields[0].Binding.Position, lineNumber())
	assert.Equal("ibar", root.Inputs[0].Types[0].Fields[1].Name, lineNumber())
	assert.Equal("File", root.Inputs[0].Types[0].Fields[1].Types[0].Type, lineNumber())
	assert.Equal(6, root.Inputs[0].Types[0].Fields[1].Binding.Position, lineNumber())
	assert.Equal("orec", root.Outputs[0].ID, lineNumber())
	assert.Equal("record", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("ofoo", root.Outputs[0].Types[0].Fields[0].Name, lineNumber())
	assert.Equal("File", root.Outputs[0].Types[0].Fields[0].Types[0].Type, lineNumber())
	assert.Equal("foo", root.Outputs[0].Types[0].Fields[0].Binding.Glob[0], lineNumber())
	assert.Equal("obar", root.Outputs[0].Types[0].Fields[1].Name, lineNumber())
	assert.Equal("File", root.Outputs[0].Types[0].Fields[1].Types[0].Type, lineNumber())
	assert.Equal("bar", root.Outputs[0].Types[0].Fields[1].Binding.Glob[0], lineNumber())
	assert.Equal("cat", root.Arguments[0].Binding.ValueFrom.Key(), lineNumber())
	assert.Equal(1, root.Arguments[0].Binding.Position, lineNumber())
	assert.Equal("> foo", root.Arguments[1].Binding.ValueFrom.Key(), lineNumber())
	assert.Equal(3, root.Arguments[1].Binding.Position, lineNumber())
	assert.Equal(false, root.Arguments[1].Binding.ShellQuote, lineNumber())
	assert.Equal("&&", root.Arguments[2].Binding.ValueFrom.Key(), lineNumber())
	assert.Equal(4, root.Arguments[2].Binding.Position, lineNumber())
	assert.Equal(false, root.Arguments[2].Binding.ShellQuote, lineNumber())
	assert.Equal("cat", root.Arguments[3].Binding.ValueFrom.Key(), lineNumber())
	assert.Equal(5, root.Arguments[3].Binding.Position, lineNumber())
	assert.Equal("> bar", root.Arguments[4].Binding.ValueFrom.Key(), lineNumber())
	assert.Equal(7, root.Arguments[4].Binding.Position, lineNumber())
	assert.Equal(false, root.Arguments[4].Binding.ShellQuote, lineNumber())
}

func recursiveInputDirectoryTest(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal(3, len(root.Requirements), lineNumber())
	assert.Equal("InlineJavascriptRequirement", root.Requirements[0].Class, lineNumber())
	assert.Equal("InitialWorkDirRequirement", root.Requirements[1].Class, lineNumber())
	assert.Equal("$(inputs.input_dir)", root.Requirements[1].Listing[0].Entry, lineNumber())
	assert.Equal("work_dir", root.Requirements[1].Listing[0].EntryName, lineNumber())
	assert.Equal(true, root.Requirements[1].Listing[0].Writable, lineNumber())
	assert.Equal("ShellCommandRequirement", root.Requirements[2].Class, lineNumber())
	assert.Equal("output.txt", root.Stdout, lineNumber())
	assert.Equal(false, root.Arguments[0].Binding.ShellQuote, lineNumber())
	assert.Equal(`touch work_dir/e;
if [ ! -w work_dir ]; then echo work_dir not writable; fi;
if [ -L work_dir ]; then echo work_dir is a symlink; fi;
if [ ! -w work_dir/a ]; then echo work_dir/a not writable; fi;
if [ -L work_dir/a ]; then echo work_dir/a is a symlink; fi;
if [ ! -w work_dir/c ]; then echo work_dir/c not writable; fi;
if [ -L work_dir/c ]; then echo work_dir/c is a symlink; fi;
if [ ! -w work_dir/c/d ]; then echo work_dir/c/d not writable; fi;
if [ -L work_dir/c/d ]; then echo work_dir/c/d is a symlink; fi;
if [ ! -w work_dir/e ]; then echo work_dir/e not writable; fi;
if [ -L work_dir/e ]; then echo work_dir/e is a symlink ; fi;
`, root.Arguments[0].Binding.ValueFrom.Key())
	assert.Equal("input_dir", root.Inputs[0].ID, lineNumber())
	assert.Equal("Directory", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal(2, len(root.Outputs), lineNumber())

	assert.Equal("output_dir", root.Outputs[0].ID, lineNumber())
	assert.Equal("Directory", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("work_dir", root.Outputs[0].Binding.Glob[0], lineNumber())
	assert.Equal("test_result", root.Outputs[1].ID, lineNumber())
	assert.Equal("stdout", root.Outputs[1].Types[0].Type, lineNumber())
}

func renameTest(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal("true", root.BaseCommands[0], lineNumber())
	assert.Equal(1, len(root.Requirements), lineNumber())
	assert.Equal("InitialWorkDirRequirement", root.Requirements[0].Class, lineNumber())
	assert.Equal("$(inputs.newname)", root.Requirements[0].Listing[0].EntryName, lineNumber())
	assert.Equal(`$(inputs.srcfile)`, root.Requirements[0].Listing[0].Entry, lineNumber())
	assert.Equal(2, len(root.Inputs), lineNumber())

	assert.Equal("newname", root.Inputs[0].ID, lineNumber())
	assert.Equal("string", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("srcfile", root.Inputs[1].ID, lineNumber())
	assert.Equal("File", root.Inputs[1].Types[0].Type, lineNumber())
	assert.Equal(1, len(root.Outputs), lineNumber())
	assert.Equal("outfile", root.Outputs[0].ID, lineNumber())
	assert.Equal("File", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("$(inputs.newname)", root.Outputs[0].Binding.Glob[0], lineNumber())
}

func revsortPackedTest(assert *a.Assertions, root *Root) {
	assert.Equal(3, len(root.Graphs), lineNumber())
	// Graph 0
	assert.Equal("Workflow", root.Graphs[0].Class, lineNumber())
	assert.Equal("#main", root.Graphs[0].ID, lineNumber())
	assert.Equal("Reverse the lines in a document, then sort those lines.", root.Graphs[0].Doc, lineNumber())
	assert.Equal("DockerRequirement", root.Graphs[0].Hints[0].Class, lineNumber())
	assert.Equal("debian:8", root.Graphs[0].Hints[0].DockerPull, lineNumber())
	assert.Equal(2, len(root.Graphs[0].Inputs), lineNumber())
	assert.Equal("File", root.Graphs[0].Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("#main/input", root.Graphs[0].Inputs[0].ID, lineNumber())
	assert.Equal("The input file to be processed.", root.Graphs[0].Inputs[0].Doc, lineNumber())
	assert.Equal("boolean", root.Graphs[0].Inputs[1].Types[0].Type, lineNumber())
	assert.Equal("#main/reverse_sort", root.Graphs[0].Inputs[1].ID, lineNumber())
	assert.Equal("If true, reverse (decending) sort", root.Graphs[0].Inputs[1].Doc, lineNumber())
	assert.Equal(1, len(root.Graphs[0].Outputs), lineNumber())
	assert.Equal("#main/output", root.Graphs[0].Outputs[0].ID, lineNumber())
	assert.Equal("File", root.Graphs[0].Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("The output with the lines reversed and sorted.", root.Graphs[0].Outputs[0].Doc[0], lineNumber())
	assert.Equal("#main/sorted/output", root.Graphs[0].Outputs[0].Source[0], lineNumber())
	assert.Equal("#main/rev", root.Graphs[0].Steps[0].ID, lineNumber())
	assert.Equal("#revtool.cwl", root.Graphs[0].Steps[0].Run.Value, lineNumber())
	assert.Equal(1, len(root.Graphs[0].Steps[0].In), lineNumber())
	assert.Equal("#main/input", root.Graphs[0].Steps[0].In[0].Source[0], lineNumber())
	assert.Equal("#main/rev/input", root.Graphs[0].Steps[0].In[0].ID, lineNumber())
	assert.Equal("#main/rev/output", root.Graphs[0].Steps[0].Out[0].ID, lineNumber())
	assert.Equal("#main/sorted/output", root.Graphs[0].Steps[1].Out[0].ID, lineNumber())
	assert.Equal("#main/sorted", root.Graphs[0].Steps[1].ID, lineNumber())
	assert.Equal("#sorttool.cwl", root.Graphs[0].Steps[1].Run.Value, lineNumber())
	assert.Equal(2, len(root.Graphs[0].Steps[1].In), lineNumber())
	assert.Equal("#main/rev/output", root.Graphs[0].Steps[1].In[0].Source[0], lineNumber())
	assert.Equal("#main/sorted/input", root.Graphs[0].Steps[1].In[0].ID, lineNumber())
	assert.Equal("#main/reverse_sort", root.Graphs[0].Steps[1].In[1].Source[0], lineNumber())
	assert.Equal("#main/sorted/reverse", root.Graphs[0].Steps[1].In[1].ID, lineNumber())
	assert.Equal("#main/sorted/output", root.Graphs[0].Steps[1].Out[0].ID, lineNumber())
	// Graph 1
	assert.Equal("CommandLineTool", root.Graphs[1].Class, lineNumber())
	assert.Equal("#revtool.cwl", root.Graphs[1].ID, lineNumber())
	assert.Equal("Reverse each line using the `rev` command", root.Graphs[1].Doc, lineNumber())
	assert.Equal("rev", root.Graphs[1].BaseCommands[0], lineNumber())
	assert.Equal("output.txt", root.Graphs[1].Stdout, lineNumber())
	assert.Equal(1, len(root.Graphs[1].Inputs), lineNumber())
	assert.Equal("File", root.Graphs[1].Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("#revtool.cwl/input", root.Graphs[1].Inputs[0].ID, lineNumber())
	assert.Equal(1, len(root.Graphs[1].Outputs), lineNumber())
	assert.Equal("#revtool.cwl/output", root.Graphs[1].Outputs[0].ID, lineNumber())
	assert.Equal("File", root.Graphs[1].Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("output.txt", root.Graphs[1].Outputs[0].Binding.Glob[0], lineNumber())
	// Graph 2
	assert.Equal("CommandLineTool", root.Graphs[2].Class, lineNumber())
	assert.Equal("#sorttool.cwl", root.Graphs[2].ID, lineNumber())
	assert.Equal("Sort lines using the `sort` command", root.Graphs[2].Doc, lineNumber())
	assert.Equal("sort", root.Graphs[2].BaseCommands[0], lineNumber())
	assert.Equal("output.txt", root.Graphs[2].Stdout, lineNumber())
	assert.Equal(2, len(root.Graphs[2].Inputs), lineNumber())
	assert.Equal("boolean", root.Graphs[2].Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("#sorttool.cwl/reverse", root.Graphs[2].Inputs[0].ID, lineNumber())
	assert.Equal(1, root.Graphs[2].Inputs[0].Binding.Position, lineNumber())
	assert.Equal("--reverse", root.Graphs[2].Inputs[0].Binding.Prefix, lineNumber())
	assert.Equal("File", root.Graphs[2].Inputs[1].Types[0].Type, lineNumber())
	assert.Equal("#sorttool.cwl/input", root.Graphs[2].Inputs[1].ID, lineNumber())
	assert.Equal(2, root.Graphs[2].Inputs[1].Binding.Position, lineNumber())
	assert.Equal(1, len(root.Graphs[2].Outputs), lineNumber())
	assert.Equal("#sorttool.cwl/output", root.Graphs[2].Outputs[0].ID, lineNumber())
	assert.Equal("File", root.Graphs[2].Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("output.txt", root.Graphs[2].Outputs[0].Binding.Glob[0], lineNumber())
}

func revsortTest(assert *a.Assertions, root *Root) {
	assert.Equal("Workflow", root.Class, lineNumber())
	assert.Equal("Reverse the lines in a document, then sort those lines.", root.Doc, lineNumber())
	assert.Equal("DockerRequirement", root.Hints[0].Class, lineNumber())
	assert.Equal("debian:8", root.Hints[0].DockerPull, lineNumber())
	assert.Equal(2, len(root.Inputs), lineNumber())
	assert.Equal("input", root.Inputs[0].ID, lineNumber())
	assert.Equal("The input file to be processed.", root.Inputs[0].Doc, lineNumber())
	assert.Equal("File", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("boolean", root.Inputs[1].Types[0].Type, lineNumber())
	assert.Equal("reverse_sort", root.Inputs[1].ID, lineNumber())
	assert.Equal("If true, reverse (decending) sort", root.Inputs[1].Doc, lineNumber())
	assert.Equal("output", root.Outputs[0].ID, lineNumber())
	assert.Equal("File", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("sorted/output", root.Outputs[0].Source[0], lineNumber())
	assert.Equal("The output with the lines reversed and sorted.", root.Outputs[0].Doc[0], lineNumber())

	assert.Equal(2, len(root.Steps), lineNumber())
	assert.Equal("rev", root.Steps[0].ID, lineNumber())
	assert.Equal("input", root.Steps[0].In[0].ID, lineNumber())
	assert.Equal("input", root.Steps[0].In[0].Source[0], lineNumber())
	assert.Equal("output", root.Steps[0].Out[0].ID, lineNumber())
	assert.Equal("revtool.cwl", root.Steps[0].Run.Value, lineNumber())
	assert.Equal("sorted", root.Steps[1].ID, lineNumber())

	assert.Equal("input", root.Steps[1].In[0].ID, lineNumber())
	assert.Equal("rev/output", root.Steps[1].In[0].Source[0], lineNumber())
	assert.Equal("reverse", root.Steps[1].In[1].ID, lineNumber())
	assert.Equal("reverse_sort", root.Steps[1].In[1].Source[0], lineNumber())
	assert.Equal("output", root.Steps[1].Out[0].ID, lineNumber())
	assert.Equal("sorttool.cwl", root.Steps[1].Run.Value, lineNumber())
}

func revtoolTest(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal("Reverse each line using the `rev` command", root.Doc, lineNumber())
	assert.Equal("input", root.Inputs[0].ID, lineNumber())
	assert.Equal("File", root.Inputs[0].Types[0].Type, lineNumber())
	// TODO check specification for this test ID and Type
	assert.Equal("output", root.Outputs[0].ID, lineNumber())
	assert.Equal("File", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("output.txt", root.Outputs[0].Binding.Glob[0], lineNumber())
	assert.Equal("rev", root.BaseCommands[0], lineNumber())
}

func scatterValueFromToolTest(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())

	assert.Equal("scattered_message", root.Inputs[1].ID, lineNumber())
	assert.Equal("string", root.Inputs[1].Types[0].Type, lineNumber())
	assert.Equal(2, root.Inputs[1].Binding.Position, lineNumber())
	assert.Equal("message", root.Inputs[0].ID, lineNumber())
	assert.Equal("string", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal(1, root.Inputs[0].Binding.Position, lineNumber())
	assert.Equal("out_message", root.Outputs[0].ID, lineNumber())
	assert.Equal("stdout", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("echo", root.BaseCommands[0], lineNumber())
}

func scatterValuefromWf1Test(assert *a.Assertions, root *Root) {
	assert.Equal("Workflow", root.Class, lineNumber())

	assert.Equal("inp", root.Inputs[0].ID, lineNumber())
	assert.Equal("array", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("record", root.Inputs[0].Types[0].Items[0].Type, lineNumber())
	assert.Equal("instr", root.Inputs[0].Types[0].Items[0].Fields[0].Name, lineNumber())
	assert.Equal("string", root.Inputs[0].Types[0].Items[0].Fields[0].Types[0].Type, lineNumber())
	assert.Equal("out", root.Outputs[0].ID, lineNumber())
	assert.Equal("array", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("string", root.Outputs[0].Types[0].Items[0].Type, lineNumber())
	assert.Equal("ScatterFeatureRequirement", root.Requirements[0].Class, lineNumber())
	assert.Equal("StepInputExpressionRequirement", root.Requirements[1].Class, lineNumber())

	assert.Equal("step1", root.Steps[0].ID, lineNumber())
	assert.Equal("echo_in", root.Steps[0].In[0].ID, lineNumber())
	assert.Equal("inp", root.Steps[0].In[0].Source[0], lineNumber())
	assert.Equal("$(self.instr)", root.Steps[0].In[0].ValueFrom, lineNumber())
	assert.Equal("first", root.Steps[0].In[1].ID, lineNumber())
	assert.Equal("inp", root.Steps[0].In[1].Source[0], lineNumber())
	assert.Equal("$(self[0].instr)", root.Steps[0].In[1].ValueFrom, lineNumber())
	assert.Equal("echo_out", root.Steps[0].Out[0].ID, lineNumber())
	assert.Equal("echo_in", root.Steps[0].Scatter[0], lineNumber())
	assert.Equal("CommandLineTool", root.Steps[0].Run.Workflow.Class, lineNumber())

	assert.Equal("first", root.Steps[0].Run.Workflow.Inputs[0].ID, lineNumber())
	assert.Equal("string", root.Steps[0].Run.Workflow.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal(1, root.Steps[0].Run.Workflow.Inputs[0].Binding.Position, lineNumber())
	assert.Equal("echo_in", root.Steps[0].Run.Workflow.Inputs[1].ID, lineNumber())
	assert.Equal("string", root.Steps[0].Run.Workflow.Inputs[1].Types[0].Type, lineNumber())
	assert.Equal(2, root.Steps[0].Run.Workflow.Inputs[1].Binding.Position, lineNumber())
	assert.Equal("echo_out", root.Steps[0].Run.Workflow.Outputs[0].ID, lineNumber())
	assert.Equal("string", root.Steps[0].Run.Workflow.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("step1_out", root.Steps[0].Run.Workflow.Outputs[0].Binding.Glob[0], lineNumber())
	assert.Equal(true, root.Steps[0].Run.Workflow.Outputs[0].Binding.LoadContents, lineNumber())
	assert.Equal("$(self[0].contents)", root.Steps[0].Run.Workflow.Outputs[0].Binding.Eval, lineNumber())
	assert.Equal("echo", root.Steps[0].Run.Workflow.BaseCommands[0], lineNumber())
	assert.Equal("-n", root.Steps[0].Run.Workflow.Arguments[0].Value, lineNumber())
	assert.Equal("foo", root.Steps[0].Run.Workflow.Arguments[1].Value, lineNumber())
	assert.Equal("step1_out", root.Steps[0].Run.Workflow.Stdout, lineNumber())
}

func scatterValuefromWf2Test(assert *a.Assertions, root *Root) {
	assert.Equal("Workflow", root.Class, lineNumber())

	assert.Equal("inp1", root.Inputs[0].ID, lineNumber())
	assert.Equal("array", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("record", root.Inputs[0].Types[0].Items[0].Type, lineNumber())
	assert.Equal("instr", root.Inputs[0].Types[0].Items[0].Fields[0].Name, lineNumber())
	assert.Equal("string", root.Inputs[0].Types[0].Items[0].Fields[0].Types[0].Type, lineNumber())
	assert.Equal("inp2", root.Inputs[1].ID, lineNumber())
	assert.Equal("array", root.Inputs[1].Types[0].Type, lineNumber())
	assert.Equal("string", root.Inputs[1].Types[0].Items[0].Type, lineNumber())
	assert.Equal("out", root.Outputs[0].ID, lineNumber())
	assert.Equal("step1/echo_out", root.Outputs[0].Source[0], lineNumber())
	assert.Equal("array", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("array", root.Outputs[0].Types[0].Items[0].Type, lineNumber())
	assert.Equal("string", root.Outputs[0].Types[0].Items[0].Items[0].Type, lineNumber())
	assert.Equal("ScatterFeatureRequirement", root.Requirements[0].Class, lineNumber())
	assert.Equal("StepInputExpressionRequirement", root.Requirements[1].Class, lineNumber())
	assert.Equal("step1", root.Steps[0].ID, lineNumber())

	assert.Equal("echo_in1", root.Steps[0].In[0].ID, lineNumber())
	assert.Equal("inp1", root.Steps[0].In[0].Source[0], lineNumber())
	assert.Equal("$(self.instr)", root.Steps[0].In[0].ValueFrom, lineNumber())
	assert.Equal("echo_in2", root.Steps[0].In[1].ID, lineNumber())
	assert.Equal("inp2", root.Steps[0].In[1].Source[0], lineNumber())
	assert.Equal("first", root.Steps[0].In[2].ID, lineNumber())
	assert.Equal("inp1", root.Steps[0].In[2].Source[0], lineNumber())
	assert.Equal("$(self[0].instr)", root.Steps[0].In[2].ValueFrom, lineNumber())
	assert.Equal("echo_out", root.Steps[0].Out[0].ID, lineNumber())
	assert.Equal("echo_in1", root.Steps[0].Scatter[0], lineNumber())
	assert.Equal("echo_in2", root.Steps[0].Scatter[1], lineNumber())
	assert.Equal("nested_crossproduct", root.Steps[0].ScatterMethod, lineNumber())
	assert.Equal("CommandLineTool", root.Steps[0].Run.Workflow.Class, lineNumber())

	assert.Equal("first", root.Steps[0].Run.Workflow.Inputs[0].ID, lineNumber())
	assert.Equal("string", root.Steps[0].Run.Workflow.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal(1, root.Steps[0].Run.Workflow.Inputs[0].Binding.Position, lineNumber())
	assert.Equal("echo_in1", root.Steps[0].Run.Workflow.Inputs[1].ID, lineNumber())
	assert.Equal("string", root.Steps[0].Run.Workflow.Inputs[1].Types[0].Type, lineNumber())
	assert.Equal(2, root.Steps[0].Run.Workflow.Inputs[1].Binding.Position, lineNumber())
	assert.Equal("echo_in2", root.Steps[0].Run.Workflow.Inputs[2].ID, lineNumber())
	assert.Equal("string", root.Steps[0].Run.Workflow.Inputs[2].Types[0].Type, lineNumber())
	assert.Equal(3, root.Steps[0].Run.Workflow.Inputs[2].Binding.Position, lineNumber())
	assert.Equal("echo_out", root.Steps[0].Run.Workflow.Outputs[0].ID, lineNumber())
	assert.Equal("string", root.Steps[0].Run.Workflow.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("step1_out", root.Steps[0].Run.Workflow.Outputs[0].Binding.Glob[0], lineNumber())
	assert.Equal(true, root.Steps[0].Run.Workflow.Outputs[0].Binding.LoadContents, lineNumber())
	assert.Equal("$(self[0].contents)", root.Steps[0].Run.Workflow.Outputs[0].Binding.Eval, lineNumber())
	assert.Equal("echo", root.Steps[0].Run.Workflow.BaseCommands[0], lineNumber())
	assert.Equal("-n", root.Steps[0].Run.Workflow.Arguments[0].Value, lineNumber())
	assert.Equal("foo", root.Steps[0].Run.Workflow.Arguments[1].Value, lineNumber())
	assert.Equal("step1_out", root.Steps[0].Run.Workflow.Stdout, lineNumber())
}

func scatterValuefromWf3Test(assert *a.Assertions, root *Root) {
	assert.Equal("echo", root.Graphs[0].ID, lineNumber())
	assert.Equal("CommandLineTool", root.Graphs[0].Class, lineNumber())

	assert.Equal("first", root.Graphs[0].Inputs[0].ID, lineNumber())
	assert.Equal("string", root.Graphs[0].Inputs[0].Types[0].Type, lineNumber())
	assert.Equal(1, root.Graphs[0].Inputs[0].Binding.Position, lineNumber())
	assert.Equal("echo_in1", root.Graphs[0].Inputs[1].ID, lineNumber())
	assert.Equal("string", root.Graphs[0].Inputs[1].Types[0].Type, lineNumber())
	assert.Equal(2, root.Graphs[0].Inputs[1].Binding.Position, lineNumber())
	assert.Equal("echo_in2", root.Graphs[0].Inputs[2].ID, lineNumber())
	assert.Equal("string", root.Graphs[0].Inputs[2].Types[0].Type, lineNumber())
	assert.Equal(3, root.Graphs[0].Inputs[2].Binding.Position, lineNumber())
	assert.Equal("echo_out", root.Graphs[0].Outputs[0].ID, lineNumber())
	assert.Equal("string", root.Graphs[0].Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("step1_out", root.Graphs[0].Outputs[0].Binding.Glob[0], lineNumber())
	assert.Equal(true, root.Graphs[0].Outputs[0].Binding.LoadContents, lineNumber())
	assert.Equal("$(self[0].contents)", root.Graphs[0].Outputs[0].Binding.Eval, lineNumber())
	assert.Equal("echo", root.Graphs[0].BaseCommands[0], lineNumber())
	assert.Equal("-n", root.Graphs[0].Arguments[0].Value, lineNumber())
	assert.Equal("foo", root.Graphs[0].Arguments[1].Value, lineNumber())
	assert.Equal("step1_out", root.Graphs[0].Stdout, lineNumber())

	assert.Equal("main", root.Graphs[1].ID, lineNumber())
	assert.Equal("Workflow", root.Graphs[1].Class, lineNumber())

	assert.Equal("inp1", root.Graphs[1].Inputs[0].ID, lineNumber())
	assert.Equal("array", root.Graphs[1].Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("record", root.Graphs[1].Inputs[0].Types[0].Items[0].Type, lineNumber())
	assert.Equal("instr", root.Graphs[1].Inputs[0].Types[0].Items[0].Fields[0].Name, lineNumber())
	assert.Equal("string", root.Graphs[1].Inputs[0].Types[0].Items[0].Fields[0].Types[0].Type, lineNumber())
	assert.Equal("inp2", root.Graphs[1].Inputs[1].ID, lineNumber())
	assert.Equal("array", root.Graphs[1].Inputs[1].Types[0].Type, lineNumber())
	assert.Equal("string", root.Graphs[1].Inputs[1].Types[0].Items[0].Type, lineNumber())

	assert.Equal("ScatterFeatureRequirement", root.Graphs[1].Requirements[0].Class, lineNumber())
	assert.Equal("StepInputExpressionRequirement", root.Graphs[1].Requirements[1].Class, lineNumber())
	assert.Equal("echo_in1", root.Graphs[1].Steps[0].Scatter[0], lineNumber())
	assert.Equal("echo_in2", root.Graphs[1].Steps[0].Scatter[1], lineNumber())
	assert.Equal("flat_crossproduct", root.Graphs[1].Steps[0].ScatterMethod, lineNumber())
	assert.Equal("step1", root.Graphs[1].Steps[0].ID, lineNumber())

	assert.Equal("echo_in1", root.Graphs[1].Steps[0].In[0].ID, lineNumber())
	assert.Equal("inp1", root.Graphs[1].Steps[0].In[0].Source[0], lineNumber())
	assert.Equal("$(self.instr)", root.Graphs[1].Steps[0].In[0].ValueFrom, lineNumber())
	assert.Equal("echo_in2", root.Graphs[1].Steps[0].In[1].ID, lineNumber())
	assert.Equal("inp2", root.Graphs[1].Steps[0].In[1].Source[0], lineNumber())
	assert.Equal("first", root.Graphs[1].Steps[0].In[2].ID, lineNumber())
	assert.Equal("inp1", root.Graphs[1].Steps[0].In[2].Source[0], lineNumber())
	assert.Equal("$(self[0].instr)", root.Graphs[1].Steps[0].In[2].ValueFrom, lineNumber())
	assert.Equal("echo_out", root.Graphs[1].Steps[0].Out[0].ID, lineNumber())
	assert.Equal("#echo", root.Graphs[1].Steps[0].Run.Value, lineNumber())
	assert.Equal("out", root.Graphs[1].Outputs[0].ID, lineNumber())
	assert.Equal("step1/echo_out", root.Graphs[1].Outputs[0].Source[0], lineNumber())
	assert.Equal("array", root.Graphs[1].Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("string", root.Graphs[1].Outputs[0].Types[0].Items[0].Type, lineNumber())
}

func scatterValuefromWf4Test(assert *a.Assertions, root *Root) {
	assert.Equal("echo", root.Graphs[0].ID, lineNumber())
	assert.Equal("CommandLineTool", root.Graphs[0].Class, lineNumber())

	assert.Equal("first", root.Graphs[0].Inputs[0].ID, lineNumber())
	assert.Equal("string", root.Graphs[0].Inputs[0].Types[0].Type, lineNumber())
	assert.Equal(1, root.Graphs[0].Inputs[0].Binding.Position, lineNumber())
	assert.Equal("echo_in1", root.Graphs[0].Inputs[1].ID, lineNumber())
	assert.Equal("string", root.Graphs[0].Inputs[1].Types[0].Type, lineNumber())
	assert.Equal(2, root.Graphs[0].Inputs[1].Binding.Position, lineNumber())
	assert.Equal("echo_in2", root.Graphs[0].Inputs[2].ID, lineNumber())
	assert.Equal("string", root.Graphs[0].Inputs[2].Types[0].Type, lineNumber())
	assert.Equal(3, root.Graphs[0].Inputs[2].Binding.Position, lineNumber())
	assert.Equal("echo_out", root.Graphs[0].Outputs[0].ID, lineNumber())
	assert.Equal("string", root.Graphs[0].Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("step1_out", root.Graphs[0].Outputs[0].Binding.Glob[0], lineNumber())
	assert.Equal(true, root.Graphs[0].Outputs[0].Binding.LoadContents, lineNumber())
	assert.Equal("$(self[0].contents)", root.Graphs[0].Outputs[0].Binding.Eval, lineNumber())
	assert.Equal("echo", root.Graphs[0].BaseCommands[0], lineNumber())
	assert.Equal("-n", root.Graphs[0].Arguments[0].Value, lineNumber())
	assert.Equal("foo", root.Graphs[0].Arguments[1].Value, lineNumber())
	assert.Equal("step1_out", root.Graphs[0].Stdout, lineNumber())

	assert.Equal("main", root.Graphs[1].ID, lineNumber())
	assert.Equal("Workflow", root.Graphs[1].Class, lineNumber())
	assert.Equal("inp1", root.Graphs[1].Inputs[0].ID, lineNumber())
	assert.Equal("array", root.Graphs[1].Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("record", root.Graphs[1].Inputs[0].Types[0].Items[0].Type, lineNumber())
	assert.Equal("instr", root.Graphs[1].Inputs[0].Types[0].Items[0].Fields[0].Name, lineNumber())
	assert.Equal("string", root.Graphs[1].Inputs[0].Types[0].Items[0].Fields[0].Types[0].Type, lineNumber())
	assert.Equal("inp2", root.Graphs[1].Inputs[1].ID, lineNumber())
	assert.Equal("array", root.Graphs[1].Inputs[1].Types[0].Type, lineNumber())
	assert.Equal("string", root.Graphs[1].Inputs[1].Types[0].Items[0].Type, lineNumber())
	assert.Equal("ScatterFeatureRequirement", root.Graphs[1].Requirements[0].Class, lineNumber())
	assert.Equal("StepInputExpressionRequirement", root.Graphs[1].Requirements[1].Class, lineNumber())
	assert.Equal("echo_in1", root.Graphs[1].Steps[0].Scatter[0], lineNumber())
	assert.Equal("echo_in2", root.Graphs[1].Steps[0].Scatter[1], lineNumber())
	assert.Equal("dotproduct", root.Graphs[1].Steps[0].ScatterMethod, lineNumber())
	assert.Equal("step1", root.Graphs[1].Steps[0].ID, lineNumber())

	assert.Equal("echo_in1", root.Graphs[1].Steps[0].In[0].ID, lineNumber())
	assert.Equal("inp1", root.Graphs[1].Steps[0].In[0].Source[0], lineNumber())
	assert.Equal("$(self.instr)", root.Graphs[1].Steps[0].In[0].ValueFrom, lineNumber())
	assert.Equal("echo_in2", root.Graphs[1].Steps[0].In[1].ID, lineNumber())
	assert.Equal("inp2", root.Graphs[1].Steps[0].In[1].Source[0], lineNumber())
	assert.Equal("first", root.Graphs[1].Steps[0].In[2].ID, lineNumber())
	assert.Equal("inp1", root.Graphs[1].Steps[0].In[2].Source[0], lineNumber())
	assert.Equal("$(self[0].instr)", root.Graphs[1].Steps[0].In[2].ValueFrom, lineNumber())
	assert.Equal("echo_out", root.Graphs[1].Steps[0].Out[0].ID, lineNumber())
	assert.Equal("#echo", root.Graphs[1].Steps[0].Run.Value, lineNumber())
	assert.Equal("out", root.Graphs[1].Outputs[0].ID, lineNumber())
	assert.Equal("step1/echo_out", root.Graphs[1].Outputs[0].Source[0], lineNumber())
	assert.Equal("array", root.Graphs[1].Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("string", root.Graphs[1].Outputs[0].Types[0].Items[0].Type, lineNumber())
}

func scatterValuefromWf5Test(assert *a.Assertions, root *Root) {
	assert.Equal("Workflow", root.Class, lineNumber())

	assert.Equal("inp", root.Inputs[0].ID, lineNumber())
	assert.Equal("array", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("record", root.Inputs[0].Types[0].Items[0].Type, lineNumber())
	assert.Equal("instr", root.Inputs[0].Types[0].Items[0].Fields[0].Name, lineNumber())
	assert.Equal("string", root.Inputs[0].Types[0].Items[0].Fields[0].Types[0].Type, lineNumber())
	assert.Equal("out", root.Outputs[0].ID, lineNumber())
	assert.Equal("array", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("string", root.Outputs[0].Types[0].Items[0].Type, lineNumber())
	assert.Equal("step1/echo_out", root.Outputs[0].Source[0], lineNumber())
	assert.Equal("ScatterFeatureRequirement", root.Requirements[0].Class, lineNumber())
	assert.Equal("StepInputExpressionRequirement", root.Requirements[1].Class, lineNumber())
	assert.Equal("step1", root.Steps[0].ID, lineNumber())

	assert.Equal("echo_in", root.Steps[0].In[0].ID, lineNumber())
	assert.Equal("inp", root.Steps[0].In[0].Source[0], lineNumber())
	assert.Equal("$(self.instr)", root.Steps[0].In[0].ValueFrom, lineNumber())
	assert.Equal("first", root.Steps[0].In[1].ID, lineNumber())
	assert.Equal("inp", root.Steps[0].In[1].Source[0], lineNumber())
	assert.Equal("$(inputs.echo_in.instr)", root.Steps[0].In[1].ValueFrom, lineNumber())
	assert.Equal("echo_out", root.Steps[0].Out[0].ID, lineNumber())
	assert.Equal("echo_in", root.Steps[0].Scatter[0], lineNumber())
	assert.Equal("CommandLineTool", root.Steps[0].Run.Workflow.Class, lineNumber())

	assert.Equal("first", root.Steps[0].Run.Workflow.Inputs[0].ID, lineNumber())
	assert.Equal("string", root.Steps[0].Run.Workflow.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal(1, root.Steps[0].Run.Workflow.Inputs[0].Binding.Position, lineNumber())
	assert.Equal("echo_in", root.Steps[0].Run.Workflow.Inputs[1].ID, lineNumber())
	assert.Equal("string", root.Steps[0].Run.Workflow.Inputs[1].Types[0].Type, lineNumber())
	assert.Equal(2, root.Steps[0].Run.Workflow.Inputs[1].Binding.Position, lineNumber())
	assert.Equal("echo_out", root.Steps[0].Run.Workflow.Outputs[0].ID, lineNumber())
	assert.Equal("string", root.Steps[0].Run.Workflow.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("step1_out", root.Steps[0].Run.Workflow.Outputs[0].Binding.Glob[0], lineNumber())
	assert.Equal(true, root.Steps[0].Run.Workflow.Outputs[0].Binding.LoadContents, lineNumber())
	assert.Equal("$(self[0].contents)", root.Steps[0].Run.Workflow.Outputs[0].Binding.Eval, lineNumber())
	assert.Equal("echo", root.Steps[0].Run.Workflow.BaseCommands[0], lineNumber())
	assert.Equal("-n", root.Steps[0].Run.Workflow.Arguments[0].Value, lineNumber())
	assert.Equal("foo", root.Steps[0].Run.Workflow.Arguments[1].Value, lineNumber())
	assert.Equal("step1_out", root.Steps[0].Run.Workflow.Stdout, lineNumber())
}

func scatterValuefromWf6Test(assert *a.Assertions, root *Root) {
	assert.Equal("Workflow", root.Class, lineNumber())
	assert.Equal("ScatterFeatureRequirement", root.Requirements[0].Class, lineNumber())
	assert.Equal("StepInputExpressionRequirement", root.Requirements[1].Class, lineNumber())
	assert.Equal("scattered_messages", root.Inputs[0].ID, lineNumber())
	assert.Equal("string[]", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("out_message", root.Outputs[0].ID, lineNumber())
	assert.Equal("File[]", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("step1/out_message", root.Outputs[0].Source[0], lineNumber())
	assert.Equal("step1", root.Steps[0].ID, lineNumber())
	assert.Equal("scatter-valueFrom-tool.cwl", root.Steps[0].Run.Value, lineNumber())
	assert.Equal("scattered_message", root.Steps[0].Scatter[0], lineNumber())
	assert.Equal("dotproduct", root.Steps[0].ScatterMethod, lineNumber())

	assert.Equal("scattered_message", root.Steps[0].In[1].ID, lineNumber())
	assert.Equal("scattered_messages", root.Steps[0].In[1].Source[0], lineNumber())
	assert.Equal("message", root.Steps[0].In[0].ID, lineNumber())
	assert.Equal("Hello", root.Steps[0].In[0].ValueFrom, lineNumber())
	assert.Equal("out_message", root.Steps[0].Out[0].ID, lineNumber())
}

func scatterWf1Test(assert *a.Assertions, root *Root) {
	assert.Equal("Workflow", root.Class, lineNumber())

	assert.Equal("inp", root.Inputs[0].ID, lineNumber())
	assert.Equal("string[]", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("out", root.Outputs[0].ID, lineNumber())
	assert.Equal("string[]", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("step1/echo_out", root.Outputs[0].Source[0], lineNumber())
	assert.Equal("ScatterFeatureRequirement", root.Requirements[0].Class, lineNumber())
	assert.Equal("step1", root.Steps[0].ID, lineNumber())
	assert.Equal("echo_in", root.Steps[0].In[0].ID, lineNumber())
	assert.Equal("inp", root.Steps[0].In[0].Source[0], lineNumber())
	assert.Equal("echo_out", root.Steps[0].Out[0].ID, lineNumber())
	assert.Equal("echo_in", root.Steps[0].Scatter[0], lineNumber())
	assert.Equal("CommandLineTool", root.Steps[0].Run.Workflow.Class, lineNumber())
	assert.Equal("echo_in", root.Steps[0].Run.Workflow.Inputs[0].ID, lineNumber())
	assert.Equal("string", root.Steps[0].Run.Workflow.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("echo_out", root.Steps[0].Run.Workflow.Outputs[0].ID, lineNumber())
	assert.Equal("string", root.Steps[0].Run.Workflow.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("step1_out", root.Steps[0].Run.Workflow.Outputs[0].Binding.Glob[0], lineNumber())
	assert.Equal(true, root.Steps[0].Run.Workflow.Outputs[0].Binding.LoadContents, lineNumber())
	assert.Equal("$(self[0].contents)", root.Steps[0].Run.Workflow.Outputs[0].Binding.Eval, lineNumber())
	assert.Equal("echo", root.Steps[0].Run.Workflow.BaseCommands[0], lineNumber())
	assert.Equal("-n", root.Steps[0].Run.Workflow.Arguments[0].Value, lineNumber())
	assert.Equal("foo", root.Steps[0].Run.Workflow.Arguments[1].Value, lineNumber())
	assert.Equal("step1_out", root.Steps[0].Run.Workflow.Stdout, lineNumber())
}

func scatterWf2Test(assert *a.Assertions, root *Root) {
	assert.Equal("Workflow", root.Class, lineNumber())

	assert.Equal("inp1", root.Inputs[0].ID, lineNumber())
	assert.Equal("string[]", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("inp2", root.Inputs[1].ID, lineNumber())
	assert.Equal("string[]", root.Inputs[1].Types[0].Type, lineNumber())
	assert.Equal("out", root.Outputs[0].ID, lineNumber())
	assert.Equal("step1/echo_out", root.Outputs[0].Source[0], lineNumber())
	assert.Equal("array", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("array", root.Outputs[0].Types[0].Items[0].Type, lineNumber())
	assert.Equal("string", root.Outputs[0].Types[0].Items[0].Items[0].Type, lineNumber())
	assert.Equal("ScatterFeatureRequirement", root.Requirements[0].Class, lineNumber())

	assert.Equal("step1", root.Steps[0].ID, lineNumber())
	assert.Equal("echo_in1", root.Steps[0].In[0].ID, lineNumber())
	assert.Equal("inp1", root.Steps[0].In[0].Source[0], lineNumber())
	assert.Equal("echo_in2", root.Steps[0].In[1].ID, lineNumber())
	assert.Equal("inp2", root.Steps[0].In[1].Source[0], lineNumber())
	assert.Equal("echo_out", root.Steps[0].Out[0].ID, lineNumber())
	assert.Equal("echo_in1", root.Steps[0].Scatter[0], lineNumber())
	assert.Equal("echo_in2", root.Steps[0].Scatter[1], lineNumber())
	assert.Equal("nested_crossproduct", root.Steps[0].ScatterMethod, lineNumber())
	assert.Equal("CommandLineTool", root.Steps[0].Run.Workflow.Class, lineNumber())
	assert.Equal("step1command", root.Steps[0].Run.Workflow.ID, lineNumber())
	assert.Equal("echo_in1", root.Steps[0].Run.Workflow.Inputs[0].ID, lineNumber())
	assert.Equal("string", root.Steps[0].Run.Workflow.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("echo_in2", root.Steps[0].Run.Workflow.Inputs[1].ID, lineNumber())
	assert.Equal("string", root.Steps[0].Run.Workflow.Inputs[1].Types[0].Type, lineNumber())
	assert.Equal("echo_out", root.Steps[0].Run.Workflow.Outputs[0].ID, lineNumber())
	assert.Equal("string", root.Steps[0].Run.Workflow.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("step1_out", root.Steps[0].Run.Workflow.Outputs[0].Binding.Glob[0], lineNumber())
	assert.Equal(true, root.Steps[0].Run.Workflow.Outputs[0].Binding.LoadContents, lineNumber())
	assert.Equal("$(self[0].contents)", root.Steps[0].Run.Workflow.Outputs[0].Binding.Eval, lineNumber())
	assert.Equal("echo", root.Steps[0].Run.Workflow.BaseCommands[0], lineNumber())
	assert.Equal("-n", root.Steps[0].Run.Workflow.Arguments[0].Value, lineNumber())
	assert.Equal("foo", root.Steps[0].Run.Workflow.Arguments[1].Value, lineNumber())
	assert.Equal("step1_out", root.Steps[0].Run.Workflow.Stdout, lineNumber())
}

func scatterWf3Test(assert *a.Assertions, root *Root) {
	assert.Equal("echo", root.Graphs[0].ID, lineNumber())
	assert.Equal("CommandLineTool", root.Graphs[0].Class, lineNumber())
	assert.Equal("echo_in1", root.Graphs[0].Inputs[0].ID, lineNumber())
	assert.Equal("string", root.Graphs[0].Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("echo_in2", root.Graphs[0].Inputs[1].ID, lineNumber())
	assert.Equal("string", root.Graphs[0].Inputs[1].Types[0].Type, lineNumber())
	assert.Equal("echo_out", root.Graphs[0].Outputs[0].ID, lineNumber())
	assert.Equal("string", root.Graphs[0].Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("step1_out", root.Graphs[0].Outputs[0].Binding.Glob[0], lineNumber())
	assert.Equal(true, root.Graphs[0].Outputs[0].Binding.LoadContents, lineNumber())
	assert.Equal("$(self[0].contents)", root.Graphs[0].Outputs[0].Binding.Eval, lineNumber())
	assert.Equal("echo", root.Graphs[0].BaseCommands[0], lineNumber())
	assert.Equal("-n", root.Graphs[0].Arguments[0].Value, lineNumber())
	assert.Equal("foo", root.Graphs[0].Arguments[1].Value, lineNumber())
	assert.Equal("step1_out", root.Graphs[0].Stdout, lineNumber())

	assert.Equal("main", root.Graphs[1].ID, lineNumber())
	assert.Equal("Workflow", root.Graphs[1].Class, lineNumber())

	assert.Equal("inp1", root.Graphs[1].Inputs[0].ID, lineNumber())
	assert.Equal("string[]", root.Graphs[1].Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("inp2", root.Graphs[1].Inputs[1].ID, lineNumber())
	assert.Equal("string[]", root.Graphs[1].Inputs[1].Types[0].Type, lineNumber())
	assert.Equal("ScatterFeatureRequirement", root.Graphs[1].Requirements[0].Class, lineNumber())
	assert.Equal("echo_in1", root.Graphs[1].Steps[0].Scatter[0], lineNumber())
	assert.Equal("echo_in2", root.Graphs[1].Steps[0].Scatter[1], lineNumber())
	assert.Equal("flat_crossproduct", root.Graphs[1].Steps[0].ScatterMethod, lineNumber())
	assert.Equal("step1", root.Graphs[1].Steps[0].ID, lineNumber())

	assert.Equal("echo_in1", root.Graphs[1].Steps[0].In[0].ID, lineNumber())
	assert.Equal("inp1", root.Graphs[1].Steps[0].In[0].Source[0], lineNumber())
	assert.Equal("echo_in2", root.Graphs[1].Steps[0].In[1].ID, lineNumber())
	assert.Equal("inp2", root.Graphs[1].Steps[0].In[1].Source[0], lineNumber())
	assert.Equal("echo_out", root.Graphs[1].Steps[0].Out[0].ID, lineNumber())
	assert.Equal("#echo", root.Graphs[1].Steps[0].Run.Value, lineNumber())

	assert.Equal("out", root.Graphs[1].Outputs[0].ID, lineNumber())
	assert.Equal("step1/echo_out", root.Graphs[1].Outputs[0].Source[0], lineNumber())
	assert.Equal("array", root.Graphs[1].Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("string", root.Graphs[1].Outputs[0].Types[0].Items[0].Type, lineNumber())
}

func scatterWf4Test(assert *a.Assertions, root *Root) {
	assert.Equal("echo", root.Graphs[0].ID, lineNumber())
	assert.Equal("CommandLineTool", root.Graphs[0].Class, lineNumber())
	assert.Equal("echo_in1", root.Graphs[0].Inputs[0].ID, lineNumber())
	assert.Equal("string", root.Graphs[0].Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("echo_in2", root.Graphs[0].Inputs[1].ID, lineNumber())
	assert.Equal("string", root.Graphs[0].Inputs[1].Types[0].Type, lineNumber())
	assert.Equal("echo_out", root.Graphs[0].Outputs[0].ID, lineNumber())
	assert.Equal("string", root.Graphs[0].Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("step1_out", root.Graphs[0].Outputs[0].Binding.Glob[0], lineNumber())
	assert.Equal(true, root.Graphs[0].Outputs[0].Binding.LoadContents, lineNumber())
	assert.Equal("$(self[0].contents)", root.Graphs[0].Outputs[0].Binding.Eval, lineNumber())
	assert.Equal("echo", root.Graphs[0].BaseCommands[0], lineNumber())
	assert.Equal("-n", root.Graphs[0].Arguments[0].Value, lineNumber())
	assert.Equal("foo", root.Graphs[0].Arguments[1].Value, lineNumber())
	assert.Equal("step1_out", root.Graphs[0].Stdout, lineNumber())

	assert.Equal("main", root.Graphs[1].ID, lineNumber())
	assert.Equal("Workflow", root.Graphs[1].Class, lineNumber())
	assert.Equal("inp1", root.Graphs[1].Inputs[0].ID, lineNumber())
	assert.Equal("string[]", root.Graphs[1].Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("inp2", root.Graphs[1].Inputs[1].ID, lineNumber())
	assert.Equal("string[]", root.Graphs[1].Inputs[1].Types[0].Type, lineNumber())
	assert.Equal("ScatterFeatureRequirement", root.Graphs[1].Requirements[0].Class, lineNumber())
	assert.Equal("echo_in1", root.Graphs[1].Steps[0].Scatter[0], lineNumber())
	assert.Equal("echo_in2", root.Graphs[1].Steps[0].Scatter[1], lineNumber())
	assert.Equal("dotproduct", root.Graphs[1].Steps[0].ScatterMethod, lineNumber())
	assert.Equal("step1", root.Graphs[1].Steps[0].ID, lineNumber())
	assert.Equal("echo_in1", root.Graphs[1].Steps[0].In[0].ID, lineNumber())
	assert.Equal("inp1", root.Graphs[1].Steps[0].In[0].Source[0], lineNumber())
	assert.Equal("echo_in2", root.Graphs[1].Steps[0].In[1].ID, lineNumber())
	assert.Equal("inp2", root.Graphs[1].Steps[0].In[1].Source[0], lineNumber())
	assert.Equal("echo_out", root.Graphs[1].Steps[0].Out[0].ID, lineNumber())
	assert.Equal("#echo", root.Graphs[1].Steps[0].Run.Value, lineNumber())

	assert.Equal("out", root.Graphs[1].Outputs[0].ID, lineNumber())
	assert.Equal("step1/echo_out", root.Graphs[1].Outputs[0].Source[0], lineNumber())
	assert.Equal("array", root.Graphs[1].Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("string", root.Graphs[1].Outputs[0].Types[0].Items[0].Type, lineNumber())

}

func schemadefToolTest(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal("schemadef-type.yml", root.Requirements[0].Import, lineNumber())
	assert.Equal("InlineJavascriptRequirement", root.Requirements[1].Class, lineNumber())

	assert.Equal("hello", root.Inputs[0].ID, lineNumber())
	assert.Equal("schemadef-type.yml#HelloType", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal(`self.a + "/" + self.b`, root.Inputs[0].Binding.ValueFrom.Key(), lineNumber())
	assert.Equal("output", root.Outputs[0].ID, lineNumber())
	assert.Equal("File", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("output.txt", root.Outputs[0].Binding.Glob[0], lineNumber())
	assert.Equal("output.txt", root.Stdout, lineNumber())
	assert.Equal("echo", root.BaseCommands[0], lineNumber())
}

func schemadefWfTest(assert *a.Assertions, root *Root) {
	assert.Equal("Workflow", root.Class, lineNumber())
	assert.Equal("schemadef-type.yml", root.Requirements[0].Import, lineNumber())
	assert.Equal("hello", root.Inputs[0].ID, lineNumber())
	assert.Equal("schemadef-type.yml#HelloType", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("output", root.Outputs[0].ID, lineNumber())
	assert.Equal("File", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("step1/output", root.Outputs[0].Source[0], lineNumber())
	assert.Equal("step1", root.Steps[0].ID, lineNumber())
	assert.Equal("hello", root.Steps[0].In[0].ID, lineNumber())
	assert.Equal("hello", root.Steps[0].In[0].Source[0], lineNumber())
	assert.Equal("output", root.Steps[0].Out[0].ID, lineNumber())
	assert.Equal("schemadef-tool.cwl", root.Steps[0].Run.Value, lineNumber())
}

func searchTest(assert *a.Assertions, root *Root) {
	assert.Equal("index", root.Graphs[0].ID, lineNumber())
	assert.Equal("CommandLineTool", root.Graphs[0].Class, lineNumber())
	assert.Equal("python", root.Graphs[0].BaseCommands[0], lineNumber())
	assert.Equal("input.txt", root.Graphs[0].Arguments[0].Binding.ValueFrom.Key(), lineNumber())
	assert.Equal(1, root.Graphs[0].Arguments[0].Binding.Position, lineNumber())
	assert.Equal("InitialWorkDirRequirement", root.Graphs[0].Requirements[0].Class, lineNumber())
	assert.Equal("input.txt", root.Graphs[0].Requirements[0].Listing[0].EntryName, lineNumber())
	assert.Equal("$(inputs.file)", root.Graphs[0].Requirements[0].Listing[0].Entry, lineNumber())
	assert.Equal("InlineJavascriptRequirement", root.Graphs[0].Requirements[1].Class, lineNumber())
	assert.Equal("DockerRequirement", root.Graphs[0].Hints[0].Class, lineNumber())
	assert.Equal("python:2-slim", root.Graphs[0].Hints[0].DockerPull, lineNumber())
	// skip Default value check
	count := 0
	for _, input := range root.Graphs[0].Inputs {
		id := input.ID
		switch id {
		case "file":
			assert.Equal("File", input.Types[0].Type, lineNumber())
			count = count + 1
		case "secondfile":
			assert.Equal("File", input.Types[0].Type, lineNumber())
			count = count + 1
		case "index.py":
			assert.Equal("File", input.Types[0].Type, lineNumber())
			assert.Equal(0, input.Binding.Position, lineNumber())
			count = count + 1
		}
	}
	assert.Equal(3, count, lineNumber())
	assert.Equal("result", root.Graphs[0].Outputs[0].ID, lineNumber())
	assert.Equal("File", root.Graphs[0].Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("input.txt", root.Graphs[0].Outputs[0].Binding.Glob[0], lineNumber())
	assert.Equal(".idx1", root.Graphs[0].Outputs[0].SecondaryFiles[0].Entry, lineNumber())
	assert.Equal("^.idx2", root.Graphs[0].Outputs[0].SecondaryFiles[1].Entry, lineNumber())
	assert.Equal(`$(self.basename).idx3`, root.Graphs[0].Outputs[0].SecondaryFiles[2].Entry, lineNumber())
	assert.Equal(`${ return self.basename+".idx4"; }`, root.Graphs[0].Outputs[0].SecondaryFiles[3].Entry, lineNumber())
	assert.Equal(`$({"path": self.path+".idx5", "class": "File"})`, root.Graphs[0].Outputs[0].SecondaryFiles[4].Entry, lineNumber())
	assert.Equal(`$(self.nameroot).idx6$(self.nameext)`, root.Graphs[0].Outputs[0].SecondaryFiles[5].Entry, lineNumber())
	assert.Equal(`${ return [self.basename+".idx7", inputs.secondfile]; }`, root.Graphs[0].Outputs[0].SecondaryFiles[6].Entry, lineNumber())
	assert.Equal("_idx8", root.Graphs[0].Outputs[0].SecondaryFiles[7].Entry, lineNumber())

	assert.Equal("search", root.Graphs[1].ID, lineNumber())
	assert.Equal("CommandLineTool", root.Graphs[1].Class, lineNumber())
	assert.Equal("python", root.Graphs[1].BaseCommands[0], lineNumber())
	assert.Equal("InlineJavascriptRequirement", root.Graphs[1].Requirements[0].Class, lineNumber())
	assert.Equal("DockerRequirement", root.Graphs[1].Hints[0].Class, lineNumber())
	assert.Equal("python:2-slim", root.Graphs[1].Hints[0].DockerPull, lineNumber())

	assert.Equal("search.py", root.Graphs[1].Inputs[0].ID, lineNumber())
	assert.Equal("File", root.Graphs[1].Inputs[0].Types[0].Type, lineNumber())
	assert.Equal(0, root.Graphs[1].Inputs[0].Binding.Position, lineNumber())
	assert.Equal("term", root.Graphs[1].Inputs[2].ID, lineNumber())
	assert.Equal("string", root.Graphs[1].Inputs[2].Types[0].Type, lineNumber())
	assert.Equal(2, root.Graphs[1].Inputs[2].Binding.Position, lineNumber())
	assert.Equal("result", root.Graphs[1].Outputs[0].ID, lineNumber())
	assert.Equal("File", root.Graphs[1].Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("result.txt", root.Graphs[1].Outputs[0].Binding.Glob[0], lineNumber())

	assert.Equal("file", root.Graphs[1].Inputs[1].ID, lineNumber())
	assert.Equal("File", root.Graphs[1].Inputs[1].Types[0].Type, lineNumber())
	assert.Equal(1, root.Graphs[1].Inputs[1].Binding.Position, lineNumber())
	assert.Equal(".idx1", root.Graphs[1].Inputs[1].SecondaryFiles[0].Entry, lineNumber())
	assert.Equal("^.idx2", root.Graphs[1].Inputs[1].SecondaryFiles[1].Entry, lineNumber())
	assert.Equal(`$(self.basename).idx3`, root.Graphs[1].Inputs[1].SecondaryFiles[2].Entry, lineNumber())
	assert.Equal(`${ return self.basename+".idx4"; }`, root.Graphs[1].Inputs[1].SecondaryFiles[3].Entry, lineNumber())
	assert.Equal(`$(self.nameroot).idx6$(self.nameext)`, root.Graphs[1].Inputs[1].SecondaryFiles[4].Entry, lineNumber())
	assert.Equal(`${ return [self.basename+".idx7"]; }`, root.Graphs[1].Inputs[1].SecondaryFiles[5].Entry, lineNumber())
	assert.Equal("_idx8", root.Graphs[1].Inputs[1].SecondaryFiles[6].Entry, lineNumber())

	assert.Equal("main", root.Graphs[2].ID, lineNumber())
	assert.Equal("Workflow", root.Graphs[2].Class, lineNumber())
	count2 := 0
	for _, input := range root.Graphs[2].Inputs {
		id := input.ID
		switch id {
		case "infile":
			assert.Equal("File", input.Types[0].Type, lineNumber())
			count2 = count2 + 1
		case "secondfile":
			assert.Equal("File", input.Types[0].Type, lineNumber())
			count2 = count2 + 1
		case "term":
			assert.Equal("string", input.Types[0].Type, lineNumber())
			count2 = count2 + 1
		}
	}
	assert.Equal(3, count2, lineNumber())

	assert.Equal("indexedfile", root.Graphs[2].Outputs[0].ID, lineNumber())
	assert.Equal("File", root.Graphs[2].Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("index/result", root.Graphs[2].Outputs[0].Source[0], lineNumber())
	assert.Equal("outfile", root.Graphs[2].Outputs[1].ID, lineNumber())
	assert.Equal("File", root.Graphs[2].Outputs[1].Types[0].Type, lineNumber())
	assert.Equal("search/result", root.Graphs[2].Outputs[1].Source[0], lineNumber())

	assert.Equal("index", root.Graphs[2].Steps[0].ID, lineNumber())
	assert.Equal("#index", root.Graphs[2].Steps[0].Run.Value, lineNumber())

	assert.Equal("file", root.Graphs[2].Steps[0].In[0].ID, lineNumber())
	assert.Equal("infile", root.Graphs[2].Steps[0].In[0].Source[0], lineNumber())
	assert.Equal("secondfile", root.Graphs[2].Steps[0].In[1].ID, lineNumber())
	assert.Equal("secondfile", root.Graphs[2].Steps[0].In[1].Source[0], lineNumber())
	assert.Equal("result", root.Graphs[2].Steps[0].Out[0].ID, lineNumber())
	assert.Equal("search", root.Graphs[2].Steps[1].ID, lineNumber())
	assert.Equal("#search", root.Graphs[2].Steps[1].Run.Value, lineNumber())
	assert.Equal("file", root.Graphs[2].Steps[1].In[0].ID, lineNumber())
	assert.Equal("index/result", root.Graphs[2].Steps[1].In[0].Source[0], lineNumber())
	assert.Equal("term", root.Graphs[2].Steps[1].In[1].ID, lineNumber())
	assert.Equal("term", root.Graphs[2].Steps[1].In[1].Source[0], lineNumber())
	assert.Equal("result", root.Graphs[2].Steps[1].Out[0].ID, lineNumber())
}

func shellchar2Test(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal("Ensure that `shellQuote: true` is the default behavior when\n"+"ShellCommandRequirement is in effect.\n", root.Doc, lineNumber())
	assert.Equal("ShellCommandRequirement", root.Requirements[0].Class, lineNumber())
	assert.Equal(0, len(root.Inputs), lineNumber())
	assert.Equal(2, len(root.Outputs), lineNumber())
	count := 0
	for _, out := range root.Outputs {
		switch out.ID {
		case "stdout_file":
			assert.Equal("stdout", out.Types[0].Type, lineNumber())
			count = count + 1
		case "stderr_file":
			assert.Equal("stderr", out.Types[0].Type, lineNumber())
			count = count + 1
		}
	}
	assert.Equal(2, count, lineNumber())
	assert.Equal("echo", root.BaseCommands[0], lineNumber())
	assert.Equal("foo 1>&2", root.Arguments[0].Value, lineNumber())
}

func shellcharTest(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal("Ensure that arguments containing shell directives are not interpreted and\n"+"that `shellQuote: false` has no effect when ShellCommandRequirement is not in\n"+"effect.\n", root.Doc, lineNumber())
	assert.Equal(0, len(root.Inputs), lineNumber())
	assert.Equal(2, len(root.Outputs), lineNumber())
	count := 0
	for _, out := range root.Outputs {
		switch out.ID {
		case "stdout_file":
			assert.Equal("stdout", out.Types[0].Type, lineNumber())
			count = count + 1
		case "stderr_file":
			assert.Equal("stderr", out.Types[0].Type, lineNumber())
			count = count + 1
		}
	}
	assert.Equal(2, count, lineNumber())
	assert.Equal("echo", root.BaseCommands[0], lineNumber())
	assert.Equal("foo 1>&2", root.Arguments[0].Binding.ValueFrom.Key(), lineNumber())
	assert.Equal(false, root.Arguments[0].Binding.ShellQuote, lineNumber())
}

func shelltestTest(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal("Reverse each line using the `rev` command then sort.", root.Doc, lineNumber())
	assert.Equal("ShellCommandRequirement", root.Requirements[0].Class, lineNumber())
	assert.Equal(1, len(root.Inputs), lineNumber())
	assert.Equal("input", root.Inputs[0].ID, lineNumber())
	assert.Equal("File", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal(1, len(root.Outputs), lineNumber())
	assert.Equal("output", root.Outputs[0].ID, lineNumber())
	assert.Equal("File", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("output.txt", root.Outputs[0].Binding.Glob[0], lineNumber())
	assert.Equal("rev", root.Arguments[0].Value, lineNumber())
	assert.Equal("inputs.input", root.Arguments[1].Binding.ValueFrom.Key(), lineNumber())
	assert.Equal(" | ", root.Arguments[2].Binding.ValueFrom.Key(), lineNumber())
	assert.Equal(false, root.Arguments[2].Binding.ShellQuote, lineNumber())
	assert.Equal("sort", root.Arguments[3].Value, lineNumber())
	assert.Equal("> output.txt", root.Arguments[4].Binding.ValueFrom.Key(), lineNumber())
	assert.Equal(false, root.Arguments[4].Binding.ShellQuote, lineNumber())
}

func sorttoolTest(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal("Sort lines using the `sort` command", root.Doc, lineNumber())
	assert.Equal(2, len(root.Inputs), lineNumber())
	assert.Equal("reverse", root.Inputs[0].ID, lineNumber())
	assert.Equal("boolean", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal(1, root.Inputs[0].Binding.Position, lineNumber())
	assert.Equal("--reverse", root.Inputs[0].Binding.Prefix, lineNumber())
	assert.Equal("input", root.Inputs[1].ID, lineNumber())
	assert.Equal("File", root.Inputs[1].Types[0].Type, lineNumber())
	assert.Equal(2, root.Inputs[1].Binding.Position, lineNumber())
	assert.Equal(1, len(root.Outputs), lineNumber())
	assert.Equal("output", root.Outputs[0].ID, lineNumber())
	assert.Equal("File", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("output.txt", root.Outputs[0].Binding.Glob[0], lineNumber())
	assert.Equal("sort", root.BaseCommands[0], lineNumber())
	assert.Equal("output.txt", root.Stdout, lineNumber())
}

func stagefileTest(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal("DockerRequirement", root.Hints[0].Class, lineNumber())
	assert.Equal("python:2-slim", root.Hints[0].DockerPull, lineNumber())
	assert.Equal("InitialWorkDirRequirement", root.Requirements[0].Class, lineNumber())
	assert.Equal("$(inputs.infile)", root.Requirements[0].Listing[0].Entry, lineNumber())
	assert.Equal("bob.txt", root.Requirements[0].Listing[0].EntryName, lineNumber())
	assert.Equal(true, root.Requirements[0].Listing[0].Writable, lineNumber())
	assert.Equal(1, len(root.Inputs), lineNumber())
	assert.Equal("infile", root.Inputs[0].ID, lineNumber())
	assert.Equal("File", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal(1, len(root.Outputs), lineNumber())
	assert.Equal("outfile", root.Outputs[0].ID, lineNumber())
	assert.Equal("File", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("bob.txt", root.Outputs[0].Binding.Glob[0], lineNumber())
	assert.Equal("python2", root.BaseCommands[0], lineNumber())
	assert.Equal("-c", root.Arguments[0].Value, lineNumber())
	assert.Equal(`f = open("bob.txt", "r+")
f.seek(8)
f.write("Bob.    ")
f.close()
`, root.Arguments[1].Value, lineNumber())
}

func stderrMediumcutTest(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal("Test of capturing stderr output in a docker container.", root.Doc, lineNumber())
	assert.Equal("ShellCommandRequirement", root.Requirements[0].Class, lineNumber())
	assert.Equal(0, len(root.Inputs), lineNumber())
	assert.Equal(1, len(root.Outputs), lineNumber())
	assert.Equal("output_file", root.Outputs[0].ID, lineNumber())
	assert.Equal("stderr", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("echo foo 1>&2", root.Arguments[0].Binding.ValueFrom.Key(), lineNumber())
	assert.Equal(false, root.Arguments[0].Binding.ShellQuote, lineNumber())
	assert.Equal("std.err", root.Stderr, lineNumber())
}

func stderrShortcutTest(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal("Test of capturing stderr output in a docker container.", root.Doc, lineNumber())
	assert.Equal("ShellCommandRequirement", root.Requirements[0].Class, lineNumber())
	assert.Equal(0, len(root.Inputs), lineNumber())
	assert.Equal(1, len(root.Outputs), lineNumber())
	assert.Equal("output_file", root.Outputs[0].ID, lineNumber())
	assert.Equal("stderr", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("echo foo 1>&2", root.Arguments[0].Binding.ValueFrom.Key(), lineNumber())
	assert.Equal(false, root.Arguments[0].Binding.ShellQuote, lineNumber())
}

func stderrTest(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal("Test of capturing stderr output in a docker container.", root.Doc, lineNumber())
	assert.Equal("ShellCommandRequirement", root.Requirements[0].Class, lineNumber())
	assert.Equal(0, len(root.Inputs), lineNumber())
	assert.Equal(1, len(root.Outputs), lineNumber())
	assert.Equal("output_file", root.Outputs[0].ID, lineNumber())
	assert.Equal("File", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("error.txt", root.Outputs[0].Binding.Glob[0], lineNumber())
	assert.Equal("echo foo 1>&2", root.Arguments[0].Binding.ValueFrom.Key(), lineNumber())
	assert.Equal(false, root.Arguments[0].Binding.ShellQuote, lineNumber())
	assert.Equal("error.txt", root.Stderr, lineNumber())
}

func stepValuefrom2WfTest(assert *a.Assertions, root *Root) {
	assert.Equal("Workflow", root.Class, lineNumber())
	assert.Equal("StepInputExpressionRequirement", root.Requirements[0].Class, lineNumber())
	assert.Equal("InlineJavascriptRequirement", root.Requirements[1].Class, lineNumber())
	assert.Equal("MultipleInputFeatureRequirement", root.Requirements[2].Class, lineNumber())

	assert.Equal("a", root.Inputs[0].ID, lineNumber())
	assert.Equal("int", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("b", root.Inputs[1].ID, lineNumber())
	assert.Equal("int", root.Inputs[1].Types[0].Type, lineNumber())
	assert.Equal("val", root.Outputs[0].ID, lineNumber())
	assert.Equal("string", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("step1/echo_out", root.Outputs[0].Source[0], lineNumber())
	assert.Equal("step1", root.Steps[0].ID, lineNumber())
	// TODO test run: id: echo
	assert.Equal("CommandLineTool", root.Steps[0].Run.Workflow.Class, lineNumber())
	assert.Equal("c", root.Steps[0].Run.Workflow.Inputs[0].ID, lineNumber())
	assert.Equal("int", root.Steps[0].Run.Workflow.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("echo_out", root.Steps[0].Run.Workflow.Outputs[0].ID, lineNumber())
	assert.Equal("string", root.Steps[0].Run.Workflow.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("step1_out", root.Steps[0].Run.Workflow.Outputs[0].Binding.Glob[0], lineNumber())
	assert.Equal(true, root.Steps[0].Run.Workflow.Outputs[0].Binding.LoadContents, lineNumber())
	assert.Equal("$(self[0].contents)", root.Steps[0].Run.Workflow.Outputs[0].Binding.Eval, lineNumber())
	assert.Equal("echo", root.Steps[0].Run.Workflow.BaseCommands[0], lineNumber())
	assert.Equal("step1_out", root.Steps[0].Run.Workflow.Stdout, lineNumber())
	assert.Equal("c", root.Steps[0].In[0].ID, lineNumber())
	assert.Equal("a", root.Steps[0].In[0].Source[0], lineNumber())
	assert.Equal("b", root.Steps[0].In[0].Source[1], lineNumber())
	assert.Equal("$(self[0] + self[1])", root.Steps[0].In[0].ValueFrom, lineNumber())
	assert.Equal("echo_out", root.Steps[0].Out[0].ID, lineNumber())
}

func stepValuefrom3WfTest(assert *a.Assertions, root *Root) {
	assert.Equal("Workflow", root.Class, lineNumber())
	assert.Equal("StepInputExpressionRequirement", root.Requirements[0].Class, lineNumber())
	assert.Equal("InlineJavascriptRequirement", root.Requirements[1].Class, lineNumber())

	assert.Equal("a", root.Inputs[0].ID, lineNumber())
	assert.Equal("int", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("b", root.Inputs[1].ID, lineNumber())
	assert.Equal("int", root.Inputs[1].Types[0].Type, lineNumber())
	assert.Equal("val", root.Outputs[0].ID, lineNumber())
	assert.Equal("string", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("step1/echo_out", root.Outputs[0].Source[0], lineNumber())
	assert.Equal("step1", root.Steps[0].ID, lineNumber())
	// TODO test run: id: echo
	assert.Equal("CommandLineTool", root.Steps[0].Run.Workflow.Class, lineNumber())
	assert.Equal("c", root.Steps[0].Run.Workflow.Inputs[0].ID, lineNumber())
	assert.Equal("int", root.Steps[0].Run.Workflow.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("echo_out", root.Steps[0].Run.Workflow.Outputs[0].ID, lineNumber())
	assert.Equal("string", root.Steps[0].Run.Workflow.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("step1_out", root.Steps[0].Run.Workflow.Outputs[0].Binding.Glob[0], lineNumber())
	assert.Equal(true, root.Steps[0].Run.Workflow.Outputs[0].Binding.LoadContents, lineNumber())
	assert.Equal("$(self[0].contents)", root.Steps[0].Run.Workflow.Outputs[0].Binding.Eval, lineNumber())
	assert.Equal("echo", root.Steps[0].Run.Workflow.BaseCommands[0], lineNumber())
	assert.Equal("step1_out", root.Steps[0].Run.Workflow.Stdout, lineNumber())

	assert.Equal("a", root.Steps[0].In[0].ID, lineNumber())
	assert.Equal("a", root.Steps[0].In[0].Source[0], lineNumber())
	assert.Equal("b", root.Steps[0].In[1].ID, lineNumber())
	assert.Equal("b", root.Steps[0].In[1].Source[0], lineNumber())
	assert.Equal("c", root.Steps[0].In[2].ID, lineNumber())
	assert.Equal("$(inputs.a + inputs.b)", root.Steps[0].In[2].ValueFrom, lineNumber())
	assert.Equal("echo_out", root.Steps[0].Out[0].ID, lineNumber())
}

func stepValuefromWfTest(assert *a.Assertions, root *Root) {
	assert.Equal("Workflow", root.Class, lineNumber())
	assert.Equal("in", root.Inputs[0].ID, lineNumber())
	assert.Equal("record", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("file1", root.Inputs[0].Types[0].Fields[0].Name, lineNumber())
	assert.Equal("File", root.Inputs[0].Types[0].Fields[0].Types[0].Type, lineNumber())
	// TODO check specification for this test ID and Type
	assert.Equal("count_output", root.Outputs[0].ID, lineNumber())
	assert.Equal("int", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("step2/output", root.Outputs[0].Source[0], lineNumber())
	assert.Equal("wc-tool.cwl", root.Steps[0].Run.Value, lineNumber())
	assert.Equal("file1", root.Steps[0].In[0].ID, lineNumber())
	assert.Equal("in", root.Steps[0].In[0].Source[0], lineNumber())
	assert.Equal("$(self.file1)", root.Steps[0].In[0].ValueFrom, lineNumber())
	assert.Equal("output", root.Steps[0].Out[0].ID, lineNumber())
	assert.Equal("parseInt-tool.cwl", root.Steps[1].Run.Value, lineNumber())
	assert.Equal("file1", root.Steps[1].In[0].ID, lineNumber())
	assert.Equal("step1/output", root.Steps[1].In[0].Source[0], lineNumber())
	assert.Equal("output", root.Steps[1].Out[0].ID, lineNumber())
}

func sumWfTest(assert *a.Assertions, root *Root) {
	assert.Equal("Workflow", root.Class, lineNumber())
	assert.Equal("StepInputExpressionRequirement", root.Requirements[0].Class, lineNumber())
	assert.Equal("MultipleInputFeatureRequirement", root.Requirements[1].Class, lineNumber())
	assert.Equal("InlineJavascriptRequirement", root.Requirements[2].Class, lineNumber())
	assert.Equal(2, len(root.Inputs), lineNumber())
	assert.Equal("int_1", root.Inputs[0].ID, lineNumber())
	assert.Equal("int", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("string", root.Inputs[0].Types[1].Type, lineNumber())
	assert.Equal("int_2", root.Inputs[1].ID, lineNumber())
	assert.Equal("int", root.Inputs[1].Types[0].Type, lineNumber())
	assert.Equal("string", root.Inputs[1].Types[1].Type, lineNumber())
	assert.Equal("result", root.Outputs[0].ID, lineNumber())
	assert.Equal("int", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("sum/result", root.Outputs[0].Source[0], lineNumber())
	assert.Equal("sum", root.Steps[0].ID, lineNumber())
	assert.Equal("data", root.Steps[0].In[0].ID, lineNumber())
	assert.Equal("int_1", root.Steps[0].In[0].Source[0], lineNumber())
	assert.Equal("int_2", root.Steps[0].In[0].Source[1], lineNumber())
	assert.Equal(`${
  var sum = 0;
  for (var i = 0; i < self.length; i++){
    sum += self[i];
  };
  return sum;
}
`, root.Steps[0].In[0].ValueFrom)
	assert.Equal("result", root.Steps[0].Out[0].ID, lineNumber())
	assert.Equal("ExpressionTool", root.Steps[0].Run.Workflow.Class, lineNumber())
	assert.Equal("data", root.Steps[0].Run.Workflow.Inputs[0].ID, lineNumber())
	assert.Equal("int", root.Steps[0].Run.Workflow.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal("result", root.Steps[0].Run.Workflow.Outputs[0].ID, lineNumber())
	assert.Equal("int", root.Steps[0].Run.Workflow.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal(`${
  return {"result": inputs.data};
}
`, root.Steps[0].Run.Workflow.Expression)
}

func templateToolTest(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())

	assert.Equal("InlineJavascriptRequirement", root.Requirements[0].Class, lineNumber())
	assert.Equal("underscore.js", root.Requirements[0].ExpressionLib[0].Value, lineNumber())
	assert.Equal("var t = function(s) { return _.template(s)({'inputs': inputs}); };", root.Requirements[0].ExpressionLib[1].Value, lineNumber())

	assert.Equal("InitialWorkDirRequirement", root.Requirements[1].Class, lineNumber())
	assert.Equal("foo.txt", root.Requirements[1].Listing[0].EntryName, lineNumber())
	assert.Equal(`$(t("The file is <%= inputs.file1.path.split('/').slice(-1)[0] %>\n"))`, root.Requirements[1].Listing[0].Entry, lineNumber())

	assert.Equal("DockerRequirement", root.Hints[0].Class, lineNumber())
	assert.Equal("debian:8", root.Hints[0].DockerPull, lineNumber())

	assert.Equal("file1", root.Inputs[0].ID, lineNumber())
	assert.Equal("File", root.Inputs[0].Types[0].Type, lineNumber())

	assert.Equal("foo", root.Outputs[0].ID, lineNumber())
	assert.Equal("File", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal([]string{"foo.txt"}, root.Outputs[0].Binding.Glob, lineNumber())

	assert.Equal("cat", root.BaseCommands[0], lineNumber())
	assert.Equal("foo.txt", root.BaseCommands[1], lineNumber())
}

func testCwlOut2Test(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal("ShellCommandRequirement", root.Requirements[0].Class, lineNumber())
	assert.Equal("DockerRequirement", root.Hints[0].Class, lineNumber())
	assert.Equal("debian:wheezy", root.Hints[0].DockerPull, lineNumber())
	assert.Equal(0, len(root.Inputs), lineNumber())
	// TODO check specification for this test ID and Type
	assert.Equal("foo", root.Outputs[0].ID, lineNumber())
	assert.Equal("File", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal(`echo foo > foo && echo '{"foo": {"location": "file://$(runtime.outdir)/foo", "class": "File"} }' > cwl.output.json
`, root.Arguments[0].Binding.ValueFrom.Key(), lineNumber())
	assert.Equal(false, root.Arguments[0].Binding.ShellQuote, lineNumber())
}

func testCwlOutTest(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal("ShellCommandRequirement", root.Requirements[0].Class, lineNumber())
	assert.Equal("DockerRequirement", root.Hints[0].Class, lineNumber())
	assert.Equal("debian:wheezy", root.Hints[0].DockerPull, lineNumber())
	assert.Equal(0, len(root.Inputs), lineNumber())
	// TODO check specification for this test ID and Type
	assert.Equal("foo", root.Outputs[0].ID, lineNumber())
	assert.Equal("File", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal(`echo foo > foo && echo '{"foo": {"path": "$(runtime.outdir)/foo", "class": "File"} }' > cwl.output.json
`, root.Arguments[0].Binding.ValueFrom.Key(), lineNumber())
}

func tmapToolTest(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())

	assert.Equal("DockerRequirement", root.Hints[0].Class, lineNumber())
	assert.Equal("python:2-slim", root.Hints[0].DockerPull, lineNumber())

	assert.Equal("#args.py", root.Inputs[0].ID, lineNumber())
	assert.Equal("File", root.Inputs[0].Types[0].Type, lineNumber())
	assert.Equal(reflect.Map, root.Inputs[0].Default.Kind, lineNumber())
	assert.Equal(-1, root.Inputs[0].Binding.Position, lineNumber())
	assert.Equal("reads", root.Inputs[1].ID, lineNumber())
	assert.Equal("File", root.Inputs[1].Types[0].Type, lineNumber())
	assert.Equal("stages", root.Inputs[2].ID, lineNumber())
	assert.Equal("array", root.Inputs[2].Types[0].Type, lineNumber())
	assert.Equal("#Stage", root.Inputs[2].Types[0].Items[0].Type, lineNumber())
	assert.Equal(1, root.Inputs[2].Binding.Position, lineNumber())

	assert.Equal("sam", root.Outputs[0].ID, lineNumber())
	assert.Equal([]string{"output.sam"}, root.Outputs[0].Binding.Glob, lineNumber())
	assert.Equal("null", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("File", root.Outputs[0].Types[1].Type, lineNumber())
	assert.Equal("args", root.Outputs[1].ID, lineNumber())
	assert.Equal("string[]", root.Outputs[1].Types[0].Type, lineNumber())

	assert.Equal("SchemaDefRequirement", root.Requirements[0].Class, lineNumber())
	// assert.Equal("Map1", root.Requirements[0].Types[0].Name, lineNumber())
	assert.Equal("record", root.Requirements[0].Types[0].Type, lineNumber())
	assert.Equal("algo", root.Requirements[0].Types[0].Fields[0].Name, lineNumber())
	assert.Equal("enum", root.Requirements[0].Types[0].Fields[0].Types[0].Type, lineNumber())
	// assert.Equal("JustMap1", root.Requirements[0].Types[0].Fields[0].Types[0].Name, lineNumber())
	assert.Equal("map1", root.Requirements[0].Types[0].Fields[0].Types[0].Symbols[0], lineNumber())
	assert.Equal(0, root.Requirements[0].Types[0].Fields[0].Binding.Position, lineNumber())
	assert.Equal("maxSeqLen", root.Requirements[0].Types[0].Fields[1].Name, lineNumber())
	assert.Equal("null", root.Requirements[0].Types[0].Fields[1].Types[0].Type, lineNumber())
	assert.Equal("int", root.Requirements[0].Types[0].Fields[1].Types[1].Type, lineNumber())
	assert.Equal(2, root.Requirements[0].Types[0].Fields[1].Binding.Position, lineNumber())
	assert.Equal("--max-seq-length", root.Requirements[0].Types[0].Fields[1].Binding.Prefix, lineNumber())
	assert.Equal("minSeqLen", root.Requirements[0].Types[0].Fields[2].Name, lineNumber())
	assert.Equal("null", root.Requirements[0].Types[0].Fields[2].Types[0].Type, lineNumber())
	assert.Equal("int", root.Requirements[0].Types[0].Fields[2].Types[1].Type, lineNumber())
	assert.Equal(2, root.Requirements[0].Types[0].Fields[2].Binding.Position, lineNumber())
	assert.Equal("--min-seq-length", root.Requirements[0].Types[0].Fields[2].Binding.Prefix, lineNumber())
	assert.Equal("seedLength", root.Requirements[0].Types[0].Fields[3].Name, lineNumber())
	assert.Equal("null", root.Requirements[0].Types[0].Fields[3].Types[0].Type, lineNumber())
	assert.Equal("int", root.Requirements[0].Types[0].Fields[3].Types[1].Type, lineNumber())
	assert.Equal(2, root.Requirements[0].Types[0].Fields[3].Binding.Position, lineNumber())
	assert.Equal("--seed-length", root.Requirements[0].Types[0].Fields[3].Binding.Prefix, lineNumber())

	// assert.Equal("Map2", root.Requirements[0].Types[1].Name, lineNumber())
	assert.Equal("record", root.Requirements[0].Types[1].Type, lineNumber())
	assert.Equal("algo", root.Requirements[0].Types[1].Fields[0].Name, lineNumber())
	assert.Equal("enum", root.Requirements[0].Types[1].Fields[0].Types[0].Type, lineNumber())
	// assert.Equal("JustMap2", root.Requirements[0].Types[1].Fields[0].Types[0].Name, lineNumber())
	assert.Equal("map2", root.Requirements[0].Types[1].Fields[0].Types[0].Symbols[0], lineNumber())
	assert.Equal(0, root.Requirements[0].Types[1].Fields[0].Binding.Position, lineNumber())
	assert.Equal("maxSeqLen", root.Requirements[0].Types[1].Fields[1].Name, lineNumber())
	assert.Equal("null", root.Requirements[0].Types[1].Fields[1].Types[0].Type, lineNumber())
	assert.Equal("int", root.Requirements[0].Types[1].Fields[1].Types[1].Type, lineNumber())
	assert.Equal(2, root.Requirements[0].Types[1].Fields[1].Binding.Position, lineNumber())
	assert.Equal("--max-seq-length", root.Requirements[0].Types[1].Fields[1].Binding.Prefix, lineNumber())
	assert.Equal("minSeqLen", root.Requirements[0].Types[1].Fields[2].Name, lineNumber())
	assert.Equal("null", root.Requirements[0].Types[1].Fields[2].Types[0].Type, lineNumber())
	assert.Equal("int", root.Requirements[0].Types[1].Fields[2].Types[1].Type, lineNumber())
	assert.Equal(2, root.Requirements[0].Types[1].Fields[2].Binding.Position, lineNumber())
	assert.Equal("--min-seq-length", root.Requirements[0].Types[1].Fields[2].Binding.Prefix, lineNumber())
	assert.Equal("maxSeedHits", root.Requirements[0].Types[1].Fields[3].Name, lineNumber())
	assert.Equal("null", root.Requirements[0].Types[1].Fields[3].Types[0].Type, lineNumber())
	assert.Equal("int", root.Requirements[0].Types[1].Fields[3].Types[1].Type, lineNumber())
	assert.Equal(2, root.Requirements[0].Types[1].Fields[3].Binding.Position, lineNumber())
	assert.Equal("--max-seed-hits", root.Requirements[0].Types[1].Fields[3].Binding.Prefix, lineNumber())

	// assert.Equal("Map3", root.Requirements[0].Types[2].Name, lineNumber())
	assert.Equal("record", root.Requirements[0].Types[2].Type, lineNumber())
	assert.Equal("algo", root.Requirements[0].Types[2].Fields[0].Name, lineNumber())
	assert.Equal("enum", root.Requirements[0].Types[2].Fields[0].Types[0].Type, lineNumber())
	// assert.Equal("JustMap3", root.Requirements[0].Types[2].Fields[0].Types[0].Name, lineNumber())
	assert.Equal("map3", root.Requirements[0].Types[2].Fields[0].Types[0].Symbols[0], lineNumber())
	assert.Equal(0, root.Requirements[0].Types[2].Fields[0].Binding.Position, lineNumber())
	assert.Equal("maxSeqLen", root.Requirements[0].Types[2].Fields[1].Name, lineNumber())
	assert.Equal("null", root.Requirements[0].Types[2].Fields[1].Types[0].Type, lineNumber())
	assert.Equal("int", root.Requirements[0].Types[2].Fields[1].Types[1].Type, lineNumber())
	assert.Equal(2, root.Requirements[0].Types[2].Fields[1].Binding.Position, lineNumber())
	assert.Equal("--max-seq-length", root.Requirements[0].Types[2].Fields[1].Binding.Prefix, lineNumber())
	assert.Equal("minSeqLen", root.Requirements[0].Types[2].Fields[2].Name, lineNumber())
	assert.Equal("null", root.Requirements[0].Types[2].Fields[2].Types[0].Type, lineNumber())
	assert.Equal("int", root.Requirements[0].Types[2].Fields[2].Types[1].Type, lineNumber())
	assert.Equal(2, root.Requirements[0].Types[2].Fields[2].Binding.Position, lineNumber())
	assert.Equal("--min-seq-length", root.Requirements[0].Types[2].Fields[2].Binding.Prefix, lineNumber())
	assert.Equal("fwdSearch", root.Requirements[0].Types[2].Fields[3].Name, lineNumber())
	assert.Equal("null", root.Requirements[0].Types[2].Fields[3].Types[0].Type, lineNumber())
	assert.Equal("boolean", root.Requirements[0].Types[2].Fields[3].Types[1].Type, lineNumber())
	assert.Equal(2, root.Requirements[0].Types[2].Fields[3].Binding.Position, lineNumber())
	assert.Equal("--fwd-search", root.Requirements[0].Types[2].Fields[3].Binding.Prefix, lineNumber())

	// assert.Equal("Map4", root.Requirements[0].Types[3].Name, lineNumber())
	assert.Equal("record", root.Requirements[0].Types[3].Type, lineNumber())
	assert.Equal("algo", root.Requirements[0].Types[3].Fields[0].Name, lineNumber())
	assert.Equal("enum", root.Requirements[0].Types[3].Fields[0].Types[0].Type, lineNumber())
	// assert.Equal("JustMap4", root.Requirements[0].Types[3].Fields[0].Types[0].Name, lineNumber())
	assert.Equal("map4", root.Requirements[0].Types[3].Fields[0].Types[0].Symbols[0], lineNumber())
	assert.Equal(0, root.Requirements[0].Types[3].Fields[0].Binding.Position, lineNumber())
	assert.Equal("maxSeqLen", root.Requirements[0].Types[3].Fields[1].Name, lineNumber())
	assert.Equal("null", root.Requirements[0].Types[3].Fields[1].Types[0].Type, lineNumber())
	assert.Equal("int", root.Requirements[0].Types[3].Fields[1].Types[1].Type, lineNumber())
	assert.Equal(2, root.Requirements[0].Types[3].Fields[1].Binding.Position, lineNumber())
	assert.Equal("--max-seq-length", root.Requirements[0].Types[3].Fields[1].Binding.Prefix, lineNumber())
	assert.Equal("minSeqLen", root.Requirements[0].Types[3].Fields[2].Name, lineNumber())
	assert.Equal("null", root.Requirements[0].Types[3].Fields[2].Types[0].Type, lineNumber())
	assert.Equal("int", root.Requirements[0].Types[3].Fields[2].Types[1].Type, lineNumber())
	assert.Equal(2, root.Requirements[0].Types[3].Fields[2].Binding.Position, lineNumber())
	assert.Equal("--min-seq-length", root.Requirements[0].Types[3].Fields[2].Binding.Prefix, lineNumber())
	assert.Equal("seedStep", root.Requirements[0].Types[3].Fields[3].Name, lineNumber())
	assert.Equal("null", root.Requirements[0].Types[3].Fields[3].Types[0].Type, lineNumber())
	assert.Equal("int", root.Requirements[0].Types[3].Fields[3].Types[1].Type, lineNumber())
	assert.Equal(2, root.Requirements[0].Types[3].Fields[3].Binding.Position, lineNumber())
	assert.Equal("--seed-step", root.Requirements[0].Types[3].Fields[3].Binding.Prefix, lineNumber())

	// assert.Equal("Stage", root.Requirements[0].Types[4].Name, lineNumber())
	assert.Equal("record", root.Requirements[0].Types[4].Type, lineNumber())
	assert.Equal("stageId", root.Requirements[0].Types[4].Fields[0].Name, lineNumber())
	assert.Equal("null", root.Requirements[0].Types[4].Fields[0].Types[0].Type, lineNumber())
	assert.Equal(0, root.Requirements[0].Types[4].Fields[0].Binding.Position, lineNumber())
	assert.Equal("stage", root.Requirements[0].Types[4].Fields[0].Binding.Prefix, lineNumber())
	assert.Equal(false, root.Requirements[0].Types[4].Fields[0].Binding.Separate, lineNumber())
	assert.Equal("stageOption1", root.Requirements[0].Types[4].Fields[1].Name, lineNumber())
	assert.Equal("null", root.Requirements[0].Types[4].Fields[1].Types[0].Type, lineNumber())
	assert.Equal("boolean", root.Requirements[0].Types[4].Fields[1].Types[1].Type, lineNumber())
	assert.Equal(1, root.Requirements[0].Types[4].Fields[1].Binding.Position, lineNumber())
	assert.Equal("-n", root.Requirements[0].Types[4].Fields[1].Binding.Prefix, lineNumber())
	assert.Equal("algos", root.Requirements[0].Types[4].Fields[2].Name, lineNumber())
	assert.Equal("array", root.Requirements[0].Types[4].Fields[2].Types[0].Type, lineNumber())
	assert.Equal("#Map1", root.Requirements[0].Types[4].Fields[2].Types[0].Items[0].Type, lineNumber())
	assert.Equal("#Map2", root.Requirements[0].Types[4].Fields[2].Types[0].Items[1].Type, lineNumber())
	assert.Equal("#Map3", root.Requirements[0].Types[4].Fields[2].Types[0].Items[2].Type, lineNumber())
	assert.Equal("#Map4", root.Requirements[0].Types[4].Fields[2].Types[0].Items[3].Type, lineNumber())
	assert.Equal(2, root.Requirements[0].Types[4].Fields[2].Binding.Position, lineNumber())
}

func wc2ToolTest(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal("InlineJavascriptRequirement", root.Requirements[0].Class, lineNumber())
	assert.Equal("file1", root.Inputs[0].ID, lineNumber())
	assert.Equal("File", root.Inputs[0].Types[0].Type, lineNumber())
	// TODO check specification for this test ID and Type
	assert.Equal("output", root.Outputs[0].ID, lineNumber())
	assert.Equal("int", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("output.txt", root.Outputs[0].Binding.Glob[0], lineNumber())
	assert.Equal(true, root.Outputs[0].Binding.LoadContents, lineNumber())
	assert.Equal("$(parseInt(self[0].contents))", root.Outputs[0].Binding.Eval, lineNumber())
	assert.Equal(1, len(root.BaseCommands), lineNumber())
	assert.Equal("wc", root.BaseCommands[0], lineNumber())
	assert.Equal("output.txt", root.Stdout, lineNumber())
}

func wc3ToolTest(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal("InlineJavascriptRequirement", root.Requirements[0].Class, lineNumber())
	assert.Equal("file1", root.Inputs[0].ID, lineNumber())
	assert.Equal("File[]", root.Inputs[0].Types[0].Type, lineNumber())
	// TODO check specification for this test ID and Type
	assert.Equal("output", root.Outputs[0].ID, lineNumber())
	assert.Equal("int", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("output.txt", root.Outputs[0].Binding.Glob[0], lineNumber())
	assert.Equal(true, root.Outputs[0].Binding.LoadContents, lineNumber())
	assert.Equal(`${
  var s = self[0].contents.split(/\r?\n/);
  return parseInt(s[s.length-2]);
}
`, root.Outputs[0].Binding.Eval)
	assert.Equal(1, len(root.BaseCommands), lineNumber())
	assert.Equal("wc", root.BaseCommands[0], lineNumber())
	assert.Equal("output.txt", root.Stdout, lineNumber())
}

func wc4ToolTest(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal("InlineJavascriptRequirement", root.Requirements[0].Class, lineNumber())
	assert.Equal("file1", root.Inputs[0].ID, lineNumber())
	assert.Equal("File", root.Inputs[0].Types[0].Type, lineNumber())
	// TODO check specification for this test ID and Type
	assert.Equal("output", root.Outputs[0].ID, lineNumber())
	assert.Equal("int", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("output.txt", root.Outputs[0].Binding.Glob[0], lineNumber())
	assert.Equal(true, root.Outputs[0].Binding.LoadContents, lineNumber())
	assert.Equal(`${
  var s = self[0].contents.split(/\r?\n/);
  return parseInt(s[s.length-2]);
}
`, root.Outputs[0].Binding.Eval)
	assert.Equal(1, len(root.BaseCommands), lineNumber())
	assert.Equal("wc", root.BaseCommands[0], lineNumber())
	assert.Equal("output.txt", root.Stdout, lineNumber())
}

func wcToolTest(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal("file1", root.Inputs[0].ID, lineNumber())
	assert.Equal("File", root.Inputs[0].Types[0].Type, lineNumber())
	// TODO check specification for this test ID and Type
	assert.Equal("output", root.Outputs[0].ID, lineNumber())
	assert.Equal("File", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("output", root.Outputs[0].Binding.Glob[0], lineNumber())
	assert.Equal(2, len(root.BaseCommands), lineNumber())
	assert.Equal("wc", root.BaseCommands[0], lineNumber())
	assert.Equal("-l", root.BaseCommands[1], lineNumber())
	assert.Equal("$(inputs.file1.path)", root.Stdin, lineNumber())
	assert.Equal("output", root.Stdout, lineNumber())
}

func writableDirTest(assert *a.Assertions, root *Root) {
	assert.Equal("CommandLineTool", root.Class, lineNumber())
	assert.Equal("InlineJavascriptRequirement", root.Requirements[1].Class, lineNumber())
	assert.Equal("InitialWorkDirRequirement", root.Requirements[0].Class, lineNumber())
	assert.Equal("emptyWritableDir", root.Requirements[0].Listing[0].EntryName, lineNumber())
	assert.Equal(true, root.Requirements[0].Listing[0].Writable, lineNumber())
	assert.Equal("$({class: 'Directory', listing: []})", root.Requirements[0].Listing[0].Entry, lineNumber())
	assert.Equal(0, len(root.Inputs), lineNumber())
	// TODO check specification for this test ID and Type
	assert.Equal("out", root.Outputs[0].ID, lineNumber())
	assert.Equal("Directory", root.Outputs[0].Types[0].Type, lineNumber())
	assert.Equal("emptyWritableDir", root.Outputs[0].Binding.Glob[0], lineNumber())
	assert.Equal(2, len(root.Arguments), lineNumber())
	assert.Equal("touch", root.Arguments[0].Value, lineNumber())
	assert.Equal("emptyWritableDir/blurg", root.Arguments[1].Value, lineNumber())
}

// test regex:
// Expect\(t, (.+)\)\.ToBe\((.+)\)
// assert.Equal($2, $1, lineNumber())
