package cwl

// StepInput represents WorkflowStepInput.
// @see http://www.commonwl.org/v1.0/Workflow.html#WorkflowStepInput
type StepInput struct {
	ID        string
	Source    []string
	LinkMerge string
	Default   *InputDefault
	ValueFrom string
}

// New constructs a StepInput struct from any interface.
func (s StepInput) New(i interface{}) StepInput {
	dest := StepInput{}
	switch x := i.(type) {
	case map[string]interface{}:
		for key, v := range x {
			if dest.ID == "" {
				dest.ID = key
			}

			if key == fieldID {
				dest.ID = v.(string)
			} else {
				switch e := v.(type) {
				case string:
					dest.Source = []string{e}
				case []interface{}:
					for _, s := range e {
						dest.Source = append(dest.Source, s.(string))
					}
				case map[string]interface{}:
					for key, v := range e {
						switch key {
						case fieldID:
							dest.ID = v.(string)
						case fieldSource:
							if list, ok := v.([]interface{}); ok {
								for _, s := range list {
									dest.Source = append(dest.Source, s.(string))
								}
							} else {
								dest.Source = append(dest.Source, v.(string))
							}
						case fieldLinkMerge:
							dest.LinkMerge = v.(string)
						case fieldDefault:
							dest.Default = InputDefault{}.New(v)
						case fieldValueFrom:
							dest.ValueFrom = v.(string)
						}
					}
				}
			}
		}
	}
	return dest
}

// StepInputs represents []StepInput
type StepInputs []StepInput

// NewList constructs a list of StepInput from interface.
func (s StepInput) NewList(i interface{}) StepInputs {
	dest := StepInputs{}
	switch x := i.(type) {
	case []interface{}:
		for _, v := range x {
			dest = append(dest, StepInput{}.New(v))
		}
	case map[string]interface{}:
		for _, key := range sortKeys(x) {
			item := make(map[string]interface{})
			item[key] = x[key]
			dest = append(dest, StepInput{}.New(item))
		}
	default:
		dest = append(dest, StepInput{}.New(x))
	}
	return dest
}
