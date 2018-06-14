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
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/robertkrimen/otto"
	"github.com/robertkrimen/otto/underscore"
)

// Command is a high-level interpretation of a concrete command line that needs
// to be run as part of a workflow.
type Command struct {
	// UniqueID is a generated string that identifies the command as coming from
	// a (particular step in a) particular workflow, resolved against particular
	// parameters.
	UniqueID string

	// Dependencies are the IDs of other Commands that you must arrange to
	// Execute() before Execute()ing this one.
	Dependencies []string

	*Resolver
}

// NewCommand creates a NewCommand that can be stored and later Execute()d.
func NewCommand(uniqueID string, dependencies []string, resolver *Resolver) *Command {
	return &Command{
		UniqueID:     uniqueID,
		Dependencies: dependencies,
		Resolver:     resolver,
	}
}

// Execute runs the Command's command line in the right working directory, with
// $HOME,  $TMPDIR and environment variables set as per configuration. The
// unique tmp dir is deleted afterwards. STDIN, OUT and ERR are also handled.
// Requiremnts are taken care of prior to execution.
//
// The return value is the decoded JSON of the file "cwl.output.json" created by
// Cmd in Cwd, if any. Otherwise it is the result of resolving the output
// binding.
//
// For Commands that are for ExpressionTools, instead of running a command line,
// it just interprets the expression to fill in the output.
func (c *Command) Execute() (interface{}, error) {
	// resolve args
	arguments, shellQuote := c.Workflow.Arguments.Resolve(c.Config)

	// resolve input that comes from prior step outputs
	for key, val := range c.OutputContext {
		//parent := filepath.Dir(key)
		child := filepath.Base(key)

		for pkey, pval := range c.Parameters {
			if p, ok := pval.(string); ok {
				parts := strings.Split(p, "/")
				if len(parts) > 1 && parts[len(parts)-2] == child {
					if out, exists := val[parts[len(parts)-1]]; exists {
						if file, ok := out.(map[interface{}]interface{}); ok {
							if t, exists := file[fieldClass]; exists && t == typeFile {
								if file[fieldLocation] != "" && !filepath.IsAbs(file[fieldLocation].(string)) {
									file[fieldLocation] = filepath.Join(c.Config.OutputDir, file[fieldLocation].(string))
									out = file
								}
							}
						}
						c.Parameters[pkey] = out
					}
				}
			}
		}
	}

	// resolve inputs
	priors, inputs, err := c.Workflow.Inputs.Resolve(c.Workflow.Requirements, c.Parameters, c.ParamsDir, c.CWLDir, c.IFC, c.InputContext)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve required inputs: %v", err)
	}

	// handle defaults for our config
	if c.Config.OutputDir == "" {
		cwd, errw := os.Getwd()
		if errw != nil {
			return nil, fmt.Errorf("failed to get working directory: %s", errw)
		}
		c.Config.OutputDir = cwd
	}

	if c.Config.TmpDirPrefix == "" {
		c.Config.TmpDirPrefix = "/tmp"
	}

	// resolve requirments
	vm, viaShell, err := c.resolveRequirments()
	if err != nil {
		return nil, fmt.Errorf("failed to resolve requirements: %s", err)
	}

	out := make(map[string]interface{})

	// ExpressionTool?
	if c.Workflow.Class == classExpression {
		if c.Workflow.Expression == "" {
			return nil, fmt.Errorf("no expression specified")
		}

		str, i, obj, errr := evaluateExpression(c.Workflow.Expression, vm)
		if errr != nil {
			return nil, errr
		}

		if str != "" || i != 0 {
			// *** not sure what to do in this case...
			return out, fmt.Errorf("expression tool returned string [%s] or numner [%f]", str, i)
		}

		for _, key := range obj.Keys() {
			val, errg := obj.Get(key)
			if errg != nil {
				return out, errg
			}
			switch {
			case val.IsNumber():
				f, errt := val.ToFloat()
				if errt != nil {
					return out, errt
				}
				i, errt := val.ToInteger()
				if errt != nil {
					return out, errt
				}
				if float64(i) == f {
					out[key] = int(i)
				} else {
					out[key] = f
				}
			case val.IsBoolean():
				v, errt := val.ToBoolean()
				if errt != nil {
					return out, errt
				}
				out[key] = v
			case val.IsString():
				v, errt := val.ToString()
				if errt != nil {
					return out, errt
				}
				out[key] = v
			}
		}

		c.OutputContext[c.UniqueID] = out
		return out, nil
	}

	// CommandLineTool
	// if no basecommands, we use the first thing out of args or inputs as the
	// base command
	cmdStrs := c.Workflow.BaseCommands
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

	if len(cmdStrs) == 0 {
		return nil, fmt.Errorf("no command specified")
	}

	// create a concrete command line to run
	args := append(priors, append(arguments, inputs...)...)
	stdInPath, _, _, err := evaluateExpression(c.Workflow.Stdin, vm)
	if err != nil {
		return nil, err
	}
	stdOutPath, _, _, err := evaluateExpression(c.Workflow.Stdout, vm)
	if err != nil {
		return nil, err
	}
	stdErrPath, _, _, err := evaluateExpression(c.Workflow.Stderr, vm)
	if err != nil {
		return nil, err
	}

	resolvedCmd := append(cmdStrs[0:], args...)

	// execute the Cmd
	if _, errs := os.Stat(c.Config.TmpDirPrefix); errs != nil && os.IsNotExist(errs) {
		errm := os.MkdirAll(c.Config.TmpDirPrefix, 0700)
		if errm != nil {
			return nil, errm
		}
	}
	tmpDir, err := ioutil.TempDir(c.Config.TmpDirPrefix, "cwlgo.tmp")
	if err != nil {
		return nil, err
	}
	defer func() {
		err = os.RemoveAll(tmpDir)
	}()

	if _, errs := os.Stat(c.Config.OutputDir); errs != nil && os.IsNotExist(errs) {
		errm := os.MkdirAll(c.Config.OutputDir, 0700)
		if errm != nil {
			return nil, errm
		}
	}
	var cmdArgs []string
	if len(resolvedCmd) > 1 {
		cmdArgs = resolvedCmd[1:]
	}

	var cmd *exec.Cmd
	if viaShell {
		cmdStr := strings.Join(append([]string{resolvedCmd[0]}, cmdArgs...), " ")
		if shellQuote {
			cmdStr = "'" + cmdStr + "'"
		}
		cmd = exec.Command("bash", "-c", cmdStr) // #nosec
	} else {
		cmd = exec.Command(resolvedCmd[0], cmdArgs...) // #nosec
	}

	cmd.Dir = c.Config.OutputDir
	cmd.Env = append(c.Config.Env, "HOME="+c.Config.OutputDir, "TMPDIR="+c.Config.TmpDirPrefix, "PATH="+os.Getenv("PATH")) // *** no PATH in container

	// handle stdout redirects
	var outFile *os.File
	if stdOutPath == "" && c.Workflow.Outputs[0].Types[0].Type == fieldStdOut {
		// this is a shortcut; StdOutPath should be set to a random file name
		outFile, err = ioutil.TempFile(c.Config.OutputDir, fieldStdOut)
		if err != nil {
			return nil, err
		}
		stdOutPath = filepath.Base(outFile.Name())
	}

	if stdOutPath != "" {
		if outFile == nil {
			outFile, err = os.Create(filepath.Join(c.Config.OutputDir, stdOutPath))
			if err != nil {
				return nil, err
			}
		}
		defer func() {
			err = outFile.Close()
		}()
		cmd.Stdout = outFile
	}

	// handle stderr redirects
	var errFile *os.File
	if stdErrPath == "" && c.Workflow.Outputs[0].Types[0].Type == fieldStdErr {
		// this is a shortcut; StdErrPath should be set to a random file name
		errFile, err = ioutil.TempFile(c.Config.OutputDir, fieldStdErr)
		if err != nil {
			return nil, err
		}
		stdErrPath = filepath.Base(errFile.Name())
	}

	if stdErrPath != "" {
		if errFile == nil {
			errFile, err = os.Create(filepath.Join(c.Config.OutputDir, stdErrPath))
			if err != nil {
				return nil, err
			}
		}
		defer func() {
			err = errFile.Close()
		}()
		cmd.Stderr = errFile
	}

	// handle stdin
	if stdInPath != "" {
		stdin, errs := cmd.StdinPipe()
		if errs != nil {
			return nil, errs
		}

		f, errs := os.Open(stdInPath)
		if errs != nil {
			return nil, errs
		}

		go func() {
			defer func() {
				err = stdin.Close()
			}()
			_, err = io.Copy(stdin, f)
		}()
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
	if jsonFile, erro := os.Open(filepath.Join(c.Config.OutputDir, "cwl.output.json")); erro == nil {
		defer func() {
			err = jsonFile.Close()
		}()
		b, errr := ioutil.ReadAll(jsonFile)
		if errr != nil {
			return nil, errr
		}

		err = json.Unmarshal(b, &out)
		c.OutputContext[c.UniqueID] = out
		return out, err
	}

	// otherwise, resolve the output binding
	for _, o := range c.Workflow.Outputs {
		result, errr := o.Resolve(c.Config.OutputDir, stdOutPath, stdErrPath, vm)
		if errr != nil {
			return nil, errr
		}
		out[o.ID] = result
	}
	c.OutputContext[c.UniqueID] = out
	return out, err
}

