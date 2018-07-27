package cwl

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/robertkrimen/otto"
)

// Input represents "CommandInputParameter".
// @see http://www.commonwl.org/v1.0/CommandLineTool.html#CommandInputParameter
type Input struct {
	ID             string         `json:"id"`
	Label          string         `json:"label"`
	Doc            string         `json:"doc"`
	Format         string         `json:"format"`
	Binding        *Binding       `json:"inputBinding"`
	Default        *InputDefault  `json:"default"`
	Types          []Type         `json:"type"`
	SecondaryFiles SecondaryFiles `json:"secondary_files"`
	// Input.Provided is what provided by parameters.(json|yaml)
	Provided interface{} `json:"-"`
	// Requirement ..
	RequiredType *Type
	Requirements Requirements
}

// New constructs "Input" struct from interface{}.
func (input Input) New(i interface{}) *Input {
	dest := &Input{}
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
func (input *Input) flatten(typ Type, binding *Binding, inputContext map[string]interface{}, paramsDir string, ifc InputFileCallback, stepUnique string, vm *otto.Otto) ([]string, error) {
	flattened := []string{}
	switch typ.Type {
	case typeInt: // Array of Int
		inputContext[input.ID] = input.Provided
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
					path, err := resolvePath(fmt.Sprintf("%v", v[fieldLocation]), paramsDir, ifc, stepUnique, input.Binding, input.ID, inputContext, input.SecondaryFiles, vm)
					if err != nil {
						return nil, err
					}
					separated = append(separated, path)
				default:
					// TODO:
				}
			}
			// In case it's Array of Files, unlike array of int,
			// it's NOT gonna be joined with .Binding.Separator.
			flattened = append(flattened, separated...)
		}
	case typeBoolean:
		inputContext[input.ID] = input.Provided
	default:
		inputContext[input.ID] = input.Provided
		if input.RequiredType != nil {
			flattened = append(flattened, input.flattenWithRequiredType()...)
		}
		//else {
		// TODO
		//}
	}
	return flattened, nil
}

