package cwl

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	. "github.com/otiai10/mint"
)

const version = "1.0"

// Provides file object for testable official .cwl files.
func cwl(name string) *os.File {
	fpath := fmt.Sprintf("./cwl/v%[1]s/v%[1]s/%s", version, name)
	f, err := os.Open(fpath)
	if err != nil {
		panic(err)
	}
	return f
}

func TestDecode_count_lines3_wf(t *testing.T) {
	f := cwl("count-lines3-wf.cwl")
	root := NewCWL()
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
	Expect(t, root.Steps[0].Run.ID).ToBe("wc2-tool.cwl")
	Expect(t, root.Steps[0].Scatter).ToBe("file1")
	Expect(t, root.Steps[0].In[0].ID).ToBe("file1")
	Expect(t, root.Steps[0].In[0].Source[0]).ToBe("file1")
	Expect(t, root.Steps[0].Out[0].ID).ToBe("output")
}

func TestDecode_count_lines4_wf(t *testing.T) {
	f := cwl("count-lines4-wf.cwl")
	root := NewCWL()
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)

	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Class).ToBe("Workflow")

	Expect(t, root.Inputs[0].ID).ToBe("file1")
	Expect(t, root.Inputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Inputs[1].ID).ToBe("file2")
	Expect(t, root.Inputs[1].Types[0].Type).ToBe("File")
	Expect(t, root.Outputs[0].ID).ToBe("count_output")
	Expect(t, root.Outputs[0].Types[0].Type).ToBe("int[]")
	Expect(t, root.Outputs[0].Source).ToBe([]string{"step1/output"})

	Expect(t, root.Requirements[0].Class).ToBe("ScatterFeatureRequirement")
	Expect(t, root.Requirements[1].Class).ToBe("MultipleInputFeatureRequirement")
	Expect(t, root.Steps[0].ID).ToBe("step1")
	Expect(t, root.Steps[0].Run.ID).ToBe("wc2-tool.cwl")
	Expect(t, root.Steps[0].Scatter).ToBe("file1")
	Expect(t, root.Steps[0].In[0].ID).ToBe("file1")
	Expect(t, root.Steps[0].In[0].Source[0]).ToBe("file1")
	Expect(t, root.Steps[0].In[0].Source[1]).ToBe("file2")
	Expect(t, root.Steps[0].Out[0].ID).ToBe("output")
}

func TestDecode_count_lines5_wf(t *testing.T) {
	f := cwl("count-lines5-wf.cwl")
	root := NewCWL()
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)

	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Class).ToBe("Workflow")

	Expect(t, root.Inputs[0].ID).ToBe("file1")
	Expect(t, root.Inputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Inputs[0].Default.Kind).ToBe(reflect.Map)
	Expect(t, root.Outputs[0].ID).ToBe("count_output")
	Expect(t, root.Outputs[0].Types[0].Type).ToBe("int")
	Expect(t, root.Outputs[0].Source).ToBe([]string{"step1/output"})

	Expect(t, root.Steps[0].ID).ToBe("step1")
	Expect(t, root.Steps[0].Run.ID).ToBe("wc2-tool.cwl")
	Expect(t, root.Steps[0].In[0].ID).ToBe("file1")
	Expect(t, root.Steps[0].In[0].Source[0]).ToBe("file1")
	Expect(t, root.Steps[0].Out[0].ID).ToBe("output")
}

func TestDecode_count_lines6_wf(t *testing.T) {
	f := cwl("count-lines6-wf.cwl")
	root := NewCWL()
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)

	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Class).ToBe("Workflow")

	Expect(t, root.Inputs[0].ID).ToBe("file1")
	Expect(t, root.Inputs[0].Types[0].Type).ToBe("File[]")
	Expect(t, root.Inputs[1].ID).ToBe("file2")
	Expect(t, root.Inputs[1].Types[0].Type).ToBe("File[]")
	Expect(t, root.Outputs[0].ID).ToBe("count_output")
	Expect(t, root.Outputs[0].Types[0].Type).ToBe("int[]")
	Expect(t, root.Outputs[0].Source).ToBe([]string{"step1/output"})

	Expect(t, root.Requirements[0].Class).ToBe("ScatterFeatureRequirement")
	Expect(t, root.Requirements[1].Class).ToBe("MultipleInputFeatureRequirement")

	Expect(t, root.Steps[0].ID).ToBe("step1")
	Expect(t, root.Steps[0].Run.ID).ToBe("wc3-tool.cwl")
	Expect(t, root.Steps[0].Scatter).ToBe("file1")
	Expect(t, root.Steps[0].In[0].ID).ToBe("file1")
	Expect(t, root.Steps[0].In[0].Source[0]).ToBe("file1")
	Expect(t, root.Steps[0].In[0].Source[1]).ToBe("file2")
	Expect(t, root.Steps[0].In[0].LinkMerge).ToBe("merge_nested")
	Expect(t, root.Steps[0].Out[0].ID).ToBe("output")
}

func TestDecode_count_lines8_wf(t *testing.T) {
	f := cwl("count-lines8-wf.cwl")
	root := NewCWL()
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)

	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Class).ToBe("Workflow")

	Expect(t, root.Inputs[0].ID).ToBe("file1")
	Expect(t, root.Inputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Outputs[0].ID).ToBe("count_output")
	Expect(t, root.Outputs[0].Types[0].Type).ToBe("int")
	Expect(t, root.Outputs[0].Source).ToBe([]string{"step1/count_output"})

	Expect(t, root.Requirements[0].Class).ToBe("SubworkflowFeatureRequirement")

	Expect(t, root.Steps[0].ID).ToBe("step1")
	Expect(t, root.Steps[0].Run.ID).ToBe("count-lines1-wf.cwl")
	Expect(t, root.Steps[0].In[0].ID).ToBe("file1")
	Expect(t, root.Steps[0].In[0].Source[0]).ToBe("file1")
	Expect(t, root.Steps[0].Out[0].ID).ToBe("count_output")
}

