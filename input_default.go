package cwl

import (
	"fmt"
	"reflect"
)

// InputDefault represents "default" field in an element of "inputs".
type InputDefault struct {
	Self interface{}
	Kind reflect.Kind
}

// New constructs new "InputDefault".
func (d InputDefault) New(i interface{}) *InputDefault {
	dest := &InputDefault{Self: i, Kind: reflect.TypeOf(i).Kind()}
	return dest
}

// Flatten returns the flattened inputs in a []string, with any of those which
// were files with relative paths being made absolute relative to the given
// cwl directory. If cwlDir is a blank string, the path is not altered. The
// optionalIFC allows you to define staging; if not supplied, defaults to
// DefaultInputFileCallback.
func (d *InputDefault) Flatten(binding *Binding, cwlDir string, optionalIFC ...InputFileCallback) []string {
	var ifc InputFileCallback
	if len(optionalIFC) == 1 && optionalIFC[0] != nil {
		ifc = optionalIFC[0]
	} else {
		ifc = DefaultInputFileCallback
	}

	var flattened []string
	switch v := d.Self.(type) {
	case map[string]interface{}:
		// TODO: more strict type casting ;(
		class, ok := v[fieldClass]
		if ok && class == typeFile {
			flattened = append(flattened, resolvePath(fmt.Sprintf("%v", v[fieldLocation]), cwlDir, ifc))
		}
	case string:
		flattened = append(flattened, d.Self.(string))
	}
	if binding != nil && binding.Prefix != "" {
		flattened = append([]string{binding.Prefix}, flattened...)
	}
	return flattened
}
