package cwl

// Steps represents "steps" field in CWL.
type Steps []Step

// New constructs "Steps" from interface.
func (steps Steps) New(i interface{}) Steps {
	dest := Steps{}
	switch x := i.(type) {
	case []interface{}:
		for _, v := range x {
			s := Step{}.New(v)
			dest = append(dest, s)
		}
	case map[string]interface{}:
		for _, key := range sortKeys(x) {
			s := Step{}.New(x[key])
			s.ID = key
			dest = append(dest, s)
		}
	}
	return dest
}

// Step represents WorkflowStep.
// @see http://www.commonwl.org/v1.0/Workflow.html#WorkflowStep
type Step struct {
	ID            string
	In            StepInputs
	Out           []StepOutput
	Run           Run
	Requirements  []Requirement
	Hints         []Hint
	Scatter       []string
	ScatterMethod string
}

// Run `run` accept string | CommandLineTool | ExpressionTool | Workflow
type Run struct {
	Value    string
	Workflow *Root
}

// New constructs "Step" from interface.
func (s Step) New(i interface{}) Step {
	dest := Step{}
	switch x := i.(type) {
	case map[string]interface{}:
		for key, v := range x {
			switch key {
			case fieldID:
				dest.ID = v.(string)
			case fieldRun:
				switch x2 := v.(type) {
				case string:
					dest.Run.Value = x2
				case map[string]interface{}:
					dest.Run.Workflow = dest.Run.Workflow.AsStep(v)
				}
			case fieldIn:
				dest.In = StepInput{}.NewList(v)
			case fieldOut:
				dest.Out = StepOutput{}.NewList(v)
			case fieldRequirements:
				dest.Requirements = Requirements{}.New(v)
			case fieldHints:
				dest.Hints = Hints{}.New(v)
			case fieldScatter:
				dest.Scatter = StringArrayable(v)
			case fieldScatterMethod:
				dest.ScatterMethod = v.(string)
			}
		}
	}
	return dest
}
