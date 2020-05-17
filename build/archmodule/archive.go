package archmodule

import (
	"fmt"
	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
	bin "github.com/Tornado9966/Lab2_Go/build/binmodule"
	"path"
)

var (
	pctx = blueprint.NewPackageContext("github.com/Tornado9966/Lab2_Go/build/archmodule")

	goArchive =  pctx.StaticRule("ziparchive", blueprint.RuleParams{

		//You can use standard linux zip utility:
		//Command:     "cd $workDir && zip -r -j $output $inputsZip", 

		//Using our own archivator:
		Command:    "cd $workDir/cmd/archivate && go run archivate.go $output $inputsZip",

		Description: "zip archive",
	}, "workDir", "output", "inputsZip")
)

type ArchiveModule struct {
	blueprint.SimpleName

	properties struct {
		Binary []string
		Srcs []string
		SrcsExclude []string
	}
}

func (am *ArchiveModule) DynamicDependencies(blueprint.DynamicDependerModuleContext) []string {
	return am.properties.Binary
}

func (am *ArchiveModule) GenerateBuildActions(ctx blueprint.ModuleContext) {
	name := ctx.ModuleName()
	config := bood.ExtractConfig(ctx)
	config.Debug.Printf("Adding build actions for go binary module '%s'", name)

	output := path.Join(config.BaseOutputDir, "archives", name + ".zip")

	dep, _ := ctx.GetDirectDep(am.properties.Binary[0])
	inputsZip := ""
	for _, file := range dep.(*bin.BinaryModule).Outputs(){
		inputsZip += file + " "
	}

	var inputs []string
	inputErors := false
	for _, src := range am.properties.Srcs {
		if matches, err := ctx.GlobWithDeps(src, am.properties.SrcsExclude); err == nil {
			inputs = append(inputs, matches...)
		} else {
			ctx.PropertyErrorf("srcs", "Cannot resolve files that match pattern %s", src)
			inputErors = true
		}
	}
	if inputErors { 
		return
	} 

	ctx.Build(pctx, blueprint.BuildParams{ 
		Description: fmt.Sprintf("Archive zipped: %s", name),
		Rule:        goArchive,
		Outputs:     []string{output},
		Implicits:   inputs,
		Args: map[string]string{
			"output":     output,
			"workDir":    ctx.ModuleDir(),
			"inputsZip":  inputsZip,
		},
	})
}

func SimpleArchiveFactory() (blueprint.Module, []interface{}) {
	mType := &ArchiveModule{}
	return mType, []interface{}{&mType.SimpleName.Properties, &mType.properties}
}
