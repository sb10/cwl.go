package cwl

// Hints ...
type Hints []Hint

// New constructs "Hints" struct.
func (h Hints) New(i interface{}) Hints {
	dest := []Hint{}
	switch x := i.(type) {
	case []interface{}:
		for _, val := range x {
			switch e := val.(type) {
			case map[string]interface{}:
				hint := Hint{}.New(e)
				dest = append(dest, hint)
			}
		}
	case map[string]interface{}:
		for _, key := range sortKeys(x) {
			val := x[key]
			switch e := val.(type) {
			case map[string]interface{}:
				hint := Hint{}.New(e)
				hint.Class = key
				dest = append(dest, hint)
			}
		}
	}
	return dest
}

// Hint ...
type Hint struct {
	Class      string
	DockerPull string   // Only appears if class is "DockerRequirement"
	CoresMin   int      // Only appears if class is "ResourceRequirement"
	Envs       []EnvDef // Only appears if class is "EnvVarRequirement"
	FakeField  string   // Only appears if class is "ex:BlibberBlubberFakeRequirement"
	Import     string
}

// New constructs Hint from interface.
func (h Hint) New(i interface{}) Hint {
	dest := Hint{}
	switch x := i.(type) {
	case map[string]interface{}:
		for key, val := range x {
			switch key {
			case fieldClass:
				dest.Class = val.(string)
			case fieldDockerPull:
				dest.DockerPull = val.(string)
			case fieldCoresMin:
				dest.CoresMin = int(val.(float64))
			case "fakeField":
				dest.FakeField = val.(string)
			case fieldEnvDef:
				dest.Envs = EnvDef{}.NewList(val)
			case "$import":
				dest.Import = val.(string)
			}
		}
	}
	return dest
}
