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
	"os"
	"path/filepath"
	"sort"
)

// Resolver is a high-level struct with the logic for interpreting CWL.
type Resolver struct {
	Workflow   *Root
	Parameters Parameters
	Outdir     string
	Quiet      bool
	Alias      map[string]interface{}
}

// NewResolver creates a new Resolver struct for the given pre-decoded Root.
func NewResolver(root *Root) (*Resolver, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return &Resolver{
		Workflow: root,
		Outdir:   cwd,
		Alias:    map[string]interface{}{},
	}, nil
}

// Resolve takes the pre-decoded parameters for a workflow and resolves
// everything to produce concrete commands to run.
func (r *Resolver) Resolve(params Parameters) (Commands, error) {
	r.Parameters = params
	arguments := r.resolveArguments()

	priors, inputs, err := r.resolveInputs()
	if err != nil {
		return nil, fmt.Errorf("failed to resolve required inputs: %v", err)
	}

	args := append(priors, append(arguments, inputs...)...)

	var cmds Commands
	if len(r.Workflow.BaseCommands) > 0 {
		cwd, err := filepath.Abs(filepath.Dir(r.Workflow.Path))
		if err != nil {
			return nil, err
		}
		cc := &Command{
			ID:         r.Workflow.ID,
			Cmd:        append(r.Workflow.BaseCommands[0:], args...),
			Cwd:        cwd,
			StdInPath:  r.Workflow.Stdin,
			StdOutPath: r.Workflow.Stdout,
			StdErrPath: r.Workflow.Stderr,
		}
		cmds = append(cmds, cc)
	} else {
		for _, step := range r.Workflow.Steps {
			subR, err := NewResolver(step.Run.Workflow)
			if err != nil {
				return nil, err
			}

			stepParams := r.resolveStepParams(step.In)

			subCmds, err := subR.Resolve(stepParams)
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
// against aliases.
func (r *Resolver) resolveArguments() []string {
	var result []string
	sort.Sort(r.Workflow.Arguments)
	for i, arg := range r.Workflow.Arguments {
		if arg.Binding != nil && arg.Binding.ValueFrom != nil {
			r.Workflow.Arguments[i].Value = r.AliasFor(arg.Binding.ValueFrom.Key())
		}
		result = append(result, r.Workflow.Arguments[i].Flatten()...)
	}
	return result
}

// resolveInputs resolves each workflow input to get the concrete command line
// arguments.
func (r *Resolver) resolveInputs() (priors []string, result []string, err error) {
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
		if in.Binding.Position < 0 {
			priors = append(priors, in.Flatten()...)
		} else {
			result = append(result, in.Flatten()...)
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

	if input.Binding == nil {
		input.Binding = &Binding{}
	}

	return nil
}

// AliasFor returns what the given key is an alias for, if anything.
func (r *Resolver) AliasFor(key string) string {
	v, ok := r.Alias[key]
	if !ok {
		return ""
	}
	return fmt.Sprintf("%v", v)
}
