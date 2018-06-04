// This file is part of cwl.go.
// Author: Sendu Bala <sb10@sanger.ac.uk>.
//
// Copyright © 2018 Genome Research Limited
//
// Initially based on github.com/otiai10/yacle/core/handle.go,
// Copyright 2017 otiai10 (Hiromu OCHIAI)
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
	"os"
	"sort"
	"strconv"
)

// ResolveConfig is used to configure a NewResolver(), specifying runtime
// behaviour.
type ResolveConfig struct {
	// OutputDir is the output directory. Defaults to current directory.
	OutputDir string

	// Env are the os.Environ() format environment variable to keep when
	// executing resolved Commands. Defaults to none.
	Env []string

	// TmpDirPrefix is the path prefix for temporary directories. Defaults to
	// /tmp.
	TmpDirPrefix string

	// TmpOutDirPrefix is the path prefix for intermediate output directories.
	// Defaults to /tmp.
	TmpOutDirPrefix string

	// IntermediateOutputHandling determins what should happen to output files
	// in intermediate output directories. Possible values are:
	// "move" (default) == move files to OutputDir and delete intermediate
	//                     output directories.
	// "leave" == leave output files in intermediate output directories.
	// "copy" == copy files to OutputDir, delete nothing.
	IntermediateOutputHandling string

	// IntermediateTmpHandling determins what should happen to intermediate tmp
	// directories. Possible values are:
	// "rm" (default) == delete them.
	// "leave" == do not delete them.
	IntermediateTmpHandling string

	// Cores is the number of CPU cores reserved for the tool process. Defaults
	// to 1.
	Cores int

	// RAM is the amount of RAM in mebibytes (2**20) reserved for the tool
	// process.
	RAM int

	// OutputDirSize is the reserved storage space available in the designated
	// output directory.

	// TmpDirSize is the reserved storage space available in the designated
	// temporary directory.
}

// RuntimeValue returns the value in the ResolveConfig for the given runtime
// key, eg. 'runtime.cores'.
func (r ResolveConfig) RuntimeValue(key string) string {
	switch key {
	case "runtime.outdir":
		return r.OutputDir
	case "runtime.tmpdir":
		return r.TmpDirPrefix
	case "runtime.cores":
		cores := r.Cores
		if cores == 0 {
			cores = 1
		}
		return strconv.Itoa(cores)
	case "runtime.ram":
		ram := r.RAM
		if ram == 0 {
			return ""
		}
		return strconv.Itoa(ram)
	default:
		return ""
	}
}

// Resolver is a high-level struct with the logic for interpreting CWL.
type Resolver struct {
	Workflow   *Root
	Parameters Parameters
	Outdir     string
	Quiet      bool
	Config     ResolveConfig
	CWLDir     string
	ParamsDir  string
}

// NewResolver creates a new Resolver struct for the given pre-decoded Root. The
// path to the decoded CWL's directory must be provided to resolve relative paths.
func NewResolver(root *Root, config ResolveConfig, cwlDir string) (*Resolver, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return &Resolver{
		Workflow: root,
		Outdir:   cwd,
		Config:   config,
		CWLDir:   cwlDir,
	}, nil
}