func TestDecode_count_lines9_wf(t *testing.T) {
	f := cwl("count-lines9-wf.cwl")
	root := NewCWL()
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)

	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Class).ToBe("Workflow")

	Expect(t, len(root.Inputs)).ToBe(0)

	Expect(t, root.Outputs[0].ID).ToBe("count_output")
	Expect(t, root.Outputs[0].Types[0].Type).ToBe("int")
	Expect(t, root.Outputs[0].Source).ToBe([]string{"step2/output"})

	Expect(t, root.Steps[0].ID).ToBe("step1")
	Expect(t, root.Steps[0].Run.ID).ToBe("wc-tool.cwl")
	Expect(t, root.Steps[0].In[0].ID).ToBe("file1")
	Expect(t, root.Steps[0].In[0].Default.Kind).ToBe(reflect.Map)
	Expect(t, root.Steps[0].Out[0].ID).ToBe("output")

	Expect(t, root.Steps[1].ID).ToBe("step2")
	Expect(t, root.Steps[1].Run.ID).ToBe("parseInt-tool.cwl")
	Expect(t, root.Steps[1].In[0].ID).ToBe("file1")
	Expect(t, root.Steps[1].In[0].Source[0]).ToBe("step1/output")
	Expect(t, root.Steps[1].Out[0].ID).ToBe("output")
}

func TestDecode_count_lines10_wf(t *testing.T) {
	f := cwl("count-lines10-wf.cwl")
	root := NewCWL()
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)

	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Class).ToBe("Workflow")

	Expect(t, root.Inputs[0].ID).ToBe("file1")
	Expect(t, root.Inputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Outputs[0].ID).ToBe("count_output")
	Expect(t, root.Outputs[0].Types[0].Type).ToBe("int")
	Expect(t, root.Outputs[0].Source).ToBe([]string{"step1/count_output"})

	Expect(t, root.Requirements[0].Class).ToBe("SubworkflowFeatureRequirement")

	Expect(t, root.Steps[0].ID).ToBe("step1")
	Expect(t, root.Steps[0].In[0].ID).ToBe("file1")
	Expect(t, root.Steps[0].In[0].Source[0]).ToBe("file1")
	Expect(t, root.Steps[0].Out[0].ID).ToBe("count_output")

	Expect(t, root.Steps[0].Run.Class).ToBe("Workflow")
	Expect(t, root.Steps[0].Run.Inputs[0].ID).ToBe("file1")
	Expect(t, root.Steps[0].Run.Inputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Steps[0].Run.Outputs[0].ID).ToBe("count_output")
	Expect(t, root.Steps[0].Run.Outputs[0].Types[0].Type).ToBe("int")
	Expect(t, root.Steps[0].Run.Outputs[0].Source).ToBe([]string{"step2/output"})
	// Recursive steps
	Expect(t, root.Steps[0].Run.Steps[0].ID).ToBe("step1")
	Expect(t, root.Steps[0].Run.Steps[0].Run.ID).ToBe("wc-tool.cwl")
	Expect(t, root.Steps[0].Run.Steps[0].In[0].ID).ToBe("file1")
	Expect(t, root.Steps[0].Run.Steps[0].In[0].Source[0]).ToBe("file1")
	Expect(t, root.Steps[0].Run.Steps[0].Out[0].ID).ToBe("output")
	Expect(t, root.Steps[0].Run.Steps[1].ID).ToBe("step2")
	Expect(t, root.Steps[0].Run.Steps[1].Run.ID).ToBe("parseInt-tool.cwl")
	Expect(t, root.Steps[0].Run.Steps[1].In[0].ID).ToBe("file1")
	Expect(t, root.Steps[0].Run.Steps[1].In[0].Source[0]).ToBe("step1/output")
	Expect(t, root.Steps[0].Run.Steps[1].Out[0].ID).ToBe("output")
}

func TestDecode_count_lines11_wf(t *testing.T) {
	f := cwl("count-lines11-wf.cwl")
	root := NewCWL()
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)

	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Class).ToBe("Workflow")

	Expect(t, root.Inputs[0].ID).ToBe("file1")
	Expect(t, root.Inputs[0].Types[0].Type).ToBe("File?")
	Expect(t, root.Outputs[0].ID).ToBe("count_output")
	Expect(t, root.Outputs[0].Types[0].Type).ToBe("int")
	Expect(t, root.Outputs[0].Source).ToBe([]string{"step2/output"})

	Expect(t, root.Steps[0].ID).ToBe("step1")
	Expect(t, root.Steps[0].Run.ID).ToBe("wc-tool.cwl")
	Expect(t, root.Steps[0].In[0].ID).ToBe("file1")
	Expect(t, root.Steps[0].In[0].Default.Kind).ToBe(reflect.Map)
	Expect(t, root.Steps[0].Out[0].ID).ToBe("output")

	Expect(t, root.Steps[1].ID).ToBe("step2")
	Expect(t, root.Steps[1].Run.ID).ToBe("parseInt-tool.cwl")
	Expect(t, root.Steps[1].In[0].ID).ToBe("file1")
	Expect(t, root.Steps[1].In[0].Source[0]).ToBe("step1/output")
	Expect(t, root.Steps[1].Out[0].ID).ToBe("output")
}

func TestDecode_count_lines12_wf(t *testing.T) {
	f := cwl("count-lines12-wf.cwl")
	root := NewCWL()
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)

	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Class).ToBe("Workflow")

	Expect(t, root.Inputs[0].ID).ToBe("file1")
	Expect(t, root.Inputs[0].Types[0].Type).ToBe("array")
	Expect(t, root.Inputs[0].Types[0].Items[0].Type).ToBe("File")
	Expect(t, root.Inputs[1].ID).ToBe("file2")
	Expect(t, root.Inputs[1].Types[0].Type).ToBe("array")
	Expect(t, root.Inputs[1].Types[0].Items[0].Type).ToBe("File")
	Expect(t, root.Outputs[0].ID).ToBe("count_output")
	Expect(t, root.Outputs[0].Types[0].Type).ToBe("int")
	Expect(t, root.Outputs[0].Source).ToBe([]string{"step1/output"})

	Expect(t, root.Requirements[0].Class).ToBe("MultipleInputFeatureRequirement")

	Expect(t, root.Steps[0].ID).ToBe("step1")
	Expect(t, root.Steps[0].Run.ID).ToBe("wc3-tool.cwl")
	Expect(t, root.Steps[0].In[0].ID).ToBe("file1")
	Expect(t, root.Steps[0].In[0].Source[0]).ToBe("file1")
	Expect(t, root.Steps[0].In[0].Source[1]).ToBe("file2")
	Expect(t, root.Steps[0].In[0].LinkMerge).ToBe("merge_flattened")
	Expect(t, root.Steps[0].Out[0].ID).ToBe("output")
}

