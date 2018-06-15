package cwl

// Requirement represent an element of "requirements".
type Requirement struct {
	Class string
	InlineJavascriptRequirement
	SchemaDefRequirement
	DockerRequirement
	SoftwareRequirement
	InitialWorkDirRequirement
	EnvVarRequirement
	ShellCommandRequirement
	ResourceRequirement
	Import string
}

// New constructs "Requirement" struct from interface.
func (r Requirement) New(i interface{}) Requirement {
	dest := Requirement{}
	switch x := i.(type) {
	case map[string]interface{}:
		for key, v := range x {
			switch key {
			case fieldClass:
				dest.Class = v.(string)
			case fieldDockerPull:
				dest.DockerPull = v.(string)
			case fieldDockerOutputDirectory:
				dest.DockerOutputDirectory = v.(string)
			case fieldTypes:
				dest.Types = Type{}.NewList(v)
			case fieldExpressionLib:
				dest.ExpressionLib = JavascriptExpression{}.NewList(v)
			case fieldEnvDef:
				dest.EnvDef = EnvDef{}.NewList(v)
			case fieldListing:
				dest.Listing = Entry{}.NewList(v)
			case "$import":
				dest.Import = v.(string)
			}
		}
	}
	return dest
}

// InlineJavascriptRequirement is supposed to be embedded to Requirement.
// @see http://www.commonwl.org/v1.0/CommandLineTool.html#InlineJavascriptRequirement
type InlineJavascriptRequirement struct {
	ExpressionLib []JavascriptExpression
}

// JavascriptExpression represents an element of "expressionLib" of InlineJavascriptRequirement.
type JavascriptExpression struct {
	Kind  string // could be "" or "$include"
	Value string
}

// NewList constructs slice of JavascriptExpression from interface.
func (j JavascriptExpression) NewList(i interface{}) []JavascriptExpression {
	dest := []JavascriptExpression{}
	switch x := i.(type) {
	case []interface{}:
		for _, v := range x {
			dest = append(dest, JavascriptExpression{}.New(v))
		}
	}
	return dest
}

// New constructs JavascriptExpression from interface.
func (j JavascriptExpression) New(i interface{}) JavascriptExpression {
	dest := JavascriptExpression{}
	switch x := i.(type) {
	case string:
		dest.Kind = "$execute"
		dest.Value = x
	case map[string]interface{}:
		for key, v := range x {
			switch key {
			case "$include":
				dest.Kind = key
				dest.Value = v.(string)
			}
		}
	}
	return dest
}

// SchemaDefRequirement is supposed to be embedded to Requirement.
// @see http://www.commonwl.org/v1.0/CommandLineTool.html#SchemaDefRequirement
type SchemaDefRequirement struct {
	Types []Type
}

// DockerRequirement is supposed to be embedded to Requirement.
// @see http://www.commonwl.org/v1.0/CommandLineTool.html#DockerRequirement
type DockerRequirement struct {
	DockerPull            string
	DockerLoad            string
	DockerFile            string
	DockerImport          string
	DockerImageID         string
	DockerOutputDirectory string
}

// SoftwareRequirement is supposed to be embedded to Requirement.
// @see http://www.commonwl.org/v1.0/CommandLineTool.html#SoftwareRequirement
type SoftwareRequirement struct {
	Packages []SoftwarePackage
}

// SoftwarePackage represents an element of SoftwarePackage.Packages
type SoftwarePackage struct {
	Package  string
	Versions []string
	Specs    []string
}

// InitialWorkDirRequirement is supposed to be embedded to Requirement.
// @see http://www.commonwl.org/v1.0/CommandLineTool.html#InitialWorkDirRequirement
type InitialWorkDirRequirement struct {
	Listing []Entry
}

// EnvVarRequirement  is supposed to be embedded to Requirement.
// @see http://www.commonwl.org/v1.0/CommandLineTool.html#EnvVarRequirement
type EnvVarRequirement struct {
	EnvDef []EnvDef
}

// ShellCommandRequirement is supposed to be embedded to Requirement.
// @see http://www.commonwl.org/v1.0/CommandLineTool.html#ShellCommandRequirement
type ShellCommandRequirement struct {
}

// ResourceRequirement is supposed to be embedded to Requirement.
// @see http://www.commonwl.org/v1.0/CommandLineTool.html#ResourceRequirement
type ResourceRequirement struct {
	CoresMin int
	CoresMax int
}

// Requirements represents "requirements" field in CWL.
type Requirements []Requirement

// New constructs "Requirements" struct from interface.
func (r Requirements) New(i interface{}) Requirements {
	dest := Requirements{}
	switch x := i.(type) {
	case []interface{}:
		for _, r := range x {
			dest = append(dest, Requirement{}.New(r))
		}
	case map[string]interface{}:
		for _, key := range sortKeys(x) {
			r := Requirement{}.New(x[key])
			r.Class = key
			dest = append(dest, r)
		}
	}
	return dest
}

// DoScatter tells you if there is a ScatterFeatureRequirement.
func (r Requirements) DoScatter() bool {
	for _, req := range r {
		if req.Class == reqScatter {
			return true
		}
	}
	return false
}
