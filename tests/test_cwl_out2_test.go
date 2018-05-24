package cwlgotest

import (
	"testing"

	. "github.com/otiai10/mint"
	cwl "github.com/sb10/cwl.go"
)

func TestDecode_test_cwl_out2(t *testing.T) {
	f := load("test-cwl-out2.cwl")
	root := cwl.NewCWL()
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)
	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Class).ToBe("CommandLineTool")
	Expect(t, root.Requirements[0].Class).ToBe("ShellCommandRequirement")
	Expect(t, root.Hints[0].Class).ToBe("DockerRequirement")
	Expect(t, root.Hints[0].DockerPull).ToBe("debian:wheezy")
	Expect(t, len(root.Inputs)).ToBe(0)
	// TODO check specification for this test ID and Type
	Expect(t, root.Outputs[0].ID).ToBe("foo")
	Expect(t, root.Outputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Arguments[0].Binding.ValueFrom.Key()).ToBe(`echo foo > foo && echo '{"foo": {"location": "file://$(runtime.outdir)/foo", "class": "File"} }' > cwl.output.json
`)
	Expect(t, root.Arguments[0].Binding.ShellQuote).ToBe(false)
}