func TestDecode_dir(t *testing.T) {
	f := cwl("dir.cwl")
	root := NewCWL()
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)

	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Class).ToBe("CommandLineTool")

	Expect(t, root.Requirements[0].Class).ToBe("ShellCommandRequirement")

	Expect(t, root.Inputs[0].Types[0].Type).ToBe("Directory")
	Expect(t, root.Outputs[0].ID).ToBe("outlist")
	Expect(t, root.Outputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Outputs[0].Binding.Glob).ToBe([]string{"output.txt"})
	Expect(t, root.Arguments[0].Value).ToBe("cd")
	Expect(t, root.Arguments[1].Value).ToBe("$(inputs.indir.path)")
	Expect(t, root.Arguments[2].Binding.ShellQuote).ToBe(false)
	Expect(t, root.Arguments[2].Binding.ValueFrom).ToBe("&&")
	Expect(t, root.Arguments[3].Value).ToBe("find")
	Expect(t, root.Arguments[4].Value).ToBe(".")
	Expect(t, root.Arguments[5].Binding.ShellQuote).ToBe(false)
	Expect(t, root.Arguments[5].Binding.ValueFrom).ToBe("|")
	Expect(t, root.Arguments[6].Value).ToBe("sort")
	Expect(t, root.Stdout).ToBe("output.txt")
}

func TestDecode_dir2(t *testing.T) {
	f := cwl("dir2.cwl")
	root := NewCWL()
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)

	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Class).ToBe("CommandLineTool")

	Expect(t, root.Hints[0].Class).ToBe("DockerRequirement")
	Expect(t, root.Hints[0].DockerPull).ToBe("debian:8")
	Expect(t, root.Hints[1].Class).ToBe("ShellCommandRequirement")
	Expect(t, root.Inputs[0].ID).ToBe("indir")
	Expect(t, root.Inputs[0].Types[0].Type).ToBe("Directory")
	Expect(t, root.Outputs[0].ID).ToBe("outlist")
	Expect(t, root.Outputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Outputs[0].Binding.Glob).ToBe([]string{"output.txt"})
	Expect(t, root.Arguments[0].Value).ToBe("cd")
	Expect(t, root.Arguments[1].Value).ToBe("$(inputs.indir.path)")
	Expect(t, root.Arguments[2].Binding.ShellQuote).ToBe(false)
	Expect(t, root.Arguments[2].Binding.ValueFrom).ToBe("&&")
	Expect(t, root.Arguments[3].Value).ToBe("find")
	Expect(t, root.Arguments[4].Value).ToBe(".")
	Expect(t, root.Arguments[5].Binding.ShellQuote).ToBe(false)
	Expect(t, root.Arguments[5].Binding.ValueFrom).ToBe("|")
	Expect(t, root.Arguments[6].Value).ToBe("sort")
	Expect(t, root.Stdout).ToBe("output.txt")
}

func TestDecode_dir3(t *testing.T) {
	f := cwl("dir3.cwl")
	root := NewCWL()
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)

	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Class).ToBe("CommandLineTool")
	Expect(t, root.BaseCommands[0]).ToBe("tar")
	Expect(t, root.BaseCommands[1]).ToBe("xvf")
	Expect(t, root.Inputs[0].ID).ToBe("inf")
	Expect(t, root.Inputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Inputs[0].Binding.Position).ToBe(1)
	Expect(t, root.Outputs[0].ID).ToBe("outdir")
	Expect(t, root.Outputs[0].Types[0].Type).ToBe("Directory")
	Expect(t, root.Outputs[0].Binding.Glob).ToBe([]string{"."})
}

func TestDecode_dir4(t *testing.T) {
	f := cwl("dir4.cwl")
	root := NewCWL()
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)

	Expect(t, root.Requirements[0].Class).ToBe("ShellCommandRequirement")
	Expect(t, root.Inputs[0].ID).ToBe("inf")
	Expect(t, root.Inputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Outputs[0].ID).ToBe("outlist")
	Expect(t, root.Outputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Outputs[0].Binding.Glob).ToBe([]string{"output.txt"})
	Expect(t, root.Arguments[0].Value).ToBe("cd")
	Expect(t, root.Arguments[1].Value).ToBe("$(inputs.inf.dirname)/xtestdir")
	Expect(t, root.Arguments[2].Binding.ShellQuote).ToBe(false)
	Expect(t, root.Arguments[2].Binding.ValueFrom).ToBe("&&")
	Expect(t, root.Arguments[3].Value).ToBe("find")
	Expect(t, root.Arguments[4].Value).ToBe(".")
	Expect(t, root.Arguments[5].Binding.ShellQuote).ToBe(false)
	Expect(t, root.Arguments[5].Binding.ValueFrom).ToBe("|")
	Expect(t, root.Arguments[6].Value).ToBe("sort")
	Expect(t, root.Stdout).ToBe("output.txt")
}

func TestDecode_dir5(t *testing.T) {
	f := cwl("dir5.cwl")
	root := NewCWL()
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)

	Expect(t, root.Requirements[0].Class).ToBe("ShellCommandRequirement")
	Expect(t, root.Requirements[1].Class).ToBe("InitialWorkDirRequirement")
	Expect(t, root.Requirements[1].Listing[0].Location).ToBe("$(inputs.indir.listing)")
	Expect(t, root.Inputs[0].ID).ToBe("indir")
	Expect(t, root.Inputs[0].Types[0].Type).ToBe("Directory")
	Expect(t, root.Outputs[0].ID).ToBe("outlist")
	Expect(t, root.Outputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Outputs[0].Binding.Glob).ToBe([]string{"output.txt"})
	Expect(t, root.Arguments[0].Value).ToBe("find")
	Expect(t, root.Arguments[1].Value).ToBe("-L")
	Expect(t, root.Arguments[2].Value).ToBe(".")
	Expect(t, root.Arguments[3].Binding.ShellQuote).ToBe(false)
	Expect(t, root.Arguments[3].Binding.ValueFrom).ToBe("|")
	Expect(t, root.Arguments[4].Value).ToBe("sort")
	Expect(t, root.Stdout).ToBe("output.txt")
}

