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

// Flatten ...
func (d *InputDefault) Flatten(binding *Binding) []string {
	flattened := []string{}
	switch v := d.Self.(type) {
	case map[string]interface{}:
		// TODO: more strict type casting ;(
		class, ok := v[fieldClass]
		if ok && class == typeFile {
			flattened = append(flattened, fmt.Sprintf("%v", v[fieldLocation]))
		}
	case string:
		flattened = append(flattened, d.Self.(string))
	}
	if binding != nil && binding.Prefix != "" {
		flattened = append([]string{binding.Prefix}, flattened...)
	}
	return flattened
}
