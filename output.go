package cwl

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
func (o Output) New(i interface{}) Output {
	dest := Output{}
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

// Outputs represents "outputs" field in "CWL".
type Outputs []Output

// New constructs "Outputs" struct.
func (o Outputs) New(i interface{}) Outputs {
	dest := Outputs{}
	switch x := i.(type) {
	case []interface{}:
		for _, v := range x {
			dest = append(dest, Output{}.New(v))
		}
	case map[string]interface{}:
		for key, v := range x {
			output := Output{}.New(v)
			output.ID = key
			dest = append(dest, output)
		}
	}
	return dest
}

// Len for sorting
func (o Outputs) Len() int {
	return len(o)
}

// Less for sorting
func (o Outputs) Less(i, j int) bool {
	prev, next := o[i].Binding, o[j].Binding
	switch [2]bool{prev == nil, next == nil} {
	case [2]bool{true, true}:
		return false
	case [2]bool{false, true}:
		return prev.Position < 0
	case [2]bool{true, false}:
		return next.Position > 0
	default:
		return prev.Position <= next.Position
	}
}

// Swap for sorting
func (o Outputs) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}