func TestDecode_dir6(t *testing.T) {
	f := cwl("dir6.cwl")
	root := NewCWL()
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)

	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Class).ToBe("CommandLineTool")

	Expect(t, root.Requirements[0].Class).ToBe("ShellCommandRequirement")
	Expect(t, root.Inputs[0].ID).ToBe("indir")
	Expect(t, root.Inputs[0].Types[0].Type).ToBe("Directory")
	Expect(t, root.Inputs[0].Binding.Position).ToBe(-1)
	Expect(t, root.Inputs[0].Binding.Prefix).ToBe("cd")
	Expect(t, root.Outputs[0].ID).ToBe("outlist")
	Expect(t, root.Outputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Outputs[0].Binding.Glob).ToBe([]string{"output.txt"})
	Expect(t, root.Arguments[0].Binding.ShellQuote).ToBe(false)
	Expect(t, root.Arguments[0].Binding.ValueFrom).ToBe("&&")
	Expect(t, root.Arguments[1].Value).ToBe("find")
	Expect(t, root.Arguments[2].Value).ToBe(".")
	Expect(t, root.Arguments[3].Binding.ShellQuote).ToBe(false)
	Expect(t, root.Arguments[3].Binding.ValueFrom).ToBe("|")
	Expect(t, root.Arguments[4].Value).ToBe("sort")
	Expect(t, root.Stdout).ToBe("output.txt")
}

func TestDecode_dir7(t *testing.T) {
	f := cwl("dir7.cwl")
	root := NewCWL()
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)

	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Class).ToBe("ExpressionTool")

	Expect(t, root.Requirements[0].Class).ToBe("InlineJavascriptRequirement")

	Expect(t, root.Inputs[0].ID).ToBe("files")
	Expect(t, root.Inputs[0].Types[0].Type).ToBe("File[]")
	Expect(t, root.Outputs[0].ID).ToBe("dir")
	Expect(t, root.Outputs[0].Types[0].Type).ToBe("Directory")
	Expect(t, root.Expression).ToBe(`${
return {"dir": {"class": "Directory", "basename": "a_directory", "listing": inputs.files}};
}`)
}

func TestDecode_cat3_nodocker(t *testing.T) {
	f := cwl("cat3-nodocker.cwl")
	root := NewCWL()
	Expect(t, root).TypeOf("*cwl.Root")
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)
	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Doc).ToBe("Print the contents of a file to stdout using 'cat'.")
	Expect(t, root.Class).ToBe("CommandLineTool")
	Expect(t, len(root.BaseCommands)).ToBe(1)
	Expect(t, root.BaseCommands[0]).ToBe("cat")
	Expect(t, root.Stdout).ToBe("output.txt")
	Expect(t, len(root.Inputs)).ToBe(1)
	Expect(t, root.Inputs[0].ID).ToBe("file1")
	Expect(t, root.Inputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Inputs[0].Label).ToBe("Input File")
	Expect(t, root.Inputs[0].Doc).ToBe("The file that will be copied using 'cat'")
	Expect(t, root.Inputs[0].Binding.Position).ToBe(1)
}

func TestDecode_cat3_tool_mediumcut(t *testing.T) {
	f := cwl("cat3-tool-mediumcut.cwl")
	root := NewCWL()
	Expect(t, root).TypeOf("*cwl.Root")
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)
	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Doc).ToBe("Print the contents of a file to stdout using 'cat' running in a docker container.")
	Expect(t, root.Class).ToBe("CommandLineTool")
	Expect(t, len(root.BaseCommands)).ToBe(1)
	Expect(t, root.BaseCommands[0]).ToBe("cat")
	Expect(t, root.Stdout).ToBe("cat-out")
	Expect(t, root.Hints).TypeOf("cwl.Hints")
	Expect(t, root.Hints[0].Class).ToBe("DockerRequirement")
	Expect(t, root.Hints[0].DockerPull).ToBe("debian:wheezy")
	Expect(t, len(root.Inputs)).ToBe(1)
	Expect(t, root.Inputs[0].ID).ToBe("file1")
	Expect(t, root.Inputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Inputs[0].Label).ToBe("Input File")
	Expect(t, root.Inputs[0].Doc).ToBe("The file that will be copied using 'cat'")
	Expect(t, root.Inputs[0].Binding.Position).ToBe(1)
}

func TestDecode_cat3_tool_shortcut(t *testing.T) {
	f := cwl("cat3-tool-shortcut.cwl")
	root := NewCWL()
	Expect(t, root).TypeOf("*cwl.Root")
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)
	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Doc).ToBe("Print the contents of a file to stdout using 'cat' running in a docker container.")
	Expect(t, root.Class).ToBe("CommandLineTool")
	Expect(t, len(root.BaseCommands)).ToBe(1)
	Expect(t, root.BaseCommands[0]).ToBe("cat")
	Expect(t, root.Hints).TypeOf("cwl.Hints")
	Expect(t, root.Hints[0].Class).ToBe("DockerRequirement")
	Expect(t, root.Hints[0].DockerPull).ToBe("debian:wheezy")
	Expect(t, len(root.Inputs)).ToBe(1)
	Expect(t, root.Inputs[0].ID).ToBe("file1")
	Expect(t, root.Inputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Inputs[0].Label).ToBe("Input File")
	Expect(t, root.Inputs[0].Doc).ToBe("The file that will be copied using 'cat'")
	Expect(t, root.Inputs[0].Binding.Position).ToBe(1)
}

func TestDecode_cat3_tool(t *testing.T) {
	f := cwl("cat3-tool.cwl")
	root := NewCWL()
	Expect(t, root).TypeOf("*cwl.Root")
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)
	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Doc).ToBe("Print the contents of a file to stdout using 'cat' running in a docker container.")
	Expect(t, root.Class).ToBe("CommandLineTool")
	Expect(t, len(root.BaseCommands)).ToBe(1)
	Expect(t, root.BaseCommands[0]).ToBe("cat")
	Expect(t, root.Stdout).ToBe("output.txt")
	Expect(t, root.Hints).TypeOf("cwl.Hints")
	Expect(t, root.Hints[0].Class).ToBe("DockerRequirement")
	Expect(t, root.Hints[0].DockerPull).ToBe("debian:wheezy")
	Expect(t, len(root.Inputs)).ToBe(1)
	Expect(t, root.Inputs[0].ID).ToBe("file1")
	Expect(t, root.Inputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Inputs[0].Label).ToBe("Input File")
	Expect(t, root.Inputs[0].Doc).ToBe("The file that will be copied using 'cat'")
	Expect(t, root.Inputs[0].Binding.Position).ToBe(1)
}

func TestDecode_env_tool1(t *testing.T) {
	f := cwl("env-tool1.cwl")
	root := NewCWL()
	Expect(t, root).TypeOf("*cwl.Root")
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)
	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, len(root.BaseCommands)).ToBe(3)
	Expect(t, root.BaseCommands[0]).ToBe("/bin/bash")
	Expect(t, root.BaseCommands[1]).ToBe("-c")
	Expect(t, root.BaseCommands[2]).ToBe("echo $TEST_ENV")
	Expect(t, len(root.Inputs)).ToBe(1)
	// TODO ignore "in: string'
	Expect(t, len(root.Outputs)).ToBe(1)
	Expect(t, root.Outputs[0].ID).ToBe("out")
	Expect(t, root.Outputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Outputs[0].Binding.Glob).ToBe([]string{"out"})
}

