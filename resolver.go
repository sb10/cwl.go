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
	"path/filepath"
	"strconv"
	"strings"
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

	// IntermediateOutputHandling determines what should happen to output files
	// in intermediate output directories. Possible values are:
	// "move" (default) == move files to OutputDir and delete intermediate
	//                     output directories.
	// "leave" == leave output files in intermediate output directories.
	// "copy" == copy files to OutputDir, delete nothing.
	IntermediateOutputHandling string

	// IntermediateTmpHandling determines what should happen to intermediate tmp
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
	Name          string
	Workflow      *Root
	Parameters    Parameters
	Outdir        string
	Quiet         bool
	Config        ResolveConfig
	CWLDir        string
	ParamsDir     string
	IFC           InputFileCallback
	InputContext  map[string]interface{}
	OutputContext map[string]map[string]interface{}
}

// NewResolver creates a new Resolver struct for the given pre-decoded Root. The
// path to the decoded CWL's directory must be provided to resolve relative paths.
func NewResolver(root *Root, config ResolveConfig, cwlDir string, optionalOutputContext ...map[string]map[string]interface{}) (*Resolver, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	outputContext := make(map[string]map[string]interface{})
	if len(optionalOutputContext) == 1 {
		outputContext = optionalOutputContext[0]
	}

	return &Resolver{
		Workflow:      root,
		Outdir:        cwd,
		Config:        config,
		CWLDir:        cwlDir,
		InputContext:  make(map[string]interface{}),
		OutputContext: outputContext,
	}, nil
}

// Resolve takes the pre-decoded parameters for a CommandLineTool,
// ExpressionTool or Workflow and resolves everything to produce commands to
// run. The path to the decoded param file's dir must be provided to resolve
// relative paths.
//
// Also resolves any requirments, carrying out anything actionble, which may
// involve creating files according to an InitialWorkDirRequirement.
//
// If resolving a Workflow, you must be sure to Execute() all of the returned
// Commands in the appropriate order, and then call Output() on this to get the
// final output value.
func (r *Resolver) Resolve(name string, params Parameters, paramsDir string, ifc InputFileCallback) (Commands, error) {
	r.Name = name
	r.Parameters = params
	r.ParamsDir = paramsDir
	r.IFC = ifc

	if r.Workflow == nil {
		return nil, fmt.Errorf("nothing specified to do")
	}

	var cmds Commands
	var err error
	switch r.Workflow.Class {
	case classCommand, classExpression:
		cmds = append(cmds, NewCommand(name, []string{}, r))
	case classWorkflow:
		if r.Workflow.Steps == nil {
			return nil, fmt.Errorf("no steps specified in workflow")
		}

		scatter, multiple := r.Workflow.Requirements.DoScatterOrMultiple()

		// fill in missing params from workflow input defaults
		for _, in := range r.Workflow.Inputs {
			if _, exists := r.Parameters[in.ID]; !exists {
				if in.Default != nil {
					if thing, ok := in.Default.Self.(map[string]interface{}); ok && thing[fieldClass] == typeFile {
						// convert to a map[interface{}]interface{} so that
						// we later interpret this as a File
						file := make(map[interface{}]interface{})
						for k, v := range thing {
							file[k] = v
						}
						r.Parameters[in.ID] = file
					} else {
						r.Parameters[in.ID] = in.Default.Self
					}
				}
			}
		}

		stepOuts := make(map[string]map[string]bool)
		for _, step := range r.Workflow.Steps {
			stepParams := r.resolveStepParams(step.In, multiple, stepOuts)

			r.resolveStepOuts(step.ID, stepOuts, step.Out)

			sps, errs := r.handleScatter(scatter, step, stepParams)
			if errs != nil {
				return nil, errs
			}

			for _, sp := range sps {
				var subR *Resolver
				var subErr error
				if step.Run.Workflow == nil {
					if step.Run.Value == "" {
						return nil, fmt.Errorf("nothing to do for step %s", step.ID)
					}
					cwlPath := filepath.Join(r.CWLDir, step.Run.Value)

					cwlF, erro := os.Open(cwlPath)
					if erro != nil {
						return nil, erro
					}
					defer func() {
						err = cwlF.Close()
					}()

					root := NewCWL()
					errd := root.Decode(cwlF)
					if errd != nil {
						return nil, errd
					}

					subR, subErr = NewResolver(root, r.Config, filepath.Dir(cwlPath), r.OutputContext)
				} else {
					subR, subErr = NewResolver(step.Run.Workflow, r.Config, r.CWLDir, r.OutputContext)
				}
				if subErr != nil {
					return nil, subErr
				}

				subCmds, errr := subR.Resolve(name+"/"+step.ID, sp, paramsDir, ifc)
				if errr != nil {
					return nil, errr
				}
				for _, c := range subCmds {
					if c.Workflow.ID == "" {
						c.Workflow.ID = step.ID
					} else {
						c.Workflow.ID = step.ID + "/" + c.Workflow.ID
					}
				}
				cmds = append(cmds, subCmds...)
			}
		}
	}

	return cmds, err
}