// Resolve takes the pre-decoded parameters for a workflow and resolves
// everything to produce concrete commands to run. The path to the decoded
// param file's dir must be provided to resolve relative paths.
func (r *Resolver) Resolve(params Parameters, paramsDir string, ifc InputFileCallback) (Commands, error) {
	r.ParamsDir = paramsDir

	// resolve args
	r.Parameters = params
	arguments := r.resolveArguments()

	// resolve inputs
	priors, inputs, err := r.resolveInputs(paramsDir, r.CWLDir, ifc)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve required inputs: %v", err)
	}

	// handle defaults for our config
	cwd := r.Config.OutputDir
	if cwd == "" {
		cwd, err = os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get working directory: %s", err)
		}
		r.Config.OutputDir = cwd
	}

	tmpDirPrefix := r.Config.TmpDirPrefix
	if tmpDirPrefix == "" {
		tmpDirPrefix = "/tmp"
		r.Config.TmpDirPrefix = tmpDirPrefix
	}

	// create a concrete Command or recurse
	var cmds Commands
	if len(r.Workflow.BaseCommands) > 0 {
		args := append(priors, append(arguments, inputs...)...)
		cc := &Command{
			ID:         r.Workflow.ID,
			Cmd:        append(r.Workflow.BaseCommands[0:], args...),
			Cwd:        cwd,
			TmpPrefix:  tmpDirPrefix,
			StdInPath:  r.Workflow.Stdin,
			StdOutPath: r.Workflow.Stdout,
			StdErrPath: r.Workflow.Stderr,
		}
		cmds = append(cmds, cc)
	} else {
		for _, step := range r.Workflow.Steps {
			subR, err := NewResolver(step.Run.Workflow, r.Config, r.CWLDir)
			if err != nil {
				return nil, err
			}

			stepParams := r.resolveStepParams(step.In)

			subCmds, err := subR.Resolve(stepParams, paramsDir, ifc)
			if err != nil {
				return nil, err
			}
			for _, c := range subCmds {
				if c.ID == "" {
					c.ID = step.ID
				} else {
					c.ID = step.ID + "." + c.ID
				}
			}
			cmds = append(cmds, subCmds...)
		}
	}

	return cmds, nil
}

// resolveStepParams compares the given step inputs to the stored user
// parameters and sets values as appropriate.
func (r *Resolver) resolveStepParams(ins StepInputs) Parameters {
	stepParams := *NewParameters()
	for _, in := range ins {
		for _, source := range in.Source {
			if val, exists := r.Parameters[source]; exists {
				stepParams[in.ID] = val
			}
		}
	}
	return stepParams
}

// resolveArguments resolves workflow arguments with "valueFrom" properties
// against the config.
func (r *Resolver) resolveArguments() []string {
	var result []string
	sort.Sort(r.Workflow.Arguments)
	for i, arg := range r.Workflow.Arguments {
		if arg.Binding != nil && arg.Binding.ValueFrom != nil {
			r.Workflow.Arguments[i].Value = r.Config.RuntimeValue(arg.Binding.ValueFrom.Key())
		}
		result = append(result, r.Workflow.Arguments[i].Flatten()...)
	}
	return result
}

// resolveInputs resolves each workflow input to get the concrete command line
// arguments.
func (r *Resolver) resolveInputs(paramsDir, cwlDir string, ifc InputFileCallback) (priors []string, result []string, err error) {
	defer func() {
		if i := recover(); i != nil {
			err = fmt.Errorf("failed to resolve required inputs against provided params: %v", i)
		}
	}()

	sort.Sort(r.Workflow.Inputs)

	for _, in := range r.Workflow.Inputs {
		err = r.resolveInput(&in)
		if err != nil {
			return priors, result, err
		}
		if in.Binding == nil {
			continue
		}

		theseIns := in.Flatten(paramsDir, cwlDir, ifc)
		if in.Binding.Position < 0 {
			priors = append(priors, theseIns...)
		} else {
			result = append(result, theseIns...)
		}
	}
	return priors, result, nil
}

// resolveInput considers user parameters and defaults to decide on a concrete
// command line argument.
func (r *Resolver) resolveInput(input *Input) error {
	if provided, ok := r.Parameters[input.ID]; ok {
		input.Provided = provided
	}

	if input.Default == nil && input.Binding == nil && input.Provided == nil {
		return fmt.Errorf("input `%s` doesn't have default field but not provided", input.ID)
	}

	if key, needed := input.Types[0].NeedRequirement(); needed {
		for _, req := range r.Workflow.Requirements {
			for _, requiredtype := range req.Types {
				if requiredtype.Name == key {
					input.RequiredType = &requiredtype
					input.Requirements = r.Workflow.Requirements
				}
			}
		}
	}

	return nil
}
