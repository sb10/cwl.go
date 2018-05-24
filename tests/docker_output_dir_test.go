package cwlgotest

import (
	"testing"

	cwl "github.com/sb10/cwl.go"
	. "github.com/otiai10/mint"
)

func TestDecode_docker_output_dir(t *testing.T) {
	f := load("docker-output-dir.cwl")
	root := cwl.NewCWL()
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)
	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Class).ToBe("CommandLineTool")
	Expect(t, root.Requirements[0].Class).ToBe("DockerRequirement")
	Expect(t, root.Requirements[0].DockerPull).ToBe("debian:8")
	Expect(t, root.Requirements[0].DockerOutputDirectory).ToBe("/other")
	Expect(t, len(root.Inputs)).ToBe(0)
	Expect(t, len(root.Outputs)).ToBe(1)
	Expect(t, root.Outputs[0].ID).ToBe("thing")
	Expect(t, root.Outputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Outputs[0].Binding.Glob[0]).ToBe("thing")
	Expect(t, len(root.BaseCommands)).ToBe(2)
	Expect(t, root.BaseCommands[0]).ToBe("touch")
	Expect(t, root.BaseCommands[1]).ToBe("/other/thing")
}
