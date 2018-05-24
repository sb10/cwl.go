package cwlgotest

import (
	"testing"

	. "github.com/otiai10/mint"
	cwl "github.com/sb10/cwl.go"
)

func TestDecode_schemadef_tool(t *testing.T) {
	f := load("schemadef-tool.cwl")
	root := cwl.NewCWL()
	err := root.Decode(f)
	Expect(t, err).ToBe(err)

	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Class).ToBe("CommandLineTool")
	Expect(t, root.Requirements[0].Import).ToBe("schemadef-type.yml")
	Expect(t, root.Requirements[1].Class).ToBe("InlineJavascriptRequirement")

	Expect(t, root.Inputs[0].ID).ToBe("hello")
	Expect(t, root.Inputs[0].Types[0].Type).ToBe("schemadef-type.yml#HelloType")
	Expect(t, root.Inputs[0].Binding.ValueFrom.Key()).ToBe(`self.a + "/" + self.b`)
	Expect(t, root.Outputs[0].ID).ToBe("output")
	Expect(t, root.Outputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Outputs[0].Binding.Glob[0]).ToBe("output.txt")
	Expect(t, root.Stdout).ToBe("output.txt")
	Expect(t, root.BaseCommands[0]).ToBe("echo")
}