func TestDecode_default_path(t *testing.T) {
	f := cwl("default_path.cwl")
	root := NewCWL()
	Expect(t, root).TypeOf("*cwl.Root")
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)
	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Class).ToBe("CommandLineTool")
	Expect(t, len(root.Inputs)).ToBe(1)
	Expect(t, root.Inputs[0].ID).ToBe("file1")
	Expect(t, root.Inputs[0].Types[0].Type).ToBe("File")
	// TODO support default: section
	// TODO support outputs: []
	Expect(t, len(root.Arguments)).ToBe(2)
	Expect(t, root.Arguments[0].Value).ToBe("cat")
	Expect(t, root.Arguments[1].Value).ToBe("$(inputs.file1.path)")
}

func TestDecode_cat4_tool(t *testing.T) {
	f := cwl("cat4-tool.cwl")
	root := NewCWL()
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

func TestDecode_cat5_tool(t *testing.T) {
	f := cwl("cat5-tool.cwl")
	root := NewCWL()
	Expect(t, root).TypeOf("*cwl.Root")
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)
	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Class).ToBe("CommandLineTool")
	Expect(t, root.Doc).ToBe("Print the contents of a file to stdout using 'cat' running in a docker container.")
	Expect(t, len(root.Hints)).ToBe(2)
	Expect(t, root.Hints).TypeOf("cwl.Hints")
	Expect(t, root.Hints[0].Class).ToBe("DockerRequirement")
	Expect(t, root.Hints[0].DockerPull).ToBe("debian:wheezy")
	Expect(t, root.Hints[1].Class).ToBe("ex:BlibberBlubberFakeRequirement")
	Expect(t, root.Hints[1].FakeField).ToBe("fraggleFroogle")
	Expect(t, len(root.Inputs)).ToBe(1)
	Expect(t, root.Inputs[0].ID).ToBe("file1")
	Expect(t, root.Inputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Inputs[0].Label).ToBe("Input File")
	Expect(t, root.Inputs[0].Doc).ToBe("The file that will be copied using 'cat'")
	Expect(t, root.Inputs[0].Binding.Position).ToBe(1)
	Expect(t, len(root.Outputs)).ToBe(1)
	Expect(t, root.Outputs[0].ID).ToBe("output_file")
	Expect(t, root.Outputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Outputs[0].Binding.Glob).ToBe([]string{"output.txt"})
	Expect(t, len(root.BaseCommands)).ToBe(1)
	Expect(t, root.BaseCommands[0]).ToBe("cat")
	Expect(t, root.Stdout).ToBe("output.txt")
	// $namespaces
	Expect(t, len(root.Namespaces)).ToBe(1)
	Expect(t, root.Namespaces[0]["ex"]).ToBe("http://example.com/")
}
func TestDecode_env_tool2(t *testing.T) {
	f := cwl("env-tool2.cwl")
	root := NewCWL()
	Expect(t, root).TypeOf("*cwl.Root")
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)
	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Class).ToBe("CommandLineTool")
	Expect(t, len(root.Hints)).ToBe(1)
	Expect(t, root.Hints).TypeOf("cwl.Hints")
	Expect(t, root.Hints[0].Class).ToBe("EnvVarRequirement")
	Expect(t, root.Hints[0].Envs[0].Name).ToBe("TEST_ENV")
	Expect(t, root.Hints[0].Envs[0].Value).ToBe("$(inputs.in)")
	Expect(t, len(root.Inputs)).ToBe(1)
	// TODO in: string
	Expect(t, len(root.Outputs)).ToBe(1)
	Expect(t, len(root.BaseCommands)).ToBe(3)
	Expect(t, root.BaseCommands[0]).ToBe("/bin/bash")
	Expect(t, root.BaseCommands[1]).ToBe("-c")
	Expect(t, root.BaseCommands[2]).ToBe("echo $TEST_ENV")
}

func TestDecode_env_wf1(t *testing.T) {
	f := cwl("env-wf1.cwl")
	root := NewCWL()
	Expect(t, root).TypeOf("*cwl.Root")
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)
	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Class).ToBe("Workflow")
	Expect(t, len(root.Inputs)).ToBe(1)
	// TODO in: string
	Expect(t, len(root.Outputs)).ToBe(1)
	Expect(t, root.Outputs[0].ID).ToBe("out")
	Expect(t, root.Outputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Outputs[0].Source).ToBe([]string{"step1/out"})
	Expect(t, len(root.Requirements)).ToBe(1)
	Expect(t, root.Requirements[0].Class).ToBe("EnvVarRequirement")
	Expect(t, root.Requirements[0].EnvDef[0].Name).ToBe("TEST_ENV")
	Expect(t, root.Requirements[0].EnvDef[0].Value).ToBe(`override`)
}

func TestDecode_envvar(t *testing.T) {
	f := cwl("envvar.cwl")
	root := NewCWL()
	Expect(t, root).TypeOf("*cwl.Root")
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)
	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Class).ToBe("CommandLineTool")
	Expect(t, len(root.Inputs)).ToBe(0)
	Expect(t, len(root.Outputs)).ToBe(0)
	Expect(t, len(root.Requirements)).ToBe(1)
	Expect(t, root.Requirements[0].Class).ToBe("ShellCommandRequirement")
	Expect(t, len(root.Arguments)).ToBe(12)
	Expect(t, root.Arguments[0].Value).ToBe("echo")
	Expect(t, root.Arguments[1].Binding.ValueFrom).ToBe("\"HOME=$HOME\"")
	Expect(t, root.Arguments[1].Binding.ShellQuote).ToBe(false)
	Expect(t, root.Arguments[2].Binding.ValueFrom).ToBe("\"TMPDIR=$TMPDIR\"")
	Expect(t, root.Arguments[2].Binding.ShellQuote).ToBe(false)
	Expect(t, root.Arguments[3].Binding.ValueFrom).ToBe("&&")
	Expect(t, root.Arguments[3].Binding.ShellQuote).ToBe(false)
	Expect(t, root.Arguments[4].Value).ToBe("test")
	Expect(t, root.Arguments[5].Binding.ValueFrom).ToBe("\"$HOME\"")
	Expect(t, root.Arguments[5].Binding.ShellQuote).ToBe(false)
	Expect(t, root.Arguments[6].Value).ToBe("=")
	Expect(t, root.Arguments[7].Value).ToBe("$(runtime.outdir)")
	Expect(t, root.Arguments[8].Value).ToBe("-a")
	Expect(t, root.Arguments[9].Binding.ValueFrom).ToBe("\"$TMPDIR\"")
	Expect(t, root.Arguments[9].Binding.ShellQuote).ToBe(false)
	Expect(t, root.Arguments[10].Value).ToBe("=")
	Expect(t, root.Arguments[11].Value).ToBe("$(runtime.tmpdir)")
}