// Output returns the final Workflow output. Only valid if called after
// Execute()ing all the Commands returned by Resolve().
func (r *Resolver) Output() interface{} {
	out := make(map[string]interface{})
	for _, o := range r.Workflow.Outputs {
		if len(o.Source) == 1 && o.Source[0] != "" {
			parts := strings.Split(o.Source[0], "/")
			if len(parts) == 2 {
				for key, val := range r.OutputContext {
					leaf := filepath.Base(key)
					if leaf == parts[0] {
						if oval, exists := val[parts[1]]; exists {
							out[o.ID] = oval
						}
					}
				}
			}
		} else {
			for key, val := range r.OutputContext {
				if key == r.Name {
					if oval, exists := val[o.ID]; exists {
						out[o.ID] = oval
					}
				}
			}
		}
	}
	return out
}

// resolveStepParams compares the given step inputs to the stored user
// parameters and other step outputs and sets values as appropriate in the
// returned Parameters.
func (r *Resolver) resolveStepParams(ins StepInputs, multiple bool, outs map[string]map[string]bool) Parameters {
	stepParams := *NewParameters()
	createdSlice := make(map[string]bool)
	for _, in := range ins {
		flatten := in.LinkMerge == mergeFlattened
		for _, source := range in.Source {
			if found, exists := r.Parameters[source]; exists {
				var vals []interface{}
				if arr, ok := found.([]interface{}); ok && flatten {
					vals = arr
				} else {
					vals = append(vals, found)
				}

				for _, val := range vals {
					if current, exists := stepParams[in.ID]; exists && multiple {
						if createdSlice[in.ID] {
							arr := stepParams[in.ID].([]interface{})
							arr = append(arr, val)
							stepParams[in.ID] = arr
						} else {
							stepParams[in.ID] = []interface{}{current, val}
							createdSlice[in.ID] = true
						}
					} else {
						stepParams[in.ID] = val
					}
				}
				continue
			}

			parts := strings.Split(source, "/")
			if len(parts) == 2 {
				if val, exists := outs[parts[0]]; exists {
					if val[parts[1]] {
						stepParams[in.ID] = source
						continue
					}
				}
			}
		}
	}
	return stepParams
}

// handleScatter creates a new Parameters for each scatter that is a copy of
// the input Parameters, altered for the particular input.
func (r *Resolver) handleScatter(scatter bool, step Step, stepParams Parameters) ([]Parameters, error) {
	var sps []Parameters
	if scatter {
		if len(step.Scatter) == 1 {
			fkey := step.Scatter[0]
			if param, exists := stepParams[fkey]; exists {
				if files, ok := param.([]interface{}); ok {
					for _, file := range files {
						theseParams := *NewParameters()
						for k, v := range stepParams {
							theseParams[k] = v
						}
						theseParams[fkey] = file
						sps = append(sps, theseParams)
					}
				} else {
					return nil, fmt.Errorf("request to scatter over non-array for %s", fkey)
				}
			}
		} else {
			scatterKeys := make(map[string]bool)
			for _, fkey := range step.Scatter {
				scatterKeys[fkey] = true
			}
			numScatter := len(step.Scatter)

			switch step.ScatterMethod {
			case scatterNestedCrossProduct, scatterFlatCrossProduct:
				scatterIndex := 0
				for i, iKey := range step.Scatter {
					if iParam, exists := stepParams[iKey]; exists {
						for j := i + 1; j < numScatter; j++ {
							jKey := step.Scatter[j]
							if jParam, exists := stepParams[jKey]; exists {
								iFiles, ok := iParam.([]interface{})
								if !ok {
									return nil, fmt.Errorf("request to scatter over non-array for %s", iKey)
								}
								jFiles, ok := jParam.([]interface{})
								if !ok {
									return nil, fmt.Errorf("request to scatter over non-array for %s", jKey)
								}

								for fi, iFile := range iFiles {
									for ji, jFile := range jFiles {
										// make some new params, copying non-scatter keys
										theseParams := *NewParameters()
										for k, v := range stepParams {
											if !scatterKeys[k] {
												theseParams[k] = v
											}
										}

										theseParams[iKey] = iFile
										theseParams[jKey] = jFile

										switch step.ScatterMethod {
										case scatterNestedCrossProduct:
											theseParams[scatterNestedInput] = [2]int{fi, ji}
										case scatterFlatCrossProduct:
											theseParams[scatterFlatInput] = scatterIndex
											scatterIndex++
										}

										sps = append(sps, theseParams)
									}
								}
							}
						}
					}
				}
			case scatterDotProduct:
				return nil, fmt.Errorf("request to scatter using dotproduct not yet implemented")
			}
		}
	} else {
		sps = []Parameters{stepParams}
	}

	return sps, nil
}

// resolveStepOuts fills out the given outs map with the give step outputs
func (r *Resolver) resolveStepOuts(name string, outs map[string]map[string]bool, stepOuts []StepOutput) {
	if _, exists := outs[name]; !exists {
		outs[name] = make(map[string]bool)
	}
	for _, sout := range stepOuts {
		outs[name][sout.ID] = true
	}
}
