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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
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

	// ParentOutput is non-nil if the Workflow this Command is part of has an
	// output with a source of this Command's output.
	ParentOutput *Output

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
// priorOuts is the return value of GetPriorOutputs() called on the same
// Resolver that made this Command.
//
// The return value is the decoded JSON of the file "cwl.output.json" created by
// Cmd in Cwd, if any. Otherwise it is the result of resolving the output
// binding. You should call SetOutput() with this returned value on the Resolver
// that made this Command.
//
// For Commands that are for ExpressionTools, instead of running a command line,
// it just interprets the expression to fill in the output.
func (c *Command) Execute(priorOuts map[string]map[string]interface{}) (map[string]interface{}, error) {
	// resolve args
	arguments, shellQuote := c.Workflow.Arguments.Resolve(c.Config)

	// resolve input that comes from prior step outputs
	for key, val := range priorOuts {
		child := filepath.Base(key)

		for pkey, pval := range c.Parameters {
			if p, ok := pval.(string); ok {
				parts := strings.Split(p, "/")
				if len(parts) > 1 && parts[len(parts)-2] == child {
					if out, exists := val[parts[len(parts)-1]]; exists {
						var adjustedOut interface{}
						if file, ok := out.(map[interface{}]interface{}); ok {
							if t, exists := file[fieldClass]; exists && t == typeFile {
								if file[fieldLocation] != "" && !filepath.IsAbs(file[fieldLocation].(string)) {
									// figure out the location of the prior
									// steps output dir based on it being
									// similar to our own *** not sure how
									// reliable this is; can't it just always
									// record an absolute path??
									priorOutputDir := strings.Replace(c.Config.OutputDir, c.Name, key, 1)
									fileCopy := make(map[interface{}]interface{})
									for fk, fv := range file {
										fileCopy[fk] = fv
									}
									fileCopy[fieldLocation] = filepath.Join(priorOutputDir, file[fieldLocation].(string))
									adjustedOut = fileCopy
								}
							}
						}
						if adjustedOut == nil {
							adjustedOut = out
						}
						c.Parameters[pkey] = adjustedOut
					}
				}
			}
		}
	}

	// resolve inputs
	var nestedOutIndexFound bool
	var nestedOutIndex [2]int
	if val, exists := c.Parameters[scatterNestedInput]; exists {
		nestedOutIndex = val.([2]int)
		nestedOutIndexFound = true
	}
	var flatOutIndexFound bool
	var flatOutIndex int
	if val, exists := c.Parameters[scatterFlatInput]; exists {
		flatOutIndex = val.(int)
		flatOutIndexFound = true
	}
	inputs, err := c.Workflow.Inputs.Resolve(c.Workflow.Requirements, c.Parameters, c.ParamsDir, c.CWLDir, c.Name, c.IFC, c.InputContext, otto.New())
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
	vm, viaShell, reqEnv, err := c.resolveRequirments()
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

		// *** should only be putting things in out where key matches something
		// in the outputs
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

		if c.ParentOutput != nil && len(c.ParentOutput.Source) == 1 {
			key := filepath.Base(c.ParentOutput.Source[0])
			if val, exists := out[key]; exists {
				parentStep := filepath.Dir(c.Name)
				parentOut := make(map[string]interface{})
				parentOut[c.ParentOutput.ID] = val
				out[parentStep] = parentOut
			}
		}

		return out, nil
	}

	// CommandLineTool
	cmdStrs := c.Workflow.BaseCommands
	arguments = append(arguments, inputs...)
	sort.Sort(arguments)

	// if no basecommands, we use the first thing out of args or inputs as the
	// base command
	if len(cmdStrs) == 0 && len(arguments) > 0 {
		cmdStrs = arguments[0].arg
		if len(arguments) > 1 {
			arguments = arguments[1:]
		} else {
			arguments = []*SortableArg{}
		}
	}

	if len(cmdStrs) == 0 {
		return nil, fmt.Errorf("no command specified")
	}

	// create a concrete command line to run
	args := arguments.flatten()
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

	if len(reqEnv) > 0 {
		cmd.Env = append(cmd.Env, reqEnv...)
	}

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

	var stderr bytes.Buffer
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
	} else {
		cmd.Stderr = &stderr
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
		// fmt.Printf("wait err %s for cmd %s (%s)\n", err, resolvedCmd, stderr.String())
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
		return out, err
	}

	// otherwise, resolve the output binding
	for _, o := range c.Workflow.Outputs {
		result, errr := o.Resolve(c.Config.OutputDir, stdOutPath, stdErrPath, vm)
		if errr != nil {
			return nil, errr
		}

		if nestedOutIndexFound {
			sliceResult := make([]interface{}, nestedOutIndex[0]+1)
			nestedResult := make([]interface{}, nestedOutIndex[1]+1)
			nestedResult[nestedOutIndex[1]] = result
			sliceResult[nestedOutIndex[0]] = nestedResult
			out[o.ID] = sliceResult
		} else if flatOutIndexFound {
			sliceResult := make([]interface{}, flatOutIndex+1)
			sliceResult[flatOutIndex] = result
			out[o.ID] = sliceResult
		} else {
			out[o.ID] = result
		}
	}

	return out, err
}

