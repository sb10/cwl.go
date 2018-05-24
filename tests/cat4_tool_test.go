package cwlgotest

import (
	"testing"

	cwl "github.com/sb10/cwl.go"
	. "github.com/otiai10/mint"
)

func TestDecode_cat4_tool(t *testing.T) {
	f := load("cat4-tool.cwl")
	root := cwl.NewCWL()
	Expect(t, root).TypeOf("*cwl.Root")
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)
	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Class).ToBe("CommandLineTool")
	Expect(t, len(root.Inputs)).ToBe(1)
	Expect(t, root.Inputs[0].ID).ToBe("file1")
	Expect(t, len(root.Outputs)).ToBe(1)
	Expect(t, root.Outputs[0].ID).ToBe("output_txt")
	Expect(t, root.Outputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Outputs[0].Binding.Glob).ToBe([]string{"output.txt"})
	Expect(t, len(root.BaseCommands)).ToBe(1)
	Expect(t, root.BaseCommands[0]).ToBe("cat")
	Expect(t, root.Stdout).ToBe("output.txt")
	Expect(t, root.Stdin).ToBe("$(inputs.file1.path)")
}
