package cwl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/otiai10/yaml2json"
)

// NewCWL ...
func NewCWL() *Root {
	root := new(Root)
	root.BaseCommands = BaseCommands{}
	root.Hints = Hints{}
	root.Inputs = Inputs{}
	return root
}

// Root ...
type Root struct {
	Version      string
	Class        string
	Hints        Hints
	Doc          string
	Graphs       Graphs
	BaseCommands BaseCommands
	Arguments    Arguments
	Namespaces   Namespaces
	Schemas      Schemas
	Stdin        string
	Stdout       string
	Stderr       string
	Inputs       Inputs `json:"inputs"`
	Outputs      Outputs
	Requirements Requirements
	Steps        Steps
	ID           string // ID only appears if this Root is a step in "steps"
	Expression   string // appears only if Class is "ExpressionTool"
}

// UnmarshalMap decode map[string]interface{} to *Root.
func (root *Root) UnmarshalMap(docs map[string]interface{}) error {
	for key, val := range docs {
		switch key {
		case fieldCWLVersion:
			root.Version = val.(string)
		case fieldClass:
			root.Class = val.(string)
		case fieldHints:
			root.Hints = root.Hints.New(val)
		case fieldDoc:
			root.Doc = val.(string)
		case fieldBaseCommand:
			root.BaseCommands = root.BaseCommands.New(val)
		case fieldArguments:
			root.Arguments = root.Arguments.New(val)
		case "$namespaces":
			root.Namespaces = root.Namespaces.New(val)
		case "$schemas":
			root.Schemas = root.Schemas.New(val)
		case "$graph":
			root.Graphs = root.Graphs.New(val)
		case fieldStdIn:
			root.Stdin = val.(string)
		case fieldStdOut:
			root.Stdout = val.(string)
		case fieldStdErr:
			root.Stderr = val.(string)
		case fieldInputs:
			root.Inputs = root.Inputs.New(val)
		case fieldOutputs:
			root.Outputs = root.Outputs.New(val)
		case fieldRequirements:
			root.Requirements = root.Requirements.New(val)
		case fieldSteps:
			root.Steps = root.Steps.New(val)
		case fieldID:
			root.ID = val.(string)
		case fieldExpression:
			root.Expression = val.(string)
		}
	}
	return nil
}

// UnmarshalJSON ...
func (root *Root) UnmarshalJSON(b []byte) error {
	docs := map[string]interface{}{}
	if err := json.Unmarshal(b, &docs); err != nil {
		return err
	}
	return root.UnmarshalMap(docs)
}

// Decode decodes specified file to this root
func (root *Root) Decode(r io.Reader) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("Parse error: %v", e)
		}
	}()
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	buf, err = yaml2json.Y2J(bytes.NewReader(buf))
	if err != nil {
		return err
	}
	if err = json.Unmarshal(buf, root); err != nil {
		return err
	}
	return nil
}

// Sort on a Decode()d Root sorts the various properties for consistent ordering
// of elements. This isn't necessary under normal circumstances, but may be
// useful for comparing (parts of) CWL files. Where possible, the original order
// of elements specified in the CWL file is kept; this sort only has an effect
// on parts of the CWL file that could have been specified with an array or a
// map and were specified with maps, in which case we sort by key when we return
// an array.
func (root *Root) Sort() {
	if root.Arguments != nil {
		sort.Sort(root.Arguments)
	}
	if root.Steps != nil {
		for _, s := range root.Steps {
			if &s.Run != nil && s.Run.Workflow != nil {
				s.Run.Workflow.Sort()
			}
		}
	}
	if root.Inputs != nil {
		sort.Sort(root.Inputs)
	}
	if root.Graphs != nil {
		for _, g := range root.Graphs {
			g.Sort()
		}
	}
}

// AsStep constructs Root as a step of "steps" from interface.
func (root *Root) AsStep(i interface{}) *Root {
	dest := new(Root)
	switch x := i.(type) {
	case string:
		dest.ID = x
	case map[string]interface{}:
		err := dest.UnmarshalMap(x)
		if err != nil {
			panic(fmt.Sprintf("Failed to parse step as CWL.Root: %v", err))
		}
	}
	return dest
}

// Resolve decodes the given CWL file, as well as the optional parameters file,
// to produce the concrete set of commands that must be run for this workflow.
// (Note that unlike Decode(), we need paths to the files to be able to resolve
// relative paths in the CWL and params file as relative to those files.)
// paramsPath can be an empty string if no parameters file is required.
//
// name should uniquely identify this combination of CWL and parameters.
//
// For CWL files that are $graphs, cwlPath can be suffixed with #[id] to select
// which node of the graph to resolve.
//
// To actually run the commands you can Execute() each of them in the right
// order.
//
// A Resolver is returned, on which you can call Output() to get the final
// output after Execute()ing all the commands.
func Resolve(name string, cwlPath, paramsPath string, config ResolveConfig, ifc InputFileCallback) (*Resolver, Commands, error) {
	var graphNode string
	parts := strings.Split(cwlPath, "#")
	if len(parts) == 2 {
		cwlPath = parts[0]
		graphNode = parts[1]
	}
	cwlPath, err := filepath.Abs(cwlPath)
	if err != nil {
		return nil, nil, err
	}

	cwlF, err := os.Open(cwlPath)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		err = cwlF.Close()
	}()

	root := NewCWL()
	err = root.Decode(cwlF)
	if err != nil {
		return nil, nil, err
	}

	r, err := NewResolver(root, config, filepath.Dir(cwlPath))
	if err != nil {
		return nil, nil, err
	}

	params := NewParameters()
	if paramsPath != "" {
		paramsPath, err = filepath.Abs(paramsPath)
		if err != nil {
			return nil, nil, err
		}

		paramsF, erro := os.Open(paramsPath)
		if erro != nil {
			return nil, nil, erro
		}
		defer func() {
			err = paramsF.Close()
		}()

		err = params.Decode(paramsF)
		if err != nil {
			return nil, nil, err
		}
	}

	cmds, err := r.Resolve(name, *params, filepath.Dir(paramsPath), ifc, graphNode, nil)

	return r, cmds, err
}
