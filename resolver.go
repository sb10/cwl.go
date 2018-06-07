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
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/robertkrimen/otto"
	"github.com/robertkrimen/otto/underscore"
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
	Workflow     *Root
	Parameters   Parameters
	Outdir       string
	Quiet        bool
	Config       ResolveConfig
	CWLDir       string
	ParamsDir    string
	inputContext map[string]interface{}
}

// NewResolver creates a new Resolver struct for the given pre-decoded Root. The
// path to the decoded CWL's directory must be provided to resolve relative paths.
func NewResolver(root *Root, config ResolveConfig, cwlDir string) (*Resolver, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return &Resolver{
		Workflow:     root,
		Outdir:       cwd,
		Config:       config,
		CWLDir:       cwlDir,
		inputContext: make(map[string]interface{}),
	}, nil
}

// Resolve takes the pre-decoded parameters for a workflow and resolves
// everything to produce concrete commands to run. The path to the decoded
// param file's dir must be provided to resolve relative paths.
//
// Also resolves any requirments, carrying out anything actionble, which may
// involve creating files according to an InitialWorkDirRequirement.
//
// The returned Otto can be used if you wish to Execute() any of the Commands.
func (r *Resolver) Resolve(params Parameters, paramsDir string, ifc InputFileCallback) (Commands, *otto.Otto, error) {
	r.ParamsDir = paramsDir

	// resolve args
	r.Parameters = params
	arguments, shellQuote := r.resolveArguments()

	// resolve inputs
	priors, inputs, err := r.resolveInputs(paramsDir, r.CWLDir, ifc)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to resolve required inputs: %v", err)
	}

	// handle defaults for our config
	cwd := r.Config.OutputDir
	if cwd == "" {
		cwd, err = os.Getwd()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get working directory: %s", err)
		}
		r.Config.OutputDir = cwd
	}

	tmpDirPrefix := r.Config.TmpDirPrefix
	if tmpDirPrefix == "" {
		tmpDirPrefix = "/tmp"
		r.Config.TmpDirPrefix = tmpDirPrefix
	}

	// resolve requirments
	vm, viaShell, err := r.resolveRequirments()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to resolve requirements: %s", err)
	}

	if r.Workflow.Class == "ExpressionTool" {
		// no command to run, but we still want to return a "Command" so the
		// user can Execute() it and get the output
		return Commands{&Command{
			ID:            r.Workflow.ID,
			Expression:    r.Workflow.Expression,
			OutputBinding: r.Workflow.Outputs,
		}}, vm, nil
	}

	// if no basecommands, we use the first thing out of args or inputs as the
	// base command
	cmdStrs := r.Workflow.BaseCommands
	if len(cmdStrs) == 0 {
		if len(priors) > 0 {
			cmdStrs = priors
			priors = []string{}
		} else if len(arguments) > 0 {
			cmdStrs = arguments
			arguments = []string{}
		} else if len(inputs) > 0 {
			cmdStrs = inputs
			inputs = []string{}
		}
	}

	// create a concrete Command or recurse
	var cmds Commands
	if len(cmdStrs) > 0 {
		args := append(priors, append(arguments, inputs...)...)
		stdinPath, _, _, err := evaluateExpression(r.Workflow.Stdin, vm)
		if err != nil {
			return nil, nil, err
		}
		stdoutPath, _, _, err := evaluateExpression(r.Workflow.Stdout, vm)
		if err != nil {
			return nil, nil, err
		}
		stderrPath, _, _, err := evaluateExpression(r.Workflow.Stderr, vm)
		if err != nil {
			return nil, nil, err
		}
		cc := &Command{
			ID:            r.Workflow.ID,
			Cmd:           append(cmdStrs[0:], args...),
			ViaShell:      viaShell,
			ShellQuote:    shellQuote,
			Cwd:           cwd,
			TmpPrefix:     tmpDirPrefix,
			StdInPath:     stdinPath,
			StdOutPath:    stdoutPath,
			StdErrPath:    stderrPath,
			OutputBinding: r.Workflow.Outputs,
		}
		cmds = append(cmds, cc)
	} else {
		for _, step := range r.Workflow.Steps {
			subR, err := NewResolver(step.Run.Workflow, r.Config, r.CWLDir)
			if err != nil {
				return nil, nil, err
			}

			stepParams := r.resolveStepParams(step.In)

			subCmds, _, err := subR.Resolve(stepParams, paramsDir, ifc)
			if err != nil {
				return nil, nil, err
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

	return cmds, vm, nil
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
// against the config. Returns a slice of command line arguements and a bool,
// which if true means shell metacharacters should be quoted.
func (r *Resolver) resolveArguments() ([]string, bool) {
	var result []string
	var shellQuote bool
	sort.Sort(r.Workflow.Arguments)
	for i, arg := range r.Workflow.Arguments {
		if arg.Binding != nil && arg.Binding.ValueFrom != nil {
			// *** need to properly evaluate this if an expression?
			str := arg.Binding.ValueFrom.string
			if strings.HasPrefix(str, "$(") {
				r.Workflow.Arguments[i].Value = r.Config.RuntimeValue(arg.Binding.ValueFrom.Key())
			} else {
				r.Workflow.Arguments[i].Value = str
			}

			if arg.Binding.ShellQuote {
				shellQuote = true
			}
		}
		result = append(result, r.Workflow.Arguments[i].Flatten()...)
	}
	return result, shellQuote
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
		err = r.resolveInput(in)
		if err != nil {
			return priors, result, err
		}

		theseIns, err := in.Flatten(r.inputContext, paramsDir, cwlDir, ifc)
		if err != nil {
			return priors, result, err
		}

		if in.Binding == nil {
			continue
		}
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

// resolveRequirments handles things like InlineJavascriptRequirement, creates
// files specified in InitialWorkDirRequirement, and returns a javascript vm for
// resolving expressions, and a bool to say if command should be run via shell.
func (r *Resolver) resolveRequirments() (*otto.Otto, bool, error) {
	// set up our javascript interpreter, first dealing with imports
	underscore.Disable()
	for _, req := range r.Workflow.Requirements {
		switch req.Class {
		case "InlineJavascriptRequirement":
			for _, jse := range req.ExpressionLib {
				if jse.Kind == "$include" {
					if jse.Value == "underscore.js" {
						underscore.Enable()
					}
					// *** else, how to import arbitrary js packages?...
				}
			}
		}
	}
	vm := otto.New()

	// set up namespace context
	vm.Set("inputs", r.inputContext)
	vm.Set("runtime", map[string]string{
		"outdir": r.Config.OutputDir,
		"tmpdir": r.Config.TmpDirPrefix,
		"cores":  r.Config.RuntimeValue("runtime.cores"),
		"ram":    r.Config.RuntimeValue("runtime.ram"),
	})

	viaShell := false
	for _, req := range r.Workflow.Requirements {
		switch req.Class {
		case "ShellCommandRequirement":
			viaShell = true
		case "InlineJavascriptRequirement":
			// parse expressions
			for _, jse := range req.ExpressionLib {
				if jse.Kind == "$execute" {
					_, err := vm.Run(jse.Value)
					if err != nil {
						return vm, viaShell, err
					}
				}
			}
		case "InitialWorkDirRequirement":
			for _, entry := range req.Listing {
				basename := entry.EntryName
				e := entry.Entry
				contents, _, _, err := evaluateExpression(e, vm)
				if err != nil {
					return vm, viaShell, err
				}

				err = ioutil.WriteFile(filepath.Join(r.Config.OutputDir, basename), []byte(contents), 0600)
				if err != nil {
					return vm, viaShell, err
				}
			}
		}
	}
	return vm, viaShell, nil
}

func trimExpression(e string) string {
	if strings.HasPrefix(e, "$(") {
		e = strings.TrimPrefix(e, "$(")
		e = strings.TrimSuffix(e, ")")
	} else if strings.HasPrefix(e, "${") {
		e = strings.TrimPrefix(e, "${")
		e = strings.TrimSuffix(e, "}")
	}
	return e
}

// evaluateExpression evaluates the given string as javascript if it starts with
// $( or ${, and returns either a string, float or an object. If not javascript,
// just returns e unchanged.
func evaluateExpression(e string, vm *otto.Otto) (string, float64, *otto.Object, error) {
	var fl float64
	if strings.HasPrefix(e, "$") {
		e = trimExpression(e)

		// evaluate as javascript
		value, err := vm.Run(e)
		if err != nil {
			if strings.Contains(err.Error(), "Unexpected token :") {
				// might be just a bare object, try again by assigning to a
				// variable
				e = "$cwlgoreturnval = " + e
				value, err = vm.Run(e)
			}
			if err != nil {
				fmt.Printf("got evaluateExpression err [%s] for [%s]\n", err, e)
				return "", fl, nil, err
			}
		}

		switch {
		case value.IsNumber():
			f, err := value.ToFloat()
			if err != nil {
				return "", fl, nil, err
			}
			return "", f, nil, nil
		case value.IsString():
			v, err := value.ToString()
			if err != nil {
				return "", fl, nil, err
			}
			return v, fl, nil, nil
		case value.IsObject():
			return "", fl, value.Object(), nil
		}
	}
	return e, fl, nil, nil
}
