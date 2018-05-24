package cwl

import (
	"fmt"
	"sort"
	"strings"
)

// Input represents "CommandInputParameter".
// @see http://www.commonwl.org/v1.0/CommandLineTool.html#CommandInputParameter
type Input struct {
	ID             string          `json:"id"`
	Label          string          `json:"label"`
	Doc            string          `json:"doc"`
	Format         string          `json:"format"`
	Binding        *Binding        `json:"inputBinding"`
	Default        *InputDefault   `json:"default"`
	Types          []Type          `json:"type"`
	SecondaryFiles []SecondaryFile `json:"secondary_files"`
	// Input.Provided is what provided by parameters.(json|yaml)
	Provided interface{} `json:"-"`
	// Requirement ..
	RequiredType *Type
	Requirements Requirements
}

// New constructs "Input" struct from interface{}.
func (input Input) New(i interface{}) Input {
	dest := Input{}
	switch x := i.(type) {
	case map[string]interface{}:
		for key, v := range x {
			switch key {
			case fieldID:
				dest.ID = v.(string)
			case fieldType:
				dest.Types = Type{}.NewList(v)
			case fieldLabel:
				dest.Label = v.(string)
			case fieldDoc:
				dest.Doc = v.(string)
			case fieldInputBinding:
				dest.Binding = Binding{}.New(v)
			case fieldDefault:
				dest.Default = InputDefault{}.New(v)
			case fieldFormat:
				dest.Format = v.(string)
			case fieldSecondaryFiles:
				dest.SecondaryFiles = SecondaryFile{}.NewList(v)
			}
		}
	case string:
		dest.Types = Type{}.NewList(x)
	case []interface{}:
		for _, v := range x {
			dest.Types = append(dest.Types, Type{}.New(v))
		}
	}
	return dest
}

// flatten
func (input Input) flatten(typ Type, binding *Binding) []string {
	flattened := []string{}
	switch typ.Type {
	case typeInt: // Array of Int
		tobejoined := []string{}
		for _, e := range input.Provided.([]interface{}) {
			tobejoined = append(tobejoined, fmt.Sprintf("%v", e))
		}
		flattened = append(flattened, strings.Join(tobejoined, input.Binding.Separator))
	case typeFile: // Array of Files
		switch arr := input.Provided.(type) {
		case []string:
			// TODO:
		case []interface{}:
			separated := []string{}
			for _, e := range arr {
				switch v := e.(type) {
				case map[interface{}]interface{}:
					if binding != nil && binding.Prefix != "" {
						separated = append(separated, binding.Prefix)
					}
					separated = append(separated, fmt.Sprintf("%v", v[fieldLocation]))
				default:
					// TODO:
				}
			}
			// In case it's Array of Files, unlike array of int,
			// it's NOT gonna be joined with .Binding.Separator.
			flattened = append(flattened, separated...)
		}
	default:
		if input.RequiredType != nil {
			flattened = append(flattened, input.flattenWithRequiredType()...)
		}
		//else {
		// TODO
		//}
	}
	return flattened
}

func (input Input) flattenWithRequiredType() []string {
	flattened := []string{}
	key, needed := input.Types[0].NeedRequirement()
	if !needed {
		return flattened
	}
	if input.RequiredType.Name != key {
		return flattened
	}
	switch provided := input.Provided.(type) {
	case []interface{}:
		for _, e := range provided {
			switch v := e.(type) {
			case map[interface{}]interface{}:
				for _, field := range input.RequiredType.Fields {
					if val, ok := v[field.Name]; ok {
						if field.Binding == nil {
							// Without thinking anything, just append it!!!
							flattened = append(flattened, fmt.Sprintf("%v", val))
						} else {
							if field.Binding.Prefix != "" {
								if field.Binding.Separate {
									flattened = append(flattened, field.Binding.Prefix, fmt.Sprintf("%v", val))
								} else {
									// TODO: Join if .Separator is given
									flattened = append(flattened, fmt.Sprintf("%s%v", field.Binding.Prefix, val))
								}
							} else {
								switch v2 := val.(type) {
								case []interface{}:
									for _, val2 := range v2 {
										switch v3 := val2.(type) {
										case []interface{}:
										case map[interface{}]interface{}:
											for _, types := range input.Requirements[0].SchemaDefRequirement.Types {
												val3array := []string{}
												var val3count int
												sort.Sort(types.Fields)
												for _, fields := range types.Fields {
													for key3, val3 := range v3 {
														if fields.Name == key3 {
															for _, val3type := range fields.Types {
																if val3type.Type == "" {
																} else {
																	switch val3type.Type {
																	case "enum":
																		for _, symbol := range val3type.Symbols {
																			if symbol == val3 {
																				val3array = append(val3array, fmt.Sprintf("%v", val3))
																				val3count = val3count + 1
																			}
																		}
																	case typeInt:
																		if fields.Binding.Prefix != "" {
																			val3array = append(val3array, fields.Binding.Prefix, fmt.Sprintf("%v", val3))
																			val3count = val3count + 1
																		} else {
																			val3array = append(val3array, fmt.Sprintf("%v", val3))
																			val3count = val3count + 1
																		}
																	}
																}
															}
														}
													}
												}
												if len(v3) == val3count {
													flattened = append(flattened, val3array...)
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return flattened
}

// Flatten ...
func (input Input) Flatten() []string {
	if input.Provided == nil {
		// In case "input.Default == nil" should be validated by usage layer.
		if input.Default != nil {
			return input.Default.Flatten(input.Binding)
		}
		return []string{}
	}
	flattened := []string{}
	if repr := input.Types[0]; len(input.Types) == 1 {
		switch repr.Type {
		case "array":
			flattened = append(flattened, input.flatten(repr.Items[0], repr.Binding)...)
		case "int":
			flattened = append(flattened, fmt.Sprintf("%v", input.Provided.(int)))
		case "File":
			switch provided := input.Provided.(type) {
			case map[interface{}]interface{}:
				// TODO: more strict type casting
				flattened = append(flattened, fmt.Sprintf("%v", provided["location"]))
			default:
			}
		default:
			flattened = append(flattened, fmt.Sprintf("%v", input.Provided))
		}
	}
	if input.Binding != nil && input.Binding.Prefix != "" {
		flattened = append([]string{input.Binding.Prefix}, flattened...)
	}

	return flattened
}

// Inputs represents "inputs" field in CWL.
type Inputs []Input

// New constructs new "Inputs" struct.
func (ins Inputs) New(i interface{}) Inputs {
	dest := Inputs{}
	switch x := i.(type) {
	case []interface{}:
		for _, v := range x {
			dest = append(dest, Input{}.New(v))
		}
	case map[string]interface{}:
		for key, v := range x {
			input := Input{}.New(v)
			input.ID = key
			dest = append(dest, input)
		}
	}
	return dest
}

// Len for sorting.
func (ins Inputs) Len() int {
	return len(ins)
}

// Less for sorting.
func (ins Inputs) Less(i, j int) bool {
	prev, next := ins[i].Binding, ins[j].Binding
	switch [2]bool{prev == nil, next == nil} {
	case [2]bool{true, true}:
		return true
	case [2]bool{false, true}:
		return prev.Position < 0
	case [2]bool{true, false}:
		return next.Position > 0
	default:
		return prev.Position <= next.Position
	}
}

// Swap for sorting.
func (ins Inputs) Swap(i, j int) {
	ins[i], ins[j] = ins[j], ins[i]
}
