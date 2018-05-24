package cwl

import "strings"

// Type represents CWL Typeable objects.
// - http://www.commonwl.org/v1.0/CommandLineTool.html#CWLType
// - http://www.commonwl.org/v1.0/CommandLineTool.html#CommandInputRecordSchema
// - http://www.commonwl.org/v1.0/CommandLineTool.html#CommandInputEnumSchema
// - http://www.commonwl.org/v1.0/CommandLineTool.html#CommandInputArraySchema
type Type struct {
	Type    string
	Label   string
	Binding *Binding
	Fields  Fields   // from CommandInputRecordSchema
	Symbols []string // from CommandInputEnumSchema
	Items   []Type   // from CommandInputArraySchema
	Name    string
}

// NewList constructs a list of Type from any interface.
// It only handles []interface{}
func (t Type) NewList(i interface{}) []Type {
	dest := []Type{}
	switch x := i.(type) {
	case []interface{}:
		for _, s := range x {
			dest = append(dest, Type{}.New(s))
		}
	default:
		dest = append(dest, Type{}.New(x))
	}
	return dest
}

// New constructs single Type struct from any interface.
func (t Type) New(i interface{}) Type {
	dest := Type{}
	switch x := i.(type) {
	case string:
		dest.Type = x
	case map[string]interface{}:
		for key, v := range x {
			switch key {
			case fieldType:
				dest.Type = v.(string)
			case fieldItems:
				dest.Items = Type{}.NewList(v)
			case fieldInputBinding:
				dest.Binding = Binding{}.New(v)
			case fieldFields:
				dest.Fields = Fields{}.New(v)
			case fieldSymbols:
				dest.Symbols = StringArrayable(v)
			case fieldName:
				dest.Name = v.(string)
			}
		}
	}
	return dest
}

// NeedRequirement ...
func (t Type) NeedRequirement() (string, bool) {
	if strings.HasPrefix(t.Type, "#") {
		return strings.TrimPrefix(t.Type, "#"), true
	}
	for _, itemtype := range t.Items {
		if key, needed := itemtype.NeedRequirement(); needed {
			return key, needed
		}
	}
	return "", false
}
