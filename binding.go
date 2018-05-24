package cwl

// Binding represents and combines "CommandLineBinding" and "CommandOutputBinding"
// @see
// - http://www.commonwl.org/v1.0/CommandLineTool.html#CommandLineBinding
// - http://www.commonwl.org/v1.0/CommandLineTool.html#CommandOutputBinding
type Binding struct {
	// Common
	LoadContents bool
	// CommandLineBinding
	Position   int    `json:"position"`
	Prefix     string `json:"prefix"`
	Separate   bool   `json:"separate"`
	Separator  string `json:"separator"`
	ShellQuote bool   `json:"shellQuote"`
	ValueFrom  *Alias `json:"valueFrom"`
	// CommandOutputBinding
	Glob     []string `json:"glob"`
	Eval     string   `json:"outputEval"`
	Contents bool     `json:"loadContents"`
}

// New constructs new "Binding".
func (binding Binding) New(i interface{}) *Binding {
	dest := new(Binding)
	switch x := i.(type) {
	case map[string]interface{}:
		for key, v := range x {
			switch key {
			case fieldPosition:
				dest.Position = int(v.(float64))
			case fieldPrefix:
				dest.Prefix = v.(string)
			case fieldItemSeparator:
				dest.Separator = v.(string)
			case fieldLoadContents:
				dest.LoadContents = v.(bool)
			case fieldGlob:
				dest.Glob = StringArrayable(v)
			case fieldShellQuote:
				dest.ShellQuote = v.(bool)
			case fieldValueFrom:
				dest.ValueFrom = &Alias{v.(string)}
			case fieldOutputEval:
				dest.Eval = v.(string)
			}
		}
	}
	return dest
}