// resolveRequirments handles things like InlineJavascriptRequirement, creates
// files specified in InitialWorkDirRequirement, and returns a javascript vm for
// resolving expressions, a bool to say if command should be run via shell, and
// any environment variables that the command should run with (in os.Environ()
// format).
func (c *Command) resolveRequirments() (*otto.Otto, bool, []string, error) {
	// set up our javascript interpreter, first dealing with imports
	underscore.Disable()
	for _, req := range c.Workflow.Requirements {
		switch req.Class {
		case reqJS:
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
		return nil, false, nil, err
	}
	err = vm.Set("runtime", map[string]string{
		"outdir": c.Config.OutputDir,
		"tmpdir": c.Config.TmpDirPrefix,
		"cores":  c.Config.RuntimeValue("runtime.cores"),
		"ram":    c.Config.RuntimeValue("runtime.ram"),
	})
	if err != nil {
		return nil, false, nil, err
	}

	// handle requirements
	viaShell := false
	var env []string
	setEnvs := make(map[string]bool)
	for _, req := range c.Workflow.Requirements {
		switch req.Class {
		case reqShell:
			viaShell = true
		case reqJS:
			// parse expressions
			for _, jse := range req.ExpressionLib {
				if jse.Kind == "$execute" {
					_, err := vm.Run(jse.Value)
					if err != nil {
						return vm, viaShell, env, err
					}
				}
			}
		case reqWorkDir:
			for _, entry := range req.Listing {
				basename, _, _, err := evaluateExpression(entry.EntryName, vm)
				if err != nil {
					return vm, viaShell, env, err
				}
				e := entry.Entry
				contents, _, obj, err := evaluateExpression(e, vm)
				if err != nil {
					return vm, viaShell, env, err
				}

				if len(contents) > 0 {
					err = ioutil.WriteFile(filepath.Join(c.Config.OutputDir, basename), []byte(contents), 0600)
					if err != nil {
						return vm, viaShell, env, err
					}
				} else if obj != nil {
					if val, err := obj.Get(fieldPath); err == nil && val.IsString() {
						path, err := val.ToString()
						if err != nil {
							return vm, viaShell, env, err
						}

						err = os.MkdirAll(c.Config.OutputDir, 0700)
						if err != nil {
							return vm, viaShell, env, err
						}

						err = copyFile(path, filepath.Join(c.Config.OutputDir, basename))
						if err != nil {
							return vm, viaShell, env, err
						}
					}
				}
			}
		case reqEnv:
			for _, ed := range req.EnvDef {
				val, _, _, err := evaluateExpression(ed.Value, vm)
				if err != nil {
					return vm, viaShell, env, err
				}
				env = append(env, fmt.Sprintf("%s=%s", ed.Name, val))
				setEnvs[ed.Name] = true
			}
		}
	}

	// handle hints
	for _, hint := range c.Workflow.Hints {
		switch hint.Class {
		case reqEnv:
			for _, ed := range hint.Envs {
				if !setEnvs[ed.Name] {
					val, _, _, err := evaluateExpression(ed.Value, vm)
					if err != nil {
						return vm, viaShell, env, err
					}
					env = append(env, fmt.Sprintf("%s=%s", ed.Name, val))
				}
			}
		}
	}

	return vm, viaShell, env, nil
}

// evaluateExpression evaluates the given string as javascript if it starts with
// $( or ${, and returns either a string, float or an object. If not javascript,
// just returns e unchanged.
func evaluateExpression(e string, vm *otto.Otto) (string, float64, *otto.Object, error) {
	e = strings.TrimSpace(e)

	var fl float64
	if strings.HasPrefix(e, "${") {
		// a js function body, which we hope is always specified as a string
		// starting with ${ and ending with }
		e = strings.TrimPrefix(e, "${")
		e = strings.TrimSuffix(e, "}")
		return evaluateJS(e, vm)
	} else if strings.Contains(e, "$(") {
		// looks like the string contains js expressions; there could be
		// multiple of them in this string, and they could contain nested
		// unmatched parentheses, so we do our best to try and pull these out
		// and replace them with the evaluated result
		var replacements [][2]string
		var startIndex, openParen int
		var inSingleQuote, inDoubleQuote bool
		var singleNum float64
		for pos := 0; pos < len(e); pos++ {
			// *** not yet ignoring stuff in comments...
			switch string(e[pos]) {
			case "$":
				if startIndex == 0 && string(e[pos+1]) == "(" {
					startIndex = pos + 2
				}
			case "'":
				if startIndex != 0 && pos > 0 && string(e[pos-1]) != escapeStr {
					if inSingleQuote {
						inSingleQuote = false
					} else if !inDoubleQuote {
						inSingleQuote = true
					}
				}
			case `"`:
				if startIndex != 0 && pos > 0 && string(e[pos-1]) != escapeStr {
					if inDoubleQuote {
						inDoubleQuote = false
					} else if !inSingleQuote {
						inDoubleQuote = true
					}
				}
			case "(":
				if startIndex != 0 && pos > 0 && string(e[pos-1]) != escapeStr && !inSingleQuote && !inDoubleQuote {
					openParen++
				}
			case ")":
				if startIndex != 0 && pos > 0 && string(e[pos-1]) != escapeStr && !inSingleQuote && !inDoubleQuote {
					openParen--
					if openParen == 0 {
						thisE := e[startIndex:pos]
						thisStr, thisNum, thisObj, thisErr := evaluateJS(thisE, vm)
						if thisObj != nil {
							// can't replace an object in to a string, assume
							// that all of e is for this obj and return now
							return thisStr, thisNum, thisObj, thisErr
						}

						if thisStr == "" {
							singleNum = thisNum
							thisStr = strconv.FormatFloat(thisNum, 'f', -1, 64)
						}
						replacements = append(replacements, [2]string{"$(" + thisE + ")", thisStr})

						startIndex = 0
					}
				}
			}
		}

		if len(replacements) == 1 && singleNum != 0 {
			return "", singleNum, nil, nil
		}

		for _, r := range replacements {
			e = strings.Replace(e, r[0], r[1], 1)
		}
		return e, fl, nil, nil
	}
	return e, fl, nil, nil
}

// evaluateJS evaluates javascript expressions and function bodies
func evaluateJS(e string, vm *otto.Otto) (string, float64, *otto.Object, error) {
	var fl float64

	value, err := vm.Run(e)
	if err != nil {
		if strings.Contains(err.Error(), "Unexpected token :") {
			// might be just a bare object, try again by assigning to a
			// variable
			e = "$cwlgoreturnval = " + e
			value, err = vm.Run(e)
		} else if strings.Contains(err.Error(), "Illegal return statement") {
			// might be returning something, wrap in a function
			e = "function cwlgoreturnfunc() { " + e + " }; cwlgoreturnfunc()"
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
	return e, fl, nil, nil
}

// fileToSelf returns a copy of the given file with path, basename, nameext and
// nameroot filled in, suitable for setting as "self" when evauluating
// expressions. path is the absolute path to the file.
func fileToSelf(path string, file map[string]interface{}) map[string]interface{} {
	self := make(map[string]interface{})
	for key, val := range file {
		self[key] = val
	}
	self[fieldPath] = path
	self[fieldBasename] = filepath.Base(path)
	self[fieldNameExt] = filepath.Ext(path)
	self[fieldNameRoot] = strings.TrimSuffix(self[fieldBasename].(string), self[fieldNameExt].(string))
	if self[fieldNameRoot].(string) == "" {
		self[fieldNameRoot] = self[fieldBasename]
	}
	return self
}

// ottoObjToFiles takes an otto object representing either a File or an array of
// paths or Files, and returns a slice of Files. Dir is supplied in order to
// determine relative location paths.
func ottoObjToFiles(obj *otto.Object, dir string) ([]interface{}, error) {
	l, err := obj.Get(fieldLocation)
	if err != nil {
		return nil, err
	}

	var location string
	if l.IsDefined() {
		location, err = l.ToString()
		if err != nil {
			return nil, err
		}
	}

	if location == "" {
		p, errg := obj.Get(fieldPath)
		if errg != nil {
			return nil, errg
		}

		if p.IsDefined() {
			thisPath, errt := p.ToString()
			if errt != nil {
				return nil, errt
			}

			location, err = filepath.Rel(dir, thisPath)
			if err != nil {
				return nil, err
			}
		}
	}

	var files []interface{}
	if location == "" {
		for _, key := range obj.Keys() {
			val, _ := obj.Get(key)
			if val.IsObject() {
				// recurse
				subObj := val.Object()
				theseFiles, errr := ottoObjToFiles(subObj, dir)
				if errr != nil {
					return nil, errr
				}
				files = append(files, theseFiles...)
			} else if val.IsString() {
				sFile := make(map[interface{}]interface{})
				sFile[fieldClass] = typeFile
				sFile[fieldLocation], err = val.ToString()
				if err != nil {
					return nil, err
				}
				files = append(files, sFile)
			}
		}
	} else {
		sFile := make(map[interface{}]interface{})
		sFile[fieldClass] = typeFile
		sFile[fieldLocation] = location
		files = append(files, sFile)
	}

	return files, err
}

// Commands is a slice of Command.
type Commands []*Command

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
	defer func() {
		err = in.Close()
	}()
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
