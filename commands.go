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
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Command is a high-level interpretation of a concrete command line that needs
// to be run as part of a workflow.
type Command struct {
	// ID is a generated string that identifies the command as coming from a
	// (particular step in a) particular workflow, resolved against particular
	// parameters.
	ID string

	// Cmd is the executable to run along with any command line arguments.
	Cmd []string

	// Cwd is the directory you should execute the Cmd in. $HOME should be set
	// to this while executing.
	Cwd string

	// TmpPrefix is the parent of a unique directory that should be created
	// before Cmd is executed. That unique dir should be set as $TMPDIR, and be
	// deleted afterwards.
	TmpPrefix string

	// Env are the environment variables the Cmd must be executed with. The
	// value is in the same format as that of os.Environ().
	Env []string

	// StdInPath, if non-blank, is the path to a file that should be piped in to
	// Cmd when executed.
	StdInPath string

	// StdOutPath, if non-blank, is the path to a file that the STDOUT of the
	// executed Cmd should be redirected to.
	StdOutPath string

	// StdErrPath, if non-blank, is the path to a file that the STDERR of the
	// executed Cmd should be redirected to.
	StdErrPath string

	// OutputBinding is the output binding specified in the CWL for this
	// command.
	OutputBinding Outputs
}

// String allows for pretty-printing of a Command.
func (c Command) String() string {
	return fmt.Sprintf("{\n Step: %s\n Cmd: %s\n Cwd: %s\n TmpPrefix: %s\n Env: %s\n StdIn: %s\n StdOut: %s\n StdErr: %s\n}", c.ID, c.Cmd, c.Cwd, c.TmpPrefix, c.Env, c.StdInPath, c.StdOutPath, c.StdErrPath)
}

// Execute runs the Command's Cmd in the right Cwd, with $HOME and $TMPDIR set
// as per Cwd and TmpPrefix, and with the environment variables from Env. The
// unique tmp dir is deleted afterwards. STDIN, OUT and ERR are also handled.
// Requiremnts are taken care of prior to execution.
// The return value is the decoded JSON of the file "cwl.output.json" created by
// Cmd in Cwd, if any. Otherwise it is the Outputs value.
func (c *Command) Execute() (interface{}, error) {
	//fmt.Printf("cmd: %+v\n", c.Cmd)

	if _, err := os.Stat(c.TmpPrefix); err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(c.TmpPrefix, 0700)
		if err != nil {
			return nil, err
		}
	}
	tmpDir, err := ioutil.TempDir(c.TmpPrefix, "cwlgo.tmp")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	if _, err := os.Stat(c.Cwd); err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(c.Cwd, 0700)
		if err != nil {
			return nil, err
		}
	}
	var cmdArgs []string
	if len(c.Cmd) > 1 {
		cmdArgs = c.Cmd[1:]
	}
	cmd := exec.Command(c.Cmd[0], cmdArgs...) // #nosec
	cmd.Dir = c.Cwd
	cmd.Env = append(c.Env, "HOME="+c.Cwd, "TMPDIR="+tmpDir, "PATH="+os.Getenv("PATH")) // *** no PATH in container

	if c.StdOutPath != "" {
		outfile, err := os.Create(filepath.Join(c.Cwd, c.StdOutPath))
		if err != nil {
			return nil, err
		}
		defer outfile.Close()
		cmd.Stdout = outfile
	}

	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	err = cmd.Wait()
	if err != nil {
		return nil, err
	}

	// return contents of cwl.output.json if it exists
	out := make(map[string]interface{})
	if jsonFile, err := os.Open(filepath.Join(c.Cwd, "cwl.output.json")); err == nil {
		defer jsonFile.Close()
		b, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(b, &out)
		return out, err
	}

	// otherwise, resolve the output binding
	for _, o := range c.OutputBinding {
		result, err := o.Resolve(c.Cwd)
		if err != nil {
			return nil, err
		}
		out[o.ID] = result
	}
	return out, nil
}

// Commands is a slice of Command.
type Commands []*Command

// String allows for pretty-printing of Commands.
func (cs Commands) String() string {
	var strs []string
	for _, c := range cs {
		strs = append(strs, c.String())
	}
	return strings.Join(strs, "\n")
}
