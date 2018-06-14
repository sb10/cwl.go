package cwl

import (
	"sort"
	"strings"
)

// Argument represents an element of "arguments" of CWL
// @see http://www.commonwl.org/v1.0/CommandLineTool.html#CommandLineTool
type Argument struct {
	Value   string
	Binding *Binding
}

// New constructs an "Argument" struct from any interface.
func (arg Argument) New(i interface{}) Argument {
	dest := Argument{}
	switch x := i.(type) {
	case string:
		dest.Value = x
	case map[string]interface{}:
		dest.Binding = Binding{}.New(x)
	}
	return dest
}

// Flatten ...
func (arg Argument) Flatten() []string {
	flattened := []string{}
	if arg.Value != "" {
		flattened = append(flattened, arg.Value)
	}
	if arg.Binding != nil {
		if arg.Binding.Prefix != "" {
			flattened = append([]string{arg.Binding.Prefix}, flattened...)
		}
	}
	return flattened
}

// Arguments represents a list of "Argument"
type Arguments []Argument

// New constructs "Arguments" struct.
func (args Arguments) New(i interface{}) Arguments {
	dest := Arguments{}
	switch x := i.(type) {
	case []interface{}:
		for _, v := range x {
			dest = append(dest, Argument{}.New(v))
		}
	default:
		dest = append(dest, Argument{}.New(x))
	}
	return dest
}

// Len for sorting.
func (args Arguments) Len() int {
	return len(args)
}

// Less for sorting.
func (args Arguments) Less(i, j int) bool {
	prev, next := args[i].Binding, args[j].Binding
	switch [2]bool{prev == nil, next == nil} {
	case [2]bool{true, true}:
		return i < j
	case [2]bool{false, true}:
		return prev.Position < 0
	case [2]bool{true, false}:
		return next.Position > 0
	default:
		if prev.Position == next.Position {
			return i < j
		}
		return prev.Position < next.Position
	}
}

// Swap for sorting.
func (args Arguments) Swap(i, j int) {
	args[i], args[j] = args[j], args[i]
}

// Resolve goes through the arguments with "valueFrom" properties and returns
// concrete values from the given config. Also returns a bool, which if true
// means shell metacharacters should be quoted.
func (args Arguments) Resolve(config ResolveConfig) ([]string, bool) {
	var result []string
	var shellQuote bool
	sort.Sort(args)
	for i, arg := range args {
		if arg.Binding != nil && arg.Binding.ValueFrom != nil {
			// *** need to properly evaluate this if an expression?
			str := arg.Binding.ValueFrom.string
			if strings.HasPrefix(str, "$(") {
				args[i].Value = config.RuntimeValue(arg.Binding.ValueFrom.Key())
			} else {
				args[i].Value = str
			}

			if arg.Binding.ShellQuote {
				shellQuote = true
			}
		}
		result = append(result, args[i].Flatten()...)
	}
	return result, shellQuote
}
