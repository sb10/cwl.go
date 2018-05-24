package cwlgotest

import (
	"sort"
	"testing"

	. "github.com/otiai10/mint"
	cwl "github.com/sb10/cwl.go"
)

func TestDecode_count_lines2_wf(t *testing.T) {
	f := load("count-lines2-wf.cwl")
	root := cwl.NewCWL()
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)

	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Class).ToBe("Workflow")

	Expect(t, root.Requirements[0].Class).ToBe("InlineJavascriptRequirement")
	Expect(t, root.Inputs[0].ID).ToBe("file1")
	Expect(t, root.Inputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Outputs[0].ID).ToBe("count_output")
	Expect(t, root.Outputs[0].Types[0].Type).ToBe("int")
	Expect(t, root.Outputs[0].Source).ToBe([]string{"step2/parseInt_output"})

	sort.Sort(root.Steps)
	Expect(t, root.Steps[0].In[0].ID).ToBe("wc_file1")
	Expect(t, root.Steps[0].In[0].Source[0]).ToBe("file1")
	Expect(t, root.Steps[0].Out[0].ID).ToBe("wc_output")
	Expect(t, root.Steps[0].Run).TypeOf("cwl.Run")
	Expect(t, root.Steps[0].Run.Workflow.ID).ToBe("wc")
	Expect(t, root.Steps[0].Run.Workflow.Class).ToBe("CommandLineTool")
	Expect(t, root.Steps[0].Run.Workflow.Inputs[0].ID).ToBe("wc_file1")
	Expect(t, root.Steps[0].Run.Workflow.Inputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Steps[0].Run.Workflow.Outputs[0].ID).ToBe("wc_output")
	Expect(t, root.Steps[0].Run.Workflow.Outputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Steps[0].Run.Workflow.Outputs[0].Binding.Glob).ToBe([]string{"output.txt"})
	Expect(t, root.Steps[0].Run.Workflow.Stdout).ToBe("output.txt")
	Expect(t, root.Steps[0].Run.Workflow.BaseCommands[0]).ToBe("wc")
	Expect(t, root.Steps[1].In[0].ID).ToBe("parseInt_file1")
	Expect(t, root.Steps[1].In[0].Source[0]).ToBe("step1/wc_output")
	Expect(t, root.Steps[1].Out[0].ID).ToBe("parseInt_output")
	Expect(t, root.Steps[1].Run.Workflow.Class).ToBe("ExpressionTool")
	Expect(t, root.Steps[1].Run.Workflow.Inputs[0].ID).ToBe("parseInt_file1")
	Expect(t, root.Steps[1].Run.Workflow.Inputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Steps[1].Run.Workflow.Inputs[0].Binding.LoadContents).ToBe(true)
	Expect(t, root.Steps[1].Run.Workflow.Outputs[0].ID).ToBe("parseInt_output")
	Expect(t, root.Steps[1].Run.Workflow.Outputs[0].Types[0].Type).ToBe("int")
	Expect(t, root.Steps[1].Run.Workflow.Expression).ToBe("${return {'parseInt_output': parseInt(inputs.parseInt_file1.contents)};}\n")
}
