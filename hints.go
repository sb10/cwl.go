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

// Merge merges in parent hints with this set, returning the new set.
// If a parent hint is already specified here, it is ignored.
func (h Hints) Merge(parentHints Hints) Hints {
	var hasEnvHint bool
	var envIndex int
	for i, hint := range h {
		switch hint.Class {
		case reqEnv:
			hasEnvHint = true
			envIndex = i
		}
	}

	merged := h[0:]
	for _, parentHint := range parentHints {
		switch parentHint.Class {
		case reqEnv:
			if !hasEnvHint {
				merged = append(merged, parentHint)
			} else {
				// merge env vars, ignoring parent values for keys set here
				alreadySet := make(map[string]bool)
				for _, ed := range merged[envIndex].Envs {
					alreadySet[ed.Name] = true
				}

				for _, ed := range parentHint.Envs {
					if !alreadySet[ed.Name] {
						merged[envIndex].Envs = append(merged[envIndex].Envs, ed)
					}
				}
			}
		}
	}
	return merged
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
