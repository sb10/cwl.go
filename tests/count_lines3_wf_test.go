package cwlgotest

import (
	"testing"

	cwl "github.com/sb10/cwl.go"
	. "github.com/otiai10/mint"
)

func TestDecode_count_lines3_wf(t *testing.T) {
	f := load("count-lines3-wf.cwl")
	root := cwl.NewCWL()
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)

	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Class).ToBe("Workflow")

	Expect(t, root.Inputs[0].ID).ToBe("file1")
	Expect(t, root.Inputs[0].Types[0].Type).ToBe("File[]")
	Expect(t, root.Outputs[0].ID).ToBe("count_output")
	Expect(t, root.Outputs[0].Types[0].Type).ToBe("int[]")
	Expect(t, root.Outputs[0].Source).ToBe([]string{"step1/output"})

	Expect(t, root.Requirements[0].Class).ToBe("ScatterFeatureRequirement")
	Expect(t, root.Steps[0].ID).ToBe("step1")
	Expect(t, root.Steps[0].Run.Value).ToBe("wc2-tool.cwl")
	Expect(t, root.Steps[0].Scatter[0]).ToBe("file1")
	Expect(t, root.Steps[0].In[0].ID).ToBe("file1")
	Expect(t, root.Steps[0].In[0].Source[0]).ToBe("file1")
	Expect(t, root.Steps[0].Out[0].ID).ToBe("output")
}