func TestDecode_envvar2(t *testing.T) {
	f := cwl("envvar2.cwl")
	root := NewCWL()
	Expect(t, root).TypeOf("*cwl.Root")
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)
	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Class).ToBe("CommandLineTool")
	Expect(t, len(root.Inputs)).ToBe(0)
	Expect(t, len(root.Outputs)).ToBe(0)
	Expect(t, len(root.Requirements)).ToBe(1)
	Expect(t, root.Requirements[0].Class).ToBe("ShellCommandRequirement")
	Expect(t, root.Hints[0].Class).ToBe("DockerRequirement")
	Expect(t, root.Hints[0].DockerPull).ToBe("debian:8")
	Expect(t, len(root.Arguments)).ToBe(12)
	Expect(t, root.Arguments[0].Value).ToBe("echo")
	Expect(t, root.Arguments[1].Binding.ValueFrom).ToBe("\"HOME=$HOME\"")
	Expect(t, root.Arguments[1].Binding.ShellQuote).ToBe(false)
	Expect(t, root.Arguments[2].Binding.ValueFrom).ToBe("\"TMPDIR=$TMPDIR\"")
	Expect(t, root.Arguments[2].Binding.ShellQuote).ToBe(false)
	Expect(t, root.Arguments[3].Binding.ValueFrom).ToBe("&&")
	Expect(t, root.Arguments[3].Binding.ShellQuote).ToBe(false)
	Expect(t, root.Arguments[4].Value).ToBe("test")
	Expect(t, root.Arguments[5].Binding.ValueFrom).ToBe("\"$HOME\"")
	Expect(t, root.Arguments[5].Binding.ShellQuote).ToBe(false)
	Expect(t, root.Arguments[6].Value).ToBe("=")
	Expect(t, root.Arguments[7].Value).ToBe("$(runtime.outdir)")
	Expect(t, root.Arguments[8].Value).ToBe("-a")
	Expect(t, root.Arguments[9].Binding.ValueFrom).ToBe("\"$TMPDIR\"")
	Expect(t, root.Arguments[9].Binding.ShellQuote).ToBe(false)
	Expect(t, root.Arguments[10].Value).ToBe("=")
	Expect(t, root.Arguments[11].Value).ToBe("$(runtime.tmpdir)")
}

func TestDecode_formattest(t *testing.T) {
	f := cwl("formattest.cwl")
	root := NewCWL()
	Expect(t, root).TypeOf("*cwl.Root")
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)
	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Class).ToBe("CommandLineTool")
	Expect(t, len(root.Inputs)).ToBe(1)
	Expect(t, root.Inputs[0].ID).ToBe("input")
	Expect(t, root.Inputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Inputs[0].Format).ToBe("edam:format_2330")
	Expect(t, root.Inputs[0].Binding.Position).ToBe(0)
	Expect(t, len(root.Outputs)).ToBe(1)
	Expect(t, root.Outputs[0].ID).ToBe("output")
	Expect(t, root.Outputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Outputs[0].Binding.Glob).ToBe([]string{"output.txt"})
	Expect(t, root.Outputs[0].Format).ToBe("edam:format_2330")
	Expect(t, len(root.BaseCommands)).ToBe(1)
	Expect(t, root.BaseCommands[0]).ToBe("rev")
	Expect(t, root.Stdout).ToBe("output.txt")
}

func TestDecode_formattest2(t *testing.T) {
	f := cwl("formattest2.cwl")
	root := NewCWL()
	Expect(t, root).TypeOf("*cwl.Root")
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)
	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Class).ToBe("CommandLineTool")
	Expect(t, len(root.Inputs)).ToBe(1)
	Expect(t, root.Inputs[0].ID).ToBe("input")
	Expect(t, root.Inputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Inputs[0].Format).ToBe("edam:format_2330")
	Expect(t, root.Inputs[0].Binding.Position).ToBe(0)
	Expect(t, len(root.Outputs)).ToBe(1)
	Expect(t, root.Outputs[0].ID).ToBe("output")
	Expect(t, root.Outputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Outputs[0].Binding.Glob).ToBe([]string{"output.txt"})
	Expect(t, root.Outputs[0].Format).ToBe("$(inputs.input.format)")
	Expect(t, len(root.BaseCommands)).ToBe(1)
	Expect(t, root.BaseCommands[0]).ToBe("rev")
	Expect(t, root.Stdout).ToBe("output.txt")
}

func TestDecode_formattest3(t *testing.T) {
	f := cwl("formattest3.cwl")
	root := NewCWL()
	Expect(t, root).TypeOf("*cwl.Root")
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)
	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Class).ToBe("CommandLineTool")
	// $namespaces
	Expect(t, len(root.Namespaces)).ToBe(2)
	Expect(t, root.Namespaces[0]["edam"]).ToBe("http://edamontology.org/")
	Expect(t, root.Namespaces[1]["gx"]).ToBe("http://galaxyproject.org/formats/")
	// $namespaces
	Expect(t, len(root.Schemas)).ToBe(2)
	Expect(t, root.Schemas[0]).ToBe("EDAM.owl")
	Expect(t, root.Schemas[1]).ToBe("gx_edam.ttl")
	Expect(t, root.Doc).ToBe("Reverse each line using the `rev` command")
	Expect(t, root.Hints[0].Class).ToBe("DockerRequirement")
	Expect(t, root.Hints[0].DockerPull).ToBe("debian:wheezy")
	Expect(t, len(root.Inputs)).ToBe(1)
	Expect(t, root.Inputs[0].ID).ToBe("input")
	Expect(t, root.Inputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Inputs[0].Format).ToBe("gx:fasta")
	Expect(t, root.Inputs[0].Binding.Position).ToBe(0)
	Expect(t, len(root.Outputs)).ToBe(1)
	Expect(t, root.Outputs[0].ID).ToBe("output")
	Expect(t, root.Outputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Outputs[0].Binding.Glob).ToBe([]string{"output.txt"})
	Expect(t, root.Outputs[0].Format).ToBe("$(inputs.input.format)")
	Expect(t, len(root.BaseCommands)).ToBe(1)
	Expect(t, root.BaseCommands[0]).ToBe("rev")
	Expect(t, root.Stdout).ToBe("output.txt")
}

