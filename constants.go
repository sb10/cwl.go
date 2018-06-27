package cwl

// *** commented out constants are in the CWL spec, but not yet handled here...

// workflow field constants
const (
	fieldInputs         = "inputs"
	fieldOutputs        = "outputs"
	fieldClass          = "class"
	fieldSteps          = "steps"
	fieldID             = "id"
	fieldRequirements   = "requirements"
	fieldHints          = "hints"
	fieldLabel          = "label"
	fieldDoc            = "doc"
	fieldCWLVersion     = "cwlVersion"
	fieldSecondaryFiles = "secondaryFiles"
	// fieldStreamable     = "streamable"
	fieldOutputBinding = "outputBinding"
	fieldFormat        = "format"
	fieldOutputSource  = "outputSource"
	fieldLinkMerge     = "linkMerge"
	fieldType          = "type"
	fieldGlob          = "glob"
	fieldLoadContents  = "loadContents"
	fieldOutputEval    = "outputEval"
	fieldLocation      = "location"
	// fieldPath           = "path"
	// fieldBasename       = "basename"
	// fieldDirname        = "dirname"
	// fieldNameRoot       = "nameroot"
	// fieldNameExt        = "nameext"
	// fieldChecksum       = "checksum"
	// fieldSize           = "size"
	// fieldContents       = "contents"
	fieldListing       = "listing"
	fieldSource        = "source"
	fieldFields        = "fields"
	fieldName          = "name"
	fieldSymbols       = "symbols"
	fieldItems         = "items"
	fieldIn            = "in"
	fieldOut           = "out"
	fieldRun           = "run"
	fieldScatter       = "scatter"
	fieldScatterMethod = "scatterMethod"
	fieldExpressionLib = "expressionLib"
	fieldPosition      = "position"
	fieldPrefix        = "prefix"
	// fieldSeparate       = "separate"
	fieldItemSeparator = "itemSeparator"
	fieldShellQuote    = "shellQuote"
	fieldDefault       = "default"
	fieldValueFrom     = "valueFrom"
	// fieldPackages       = "packages"
	// fieldPackage        = "package"
	// fieldVersion        = "vesion"
	// fieldSpecs          = "specs"
	fieldEntry        = "entry"
	fieldEntryName    = "entryname"
	fieldWritable     = "writable"
	fieldInputBinding = "inputBinding"
	fieldTypes        = "types"
	fieldExpression   = "expression"
)

// commandlinetool field constants
const (
	fieldBaseCommand = "baseCommand"
	fieldArguments   = "arguments"
	fieldStdIn       = "stdin"
	fieldStdErr      = "stderr"
	fieldStdOut      = "stdout"
	// fieldSuccessCodes          = "successCodes"
	// fieldTemporaryFailCodes    = "temporaryFailCodes"
	// fieldPermanentFailCodes    = "permanentFailCodes"
	fieldDockerPull = "dockerPull"
	// fieldDockerLoad            = "dockerLoad"
	// fieldDockerFile            = "dockerFile"
	// fieldDockerImport          = "dockerImport"
	// fieldDockerImageID         = "dockerImageId"
	fieldDockerOutputDirectory = "dockerOutputDirectory"
	fieldEnvDef                = "envDef"
	// fieldEnvName               = "envName"
	// fieldEnvValue              = "envValue"
	fieldCoresMin = "coresMin"
	// fieldCoresMax              = "coresMax"
	// fieldRamMin                = "ramMin"
	// fieldRamMax                = "ramMax"
	// fieldTmpDirMin             = "tmpdirMin"
	// fieldTmpDirMax             = "tmpdirMax"
	// fieldOutDirMin             = "outdirMin"
	// fieldOutDirMax             = "outdirMax"
)

// cwltype symbols
const (
	// typeNull      = "null"
	// typeBoolean   = "boolean"
	typeInt = "int"
	// typeLong      = "long"
	// typeFloat     = "float"
	// typeDouble    = "double"
	typeString    = "string"
	typeFile      = "File"
	typeFileSlice = "File[]"
	// typeDirectory = "directory"
	// typeArray = "array"
)

// classes
const (
	classWorkflow   = "Workflow"
	classCommand    = "CommandLineTool"
	classExpression = "ExpressionTool"
)

// requirements
const (
	reqShell    = "ShellCommandRequirement"
	reqWorkDir  = "InitialWorkDirRequirement"
	reqEnv      = "EnvVarRequirement"
	reqJS       = "InlineJavascriptRequirement"
	reqScatter  = "ScatterFeatureRequirement"
	reqMultiple = "MultipleInputFeatureRequirement"
	// reqSubWorkFlow = "SubworkflowFeatureRequirement"
)

// merges
const (
	mergeFlattened = "merge_flattened"
)

// scatter methods
const (
	scatterDotProduct         = "dotproduct"
	scatterNestedCrossProduct = "nested_crossproduct"
	scatterFlatCrossProduct   = "flat_crossproduct"
	scatterNestedInput        = "_!_scatter_nested_input_!_"
	scatterFlatInput          = "_!_scatter_flat_input_!_"
)