func (input *Input) flattenWithRequiredType() []string {
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

// InputFileCallback functions take a unique id and a path to an input file
// (which may be an absolute path on the local file system, or potentially an
// s3:// path), and return a path to where an executable should access that file
// from. This allows you to stage the file if desired.
type InputFileCallback func(id, path string) string

// DefaultInputFileCallback is an InputFileCallback that simply returns the
// path unaltered, for when you don't wish to do any staging, allowing direct
// access to input files at their original paths.
var DefaultInputFileCallback = func(id, path string) string {
	return path
}

// Flatten returns the flattened inputs in a []string, with any of those which
// were files having their paths made absolute relative to the given paramsDir
// (for files specified in the params file) or cwlDir (for files coming from a
// default specified in the CWL file), then altered by your callback (allowing
// you to stage input files).
//
// The directory options can be supplied as blank strings, in which case
// relative paths are not changed.
//
// Supplying the callback is optional, and if not supplied defaults to
// DefaultInputFileCallback.
//
// The provided inputContext will have an entry filled in with key of this
// input's ID, and value being some kind of resolved value.
func (input *Input) Flatten(inputContext map[string]interface{}, paramsDir, cwlDir string, stepUnique string, vm *otto.Otto, optionalIFC ...InputFileCallback) ([]string, error) {
	var ifc InputFileCallback
	if len(optionalIFC) == 1 && optionalIFC[0] != nil {
		ifc = optionalIFC[0]
	} else {
		ifc = DefaultInputFileCallback
	}

	var flattened []string
	if input.Provided == nil {
		// In case "input.Default == nil" should be validated by usage layer.
		if input.Default != nil {
			return input.Default.Flatten(input.Binding, input.ID, inputContext, cwlDir, ifc, stepUnique, vm)
		}
		return flattened, nil
	}

	addPrefix := true
	if repr := input.Types[0]; len(input.Types) == 1 {
		switch repr.Type {
		case typeArray:
			f, err := input.flatten(repr.Items[0], repr.Binding, inputContext, paramsDir, ifc, stepUnique, vm)
			if err != nil {
				return nil, err
			}
			flattened = append(flattened, f...)
		case typeInt:
			inputContext[input.ID] = input.Provided.(int)
			flattened = append(flattened, fmt.Sprintf("%v", input.Provided.(int)))
		case typeFile:
			if file, ok := input.Provided.(map[interface{}]interface{}); ok {
				// TODO: more strict type casting
				path, err := resolvePath(fmt.Sprintf("%v", file[fieldLocation]), paramsDir, ifc, stepUnique, input.Binding, input.ID, inputContext, input.SecondaryFiles, vm)
				if err != nil {
					return nil, err
				}
				flattened = append(flattened, path)
			} else {
				return nil, fmt.Errorf("expected a File, got %+v", input.Provided)
			}
		case typeFileSlice:
			if slice, ok := input.Provided.([]interface{}); ok {
				for _, maybeFile := range slice {
					if file, ok := maybeFile.(map[interface{}]interface{}); ok {
						path, err := resolvePath(fmt.Sprintf("%v", file[fieldLocation]), paramsDir, ifc, stepUnique, input.Binding, input.ID, inputContext, input.SecondaryFiles, vm)
						if err != nil {
							return nil, err
						}
						flattened = append(flattened, path)
					} else {
						return nil, fmt.Errorf("expected a File within the File[], got %+v", maybeFile)
					}
				}
			} else {
				return nil, fmt.Errorf("expected a File[], got %+v", input.Provided)
			}
		case typeBoolean:
			inputContext[input.ID] = input.Provided
			addPrefix = input.Provided.(bool)
		default:
			inputContext[input.ID] = input.Provided
			if input.RequiredType != nil {
				flattened = append(flattened, input.flattenWithRequiredType()...)
			} else {
				flattened = append(flattened, fmt.Sprintf("%v", input.Provided))
			}
		}
	}

	if addPrefix && input.Binding != nil && input.Binding.Prefix != "" {
		flattened = append([]string{input.Binding.Prefix}, flattened...)
	}

	return flattened, nil
}

// Resolve sets Provided if the supplied params have our ID. It also handles
// requirements.
func (input *Input) Resolve(params Parameters, reqs Requirements, customTypes map[string]*Type, vm *otto.Otto) error {
	if provided, ok := params[input.ID]; ok {
		if input.Binding != nil && input.Binding.ValueFrom != nil {
			// *** resolving valueFrom should also apply if not provided, but
			// we get a value from defaults...
			wasCustom := false
			if customTypes != nil && len(input.Types) == 1 {
				t := customTypes[input.Types[0].Type]
				if t != nil {
					self := make(map[string]interface{})
					for _, f := range t.Fields {
						// *** not actually paying attention to how the type is
						// defined...
						if val, exists := provided.(map[interface{}]interface{})[f.Name]; exists {
							self[f.Name] = val
						}
					}
					err := vm.Set("self", self)
					if err != nil {
						return err
					}
					wasCustom = true
				}
			}
			if !wasCustom {
				err := vm.Set("self", provided)
				if err != nil {
					return err
				}
			}

			str, _, _, err := evaluateExpression(input.Binding.ValueFrom.string, vm)
			if err != nil {
				return err
			}
			input.Provided = str
		} else {
			input.Provided = provided
		}
	}

	if input.Default == nil && input.Binding == nil && input.Provided == nil {
		return fmt.Errorf("input `%s` doesn't have a default, and not provided", input.ID)
	}

	if key, needed := input.Types[0].NeedRequirement(); needed {
		for _, req := range reqs {
			for _, requiredtype := range req.Types {
				if requiredtype.Name == key {
					input.RequiredType = &requiredtype
					input.Requirements = reqs
				}
			}
		}
	}

	return nil
}

// resolvePath converts relative paths to absolute ones, converted by the given
// callback. It also sets inputContext path and reads file content according to
// LoadContents binding.
func resolvePath(path, dir string, ifc InputFileCallback, stepUnique string, binding *Binding, id string, inputContext map[string]interface{}, secondaryFiles SecondaryFiles, vm *otto.Otto) (string, error) {
	if dir != "" && !filepath.IsAbs(path) {
		path = filepath.Join(dir, path)
	}

	inputContext[id] = map[string]string{}

	if binding != nil && binding.LoadContents {
		content, err := getFileContents(path)
		if err != nil {
			return "", err
		}
		inputContext[id].(map[string]string)[fieldContents] = content
	}

	stagedPath := ifc(stepUnique, path)
	inputContext[id].(map[string]string)["path"] = stagedPath

	// if there are any secondary files, pass those through the ifc as well to
	// get them staged if necessary
	if secondaryFiles != nil && vm != nil {
		sfs, err := secondaryFiles.ToFiles(filepath.Dir(path), path, vm, false)
		if err != nil {
			return "", err
		}
		parent := filepath.Dir(path)
		for _, f := range sfs {
			file := f.(map[interface{}]interface{})
			ifc(stepUnique, filepath.Join(parent, file[fieldLocation].(string)))
		}
	}

	return stagedPath, nil
}

func getFileContents(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() {
		err = f.Close()
	}()

	var header [65536]byte
	n, errr := io.ReadFull(f, header[:])
	if errr != nil && errr != io.ErrUnexpectedEOF {
		return "", errr
	}
	return string(header[0:n]), err
}

// Inputs represents "inputs" field in CWL.
type Inputs []*Input

// New constructs new "Inputs" struct.
func (ins Inputs) New(i interface{}) Inputs {
	var dest Inputs
	switch x := i.(type) {
	case []interface{}:
		for _, v := range x {
			dest = append(dest, Input{}.New(v))
		}
	case map[string]interface{}:
		for _, key := range sortKeys(x) {
			input := Input{}.New(x[key])
			input.ID = key
			dest = append(dest, input)
		}
	}
	return dest
}

// Resolve gets the concrete command line arguments (those that should go before
// anything else, and those that should go afterwards) for each input, by
// considering the provided params (and directories, for resolving relative
// paths). It also fills in the given context, for use with interpreting
// expressions.
func (ins Inputs) Resolve(reqs Requirements, params Parameters, paramsDir, CWLDir, stepUnique string, customTypes map[string]*Type, ifc InputFileCallback, inputContext map[string]interface{}, vm *otto.Otto) (result SortableArgs, err error) {
	defer func() {
		if i := recover(); i != nil {
			err = fmt.Errorf("failed to resolve required inputs against provided params: %v", i)
		}
	}()

	sort.Sort(ins)

	for i, in := range ins {
		err = in.Resolve(params, reqs, customTypes, vm)
		if err != nil {
			return result, err
		}

		theseIns, err := in.Flatten(inputContext, paramsDir, CWLDir, stepUnique, vm, ifc)
		if err != nil {
			return result, err
		}

		if in.Binding == nil {
			continue
		}

		result = append(result, &SortableArg{
			arg:      theseIns,
			position: in.Binding.Position,
			i:        i,
		})
	}
	return result, nil
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
		return ins[i].ID < ins[j].ID
	case [2]bool{false, true}:
		return prev.Position < 0
	case [2]bool{true, false}:
		return next.Position > 0
	default:
		if prev.Position == next.Position {
			return ins[i].ID < ins[j].ID
		}
		return prev.Position < next.Position
	}
}

// Swap for sorting.
func (ins Inputs) Swap(i, j int) {
	ins[i], ins[j] = ins[j], ins[i]
}
