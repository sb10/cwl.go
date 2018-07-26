package cwl

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/robertkrimen/otto"
)

// Output represents and combines "CommandOutputParameter" and "WorkflowOutputParameter"
// @see
// - http://www.commonwl.org/v1.0/CommandLineTool.html#CommandOutputParameter
// - http://www.commonwl.org/v1.0/Workflow.html#WorkflowOutputParameter
type Output struct {
	ID             string   `json:"id"`
	Label          string   `json:"label"`
	Doc            []string `json:"doc"`
	Format         string   `json:"format"`
	Binding        *Binding `json:"outputBinding"`
	Source         []string `json:"outputSource"`
	Types          []Type   `json:"type"`
	SecondaryFiles SecondaryFiles
}

// New constructs "Output" struct from interface.
func (o Output) New(i interface{}) *Output {
	dest := &Output{}
	switch x := i.(type) {
	case map[string]interface{}:
		for key, v := range x {
			switch key {
			case fieldID:
				dest.ID = v.(string)
			case fieldType:
				dest.Types = Type{}.NewList(v)
			case fieldOutputBinding:
				dest.Binding = Binding{}.New(v)
			case fieldOutputSource:
				dest.Source = StringArrayable(v)
			case fieldDoc:
				dest.Doc = StringArrayable(v)
			case fieldFormat:
				dest.Format = v.(string)
			case fieldSecondaryFiles:
				dest.SecondaryFiles = SecondaryFile{}.NewList(v)
			}
		}
	case string:
		dest.Types = Type{}.NewList(x)
	}
	return dest
}

// Resolve generates an output parameter based on the files produced by a
// CommandLineTool in the given output directory, specfied in the binding.
// stdoutPath is used if the type is 'stdout', to determine the path of the
// output file. Likewise for stderr. Expressions are evaluated with the given
// javascript vm. dir is the main configured output directory.
func (o *Output) Resolve(dir, stdoutPath, stderrPath string, vm *otto.Otto) (interface{}, error) {
	var result map[string]interface{}
	var results []map[string]interface{}
	var t string
	if repr := o.Types[0]; len(o.Types) == 1 {
		t = repr.Type
		switch repr.Type {
		case typeFile, typeInt, typeString:
			paths, err := globPaths(o.Binding, dir, vm)
			if err != nil {
				return nil, err
			}

			for _, path := range paths {
				thisResult, err := outputFileStats(dir, path, o.Binding.LoadContents)
				if err != nil {
					return nil, err
				}

				// add in details for any secondary files
				if len(o.SecondaryFiles) > 0 {
					files, err := o.SecondaryFiles.ToFiles(dir, path, vm, true)
					if err != nil {
						return nil, err
					}
					thisResult[fieldSecondaryFiles] = files
				}

				results = append(results, thisResult)
			}
		case fieldStdOut:
			if stdoutPath != "" {
				var err error
				result, err = outputFileStats(dir, filepath.Join(dir, stdoutPath), false)
				if err != nil {
					return nil, err
				}
			}
		case fieldStdErr:
			if stderrPath != "" {
				var err error
				result, err = outputFileStats(dir, filepath.Join(dir, stderrPath), false)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	if result == nil && len(results) > 0 {
		result = results[0] // *** not sure what to do with the other files
	}

	if o.Binding != nil && o.Binding.Eval != "" {
		err := vm.Set("self", results)
		if err != nil {
			return nil, err
		}
		str, fl, obj, err := evaluateExpression(o.Binding.Eval, vm)
		if err != nil {
			return nil, err
		}

		switch t {
		case typeInt:
			return int(fl), nil
		default:
			// *** don't know what to do if this was a File or something else,
			// just treat it as a str?
			if str != "" {
				return str, nil
			}
			return obj, nil
		}
	}

	// (to evaulate results we needed a map[string], but the final result we
	//  return must be map[interface{}])
	finalResult := make(map[interface{}]interface{})
	for key, val := range result {
		finalResult[key] = val
	}

	return finalResult, nil
}

func globPaths(binding *Binding, dir string, vm *otto.Otto) ([]string, error) {
	if binding != nil && binding.Glob != nil {
		var paths []string
		for _, glob := range binding.Glob {
			eglob, _, _, err := evaluateExpression(glob, vm)
			if err != nil {
				return nil, err
			}
			files, err := filepath.Glob(dir + "/" + eglob)
			if err != nil {
				return nil, err
			}
			paths = append(paths, files...)
		}
		return paths, nil
	}
	return nil, nil
}

func outputFileStats(dir, path string, loadContents bool) (map[string]interface{}, error) {
	// we need the file size and sha1 hash of the file contents
	size, sha, err := fileSizeAndSha(path)
	if err != nil {
		// we already know the file exists, so errors here
		// should not be ignored
		return nil, err
	}

	var content string
	if loadContents {
		content, err = getFileContents(path)
		if err != nil {
			return nil, err
		}
	}

	rel, err := filepath.Rel(dir, path)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		fieldClass:    typeFile,
		fieldLocation: rel,
		fieldSize:     size,
		fieldChecksum: sha,
	}

	if content != "" {
		result[fieldContents] = content
	}

	return result, err
}

func fileSizeAndSha(path string) (size int, sha string, err error) {
	info, err := os.Stat(path)
	if err != nil {
		// we already know the file exists, so errors here
		// should not be ignored
		return 0, "", err
	}

	f, err := os.Open(path)
	if err != nil {
		return 0, "", err
	}
	defer func() {
		err = f.Close()
	}()

	hash := sha1.New()
	_, err = io.Copy(hash, f)
	if err != nil {
		return 0, "", err
	}

	return int(info.Size()), fmt.Sprintf("sha1$%x", hash.Sum(nil)), err
}

// Outputs represents "outputs" field in "CWL".
type Outputs []*Output

// New constructs "Outputs" struct.
func (o Outputs) New(i interface{}) Outputs {
	dest := Outputs{}
	switch x := i.(type) {
	case []interface{}:
		for _, v := range x {
			dest = append(dest, Output{}.New(v))
		}
	case map[string]interface{}:
		for _, key := range sortKeys(x) {
			output := Output{}.New(x[key])
			output.ID = key
			dest = append(dest, output)
		}
	}
	return dest
}