func TestDecode_glob_expr_list(t *testing.T) {
	f := cwl("glob-expr-list.cwl")
	root := NewCWL()
	Expect(t, root).TypeOf("*cwl.Root")
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)
	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Class).ToBe("CommandLineTool")
	Expect(t, len(root.Inputs)).ToBe(1)
	Expect(t, root.Inputs[0].ID).ToBe("ids")
	Expect(t, root.Inputs[0].Types[0].Type).ToBe("string[]")
	Expect(t, root.Inputs[0].Binding.Position).ToBe(1)
	Expect(t, len(root.Outputs)).ToBe(1)
	Expect(t, root.Outputs[0].ID).ToBe("files")
	Expect(t, root.Outputs[0].Types[0].Type).ToBe("File[]")
	Expect(t, root.Outputs[0].Binding.Glob).ToBe([]string{"$(inputs.ids)"})
	Expect(t, len(root.BaseCommands)).ToBe(1)
	Expect(t, root.BaseCommands[0]).ToBe("touch")
}

func TestDecode_imported_hint(t *testing.T) {
	f := cwl("imported-hint.cwl")
	root := NewCWL()
	Expect(t, root).TypeOf("*cwl.Root")
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)
	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Class).ToBe("CommandLineTool")
	Expect(t, len(root.Inputs)).ToBe(0)
	Expect(t, len(root.Outputs)).ToBe(1)
	Expect(t, root.Outputs[0].ID).ToBe("out")
	// TODO test out: stdout 's stdout
	Expect(t, root.Hints).TypeOf("cwl.Hints")
	Expect(t, root.Hints[0].Import).ToBe("envvar.yml")
	Expect(t, len(root.BaseCommands)).ToBe(3)
	Expect(t, root.BaseCommands[0]).ToBe("/bin/bash")
	Expect(t, root.BaseCommands[1]).ToBe("-c")
	Expect(t, root.BaseCommands[2]).ToBe("echo $TEST_ENV")
	Expect(t, root.Stdout).ToBe("out")
}

func TestDecode_initialwork_path(t *testing.T) {
	f := cwl("initialwork-path.cwl")
	root := NewCWL()
	Expect(t, root).TypeOf("*cwl.Root")
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)
	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Class).ToBe("CommandLineTool")
	Expect(t, len(root.Inputs)).ToBe(1)
	Expect(t, root.Inputs[0].ID).ToBe("file1")
	Expect(t, root.Inputs[0].Types[0].Type).ToBe("File")
	Expect(t, len(root.Outputs)).ToBe(0)
	Expect(t, root.Requirements[0].Class).ToBe("InitialWorkDirRequirement")
	Expect(t, root.Requirements[0].Listing[0].EntryName).ToBe("bob.txt")
	Expect(t, root.Requirements[0].Listing[0].Entry).ToBe(`$(inputs.file1)`)
	Expect(t, root.Requirements[1].Class).ToBe("ShellCommandRequirement")
	Expect(t, root.Arguments[0].Binding.ShellQuote).ToBe(false)
	Expect(t, root.Arguments[0].Binding.ValueFrom).ToBe(`test "$(inputs.file1.path)" = "$(runtime.outdir)/bob.txt"
`)
	// TODO write basecommand
}

func TestDecode_initialworkdirrequirement_docker_out(t *testing.T) {
	f := cwl("initialworkdirrequirement-docker-out.cwl")
	root := NewCWL()
	Expect(t, root).TypeOf("*cwl.Root")
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)
	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Class).ToBe("CommandLineTool")
	Expect(t, len(root.Inputs)).ToBe(1)
	Expect(t, root.Inputs[0].ID).ToBe("INPUT")
	Expect(t, root.Inputs[0].Types[0].Type).ToBe("File")
	Expect(t, len(root.Outputs)).ToBe(1)
	Expect(t, root.Outputs[0].ID).ToBe("OUTPUT")
	Expect(t, root.Outputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Outputs[0].Binding.Glob).ToBe([]string{"$(inputs.INPUT.basename)"})
	Expect(t, root.Outputs[0].SecondaryFiles[0].Entry).ToBe(".fai")
	// TODO outputs
	Expect(t, len(root.Requirements)).ToBe(2)
	Expect(t, root.Requirements[0].Class).ToBe("DockerRequirement")
	Expect(t, root.Requirements[0].DockerPull).ToBe("debian:8")
	Expect(t, root.Requirements[1].Class).ToBe("InitialWorkDirRequirement")
	Expect(t, root.Requirements[1].Listing[0].Location).ToBe("$(inputs.INPUT)")
	Expect(t, root.Arguments[0].Binding.ValueFrom).ToBe("$(inputs.INPUT.basename).fai")
	// TODO test against "position" but currently just put 0 is failed
	Expect(t, len(root.BaseCommands)).ToBe(1)
	Expect(t, root.BaseCommands[0]).ToBe("touch")
}

func TestDecode_inline_js(t *testing.T) {
	f := cwl("inline-js.cwl")
	root := NewCWL()
	Expect(t, root).TypeOf("*cwl.Root")
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)
	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Class).ToBe("CommandLineTool")
	// TODO test BaseCommand because this file has two baseCommand fields
	//fmt.Println(root.BaseCommands)
	//Expect(t, len(root.BaseCommands)).ToBe(0)
	//Expect(t, root.BaseCommands[0]).ToBe("touch")
	Expect(t, len(root.Requirements)).ToBe(1)
	Expect(t, root.Requirements[0].Class).ToBe("InlineJavascriptRequirement")
	Expect(t, len(root.Requirements)).ToBe(1)
	Expect(t, root.Hints[0].Class).ToBe("DockerRequirement")
	Expect(t, root.Hints[0].DockerPull).ToBe("python:2-slim")
	Expect(t, len(root.Inputs)).ToBe(1)
	Expect(t, root.Inputs[0].ID).ToBe("args.py")
	Expect(t, root.Inputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Inputs[0].Default.Kind).ToBe(reflect.Map)
	Expect(t, root.Inputs[0].Binding.Position).ToBe(-1)
	Expect(t, len(root.Outputs)).ToBe(1)
	Expect(t, root.Outputs[0].ID).ToBe("args")
	Expect(t, root.Outputs[0].Types[0].Type).ToBe("array")
	Expect(t, root.Outputs[0].Types[0].Items[0].Type).ToBe("string")
	Expect(t, len(root.Arguments)).ToBe(3)
	Expect(t, root.Arguments[0].Binding.Prefix).ToBe("-A")
	Expect(t, root.Arguments[0].Binding.ValueFrom).ToBe("$(1+1)")
	Expect(t, root.Arguments[1].Binding.Prefix).ToBe("-B")
	Expect(t, root.Arguments[1].Binding.ValueFrom).ToBe(`$("/foo/bar/baz".split('/').slice(-1)[0])`)
	Expect(t, root.Arguments[2].Binding.Prefix).ToBe("-C")
	Expect(t, root.Arguments[2].Binding.ValueFrom).ToBe(`${
  var r = [];
  for (var i = 10; i >= 1; i--) {
    r.push(i);
  }
  return r;
}
`)
}

