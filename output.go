package cwl

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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
	SecondaryFiles []SecondaryFile
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
func (o *Output) Resolve(dir string) (interface{}, error) {
	result := make(map[interface{}]interface{})
	if repr := o.Types[0]; len(o.Types) == 1 {
		switch repr.Type {
		case typeFile:
			if o.Binding != nil && o.Binding.Glob != nil {
				var paths []string
				for _, glob := range o.Binding.Glob {
					files, err := filepath.Glob(dir + "/" + glob)
					if err != nil {
						return nil, err
					}
					paths = append(paths, files...)
				}

				for _, path := range paths {
					// we need the file size
					info, err := os.Stat(path)
					if err != nil {
						continue
					}

					// and the sha1 hash of the file contents
					f, err := os.Open(path)
					if err != nil {
						return nil, err
					}
					defer f.Close()

					hash := sha1.New()
					_, err = io.Copy(hash, f)
					if err != nil {
						return nil, err
					}

					result = map[interface{}]interface{}{
						"class":    "File",
						"location": strings.TrimPrefix(path, dir+"/"),
						"size":     int(info.Size()),
						"checksum": fmt.Sprintf("sha1$%x", hash.Sum(nil)),
					}
				}
			}
		}
	}
	return result, nil
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