// resolveRequirments handles things like InlineJavascriptRequirement, creates
// files specified in InitialWorkDirRequirement, and returns a javascript vm for
// resolving expressions, and a bool to say if command should be run via shell.
func (c *Command) resolveRequirments() (*otto.Otto, bool, error) {
	// set up our javascript interpreter, first dealing with imports
	underscore.Disable()
	for _, req := range c.Workflow.Requirements {
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
	err := vm.Set("inputs", c.InputContext)
	if err != nil {
		return nil, false, err
	}
	err = vm.Set("runtime", map[string]string{
		"outdir": c.Config.OutputDir,
		"tmpdir": c.Config.TmpDirPrefix,
		"cores":  c.Config.RuntimeValue("runtime.cores"),
		"ram":    c.Config.RuntimeValue("runtime.ram"),
	})
	if err != nil {
		return nil, false, err
	}

	viaShell := false
	for _, req := range c.Workflow.Requirements {
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

				err = ioutil.WriteFile(filepath.Join(c.Config.OutputDir, basename), []byte(contents), 0600)
				if err != nil {
					return vm, viaShell, err
				}
			}
		}
	}
	return vm, viaShell, nil
}

func trimExpression(e string) string {
	e = strings.TrimSpace(e)
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
			} else if strings.Contains(err.Error(), "Illegal return statement") && strings.HasPrefix(e, "return") {
				// might be returning an object, try again by assigning to a
				// variable instead of return
				e = "$cwlgoreturnval = " + strings.TrimPrefix(e, "return")
				value, err = vm.Run(e)
			}
			if err != nil {
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

// Commands is a slice of Command.
type Commands []*Command
