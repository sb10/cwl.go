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
	"fmt"
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

	// Cwd is the directory you should run the Cmd in.
	Cwd string

	// StdInPath, if non-blank, is the path to a file that should be piped in to
	// Cmd when executed.
	StdInPath string

	// StdOutPath, if non-blank, is the path to a file that the STDOUT of the
	// executed Cmd should be redirected to.
	StdOutPath string

	// StdErrPath, if non-blank, is the path to a file that the STDERR of the
	// executed Cmd should be redirected to.
	StdErrPath string
}

// String allows for pretty-printing of a Command.
func (c Command) String() string {
	return fmt.Sprintf("{\n Step: %s\n Cmd: %s\n Cwd: %s\n StdIn: %s\n StdOut: %s\n StdErr: %s\n}", c.ID, c.Cmd, c.Cwd, c.StdInPath, c.StdOutPath, c.StdErrPath)
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