func TestDecode_js_expr_req_wf(t *testing.T) {
	f := cwl("js-expr-req-wf.cwl")
	root := NewCWL()
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)
	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, len(root.Graphs)).ToBe(2)
	// 0
	Expect(t, root.Graphs[0].ID).ToBe("tool")
	Expect(t, root.Graphs[0].Class).ToBe("CommandLineTool")
	Expect(t, root.Graphs[0].Requirements[0].Class).ToBe("InlineJavascriptRequirement")
	Expect(t, root.Graphs[0].Requirements[0].ExpressionLib[0].Value).ToBe("function foo() { return 2; }")
	Expect(t, len(root.Graphs[0].Inputs)).ToBe(0)
	Expect(t, root.Graphs[0].Arguments[0].Value).ToBe("echo")
	Expect(t, root.Graphs[0].Stdout).ToBe("whatever.txt")
	Expect(t, len(root.Graphs[0].Outputs)).ToBe(1)
	Expect(t, root.Graphs[0].Outputs[0].ID).ToBe("out")
	Expect(t, root.Graphs[0].Outputs[0].Types[0].Type).ToBe("stdout")
	// 1
	Expect(t, root.Graphs[1].ID).ToBe("wf")
	Expect(t, root.Graphs[1].Class).ToBe("Workflow")
	Expect(t, root.Graphs[1].Requirements[0].Class).ToBe("InlineJavascriptRequirement")
	Expect(t, root.Graphs[1].Requirements[0].ExpressionLib[0].Value).ToBe("function bar() { return 1; }")
	Expect(t, len(root.Graphs[1].Inputs)).ToBe(0)
	Expect(t, root.Graphs[1].Outputs[0].ID).ToBe("out")
	Expect(t, root.Graphs[1].Outputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Graphs[1].Outputs[0].Source[0]).ToBe("tool/out")
	Expect(t, root.Graphs[1].Steps[0].ID).ToBe("tool")
	Expect(t, root.Graphs[1].Steps[0].Run.ID).ToBe("#tool")
	Expect(t, len(root.Graphs[1].Steps[0].In)).ToBe(1)
	// TODO check empty In
	Expect(t, len(root.Graphs[1].Steps[0].Out)).ToBe(1)
	Expect(t, root.Graphs[1].Steps[0].Out[0].ID).ToBe("out")
}

func TestDecode_nameroot(t *testing.T) {
	f := cwl("nameroot.cwl")
	root := NewCWL()
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)
	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Class).ToBe("CommandLineTool")
	Expect(t, len(root.Inputs)).ToBe(1)
	Expect(t, root.Inputs[0].ID).ToBe("file1")
	Expect(t, root.Inputs[0].Types[0].Type).ToBe("File")
	Expect(t, len(root.Outputs)).ToBe(1)
	Expect(t, root.Outputs[0].ID).ToBe("b")
	Expect(t, len(root.BaseCommands)).ToBe(0)
	Expect(t, len(root.Arguments)).ToBe(4)
	Expect(t, root.Arguments[0].Value).ToBe("echo")
	Expect(t, root.Arguments[1].Value).ToBe("$(inputs.file1.basename)")
	Expect(t, root.Arguments[2].Value).ToBe("$(inputs.file1.nameroot)")
	Expect(t, root.Arguments[3].Value).ToBe("$(inputs.file1.nameext)")
}

func TestDecode_nested_array(t *testing.T) {
	f := cwl("nested-array.cwl")
	root := NewCWL()
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)
	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Class).ToBe("CommandLineTool")
	Expect(t, root.BaseCommands[0]).ToBe("echo")
	Expect(t, len(root.Inputs)).ToBe(1)
	Expect(t, root.Inputs[0].ID).ToBe("letters")
	Expect(t, root.Inputs[0].Types[0].Type).ToBe("array")
	Expect(t, root.Inputs[0].Types[0].Items[0].Type).ToBe("array")
	Expect(t, root.Inputs[0].Types[0].Items[0].Items[0].Type).ToBe("string")
	Expect(t, root.Inputs[0].Binding.Position).ToBe(1)
	Expect(t, root.Stdout).ToBe("echo.txt")
	Expect(t, root.Outputs[0].ID).ToBe("echo")
	Expect(t, root.Outputs[0].Types[0].Type).ToBe("stdout")
}

func TestDecode_null_defined(t *testing.T) {
	f := cwl("null-defined.cwl")
	root := NewCWL()
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)
	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Class).ToBe("CommandLineTool")
	Expect(t, len(root.Requirements)).ToBe(1)
	Expect(t, root.Requirements[0].Class).ToBe("InlineJavascriptRequirement")
	Expect(t, len(root.Inputs)).ToBe(1)
	Expect(t, root.Inputs[0].ID).ToBe("file1")
	Expect(t, root.Inputs[0].Types[0].Type).ToBe("File?")
	Expect(t, root.Outputs[0].ID).ToBe("out")
	Expect(t, root.Outputs[0].Types[0].Type).ToBe("string")
	Expect(t, root.Outputs[0].Binding.Glob[0]).ToBe("out.txt")
	Expect(t, root.Outputs[0].Binding.Contents).ToBe(false)
	Expect(t, root.Outputs[0].Binding.Eval).ToBe("$(self[0].contents)")
	Expect(t, root.Stdout).ToBe("out.txt")
	Expect(t, len(root.Arguments)).ToBe(2)
	Expect(t, root.Arguments[0].Value).ToBe("echo")
	Expect(t, root.Arguments[1].Value).ToBe(`$(inputs.file1 === null ? "t" : "f")`)
}
