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
// ifc allows you to define staging.
func (d *InputDefault) Flatten(binding *Binding, id string, inputContext map[string]interface{}, cwlDir string, ifc InputFileCallback) []string {
	var flattened []string
	switch v := d.Self.(type) {
	case map[string]interface{}:
		// TODO: more strict type casting ;(
		class, ok := v[fieldClass]
		if ok && class == typeFile {
			path := resolvePath(fmt.Sprintf("%v", v[fieldLocation]), cwlDir, ifc)
			inputContext[id] = map[string]string{"path": path}
			flattened = append(flattened, path)
		}
	case string:
		inputContext[id] = d.Self.(string)
		flattened = append(flattened, d.Self.(string))
	default:
		inputContext[id] = d.Self
	}
	if binding != nil && binding.Prefix != "" {
		flattened = append([]string{binding.Prefix}, flattened...)
	}
	return flattened
}
