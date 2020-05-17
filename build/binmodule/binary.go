package binmodule

import (
	"fmt"
	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
	"path"
)

var (
	pctx = blueprint.NewPackageContext("github.com/Tornado9966/Lab2_Go/build/binmodule")

	goBuild = pctx.StaticRule("binaryBuild", blueprint.RuleParams{
		Command:     "cd $workDir && go build -o $outputPath $pkg",
		Description: "build go command $pkg",
	}, "workDir", "outputPath", "pkg")

	goTest = pctx.StaticRule("gotest", blueprint.RuleParams{
		Command:     "cd $workDir && go test -v -bench=. -benchtime=100x $pkgTest > $outputFile",
		Description: "test $pkgTest",
	}, "workDir", "pkgTest", "outputFile")

	goVendor = pctx.StaticRule("vendor", blueprint.RuleParams{
		Command:     "cd $workDir && go mod vendor",
		Description: "vendor dependencies of $name",
	}, "workDir", "name")

	config *(bood.Config)
)

type BinaryModule struct {
	blueprint.SimpleName

	properties struct {
		Pkg string
                TestPkg string
		OutTestFile string
		Srcs []string
		SrcsExclude []string
		Optional bool
		VendorFirst bool
	}
} 

func (bm *BinaryModule) GenerateBuildActions(ctx blueprint.ModuleContext) {
	name := ctx.ModuleName()
	config = bood.ExtractConfig(ctx)
	config.Debug.Printf("Adding build actions for go binary module '%s'", name)

	outputPath := path.Join(config.BaseOutputDir, "bin", name)
	outputFile := path.Join(config.BaseOutputDir, "bin", bm.properties.OutTestFile)

	var inputs []string
	inputErors := false
	for _, src := range bm.properties.Srcs {
		if matches, err := ctx.GlobWithDeps(src, bm.properties.SrcsExclude); err == nil {
			inputs = append(inputs, matches...)
		} else {
			ctx.PropertyErrorf("srcs", "Cannot resolve files that match pattern %s", src)
			inputErors = true
		}
	}
	if inputErors {
		return
	}

	var inputsTest []string
	for _, src := range bm.properties.Srcs {
		if matches, err := ctx.GlobWithDeps(src, nil); err == nil {
			inputsTest = append(inputsTest, matches...)
		}
	}

	if bm.properties.VendorFirst {
		vendorDirPath := path.Join(ctx.ModuleDir(), "vendor")
		ctx.Build(pctx, blueprint.BuildParams{
			Description: fmt.Sprintf("Vendor dependencies of %s", name),
			Rule:        goVendor,
			Outputs:     []string{vendorDirPath},
			Implicits:   []string{path.Join(ctx.ModuleDir(), "go.mod")},
			Optional:    true,
			Args: map[string]string{
				"workDir": ctx.ModuleDir(),
				"name":    name,
			},
		})
		inputs = append(inputs, vendorDirPath)
	}

	if(len(inputs) != 0) { 
		ctx.Build(pctx, blueprint.BuildParams{
			Description: fmt.Sprintf("Build %s as Go binary", name),
			Rule:        goBuild,
			Outputs:     []string{outputPath},
			Implicits:   inputs,
			Args: map[string]string{
				"outputPath": outputPath,
				"workDir":    ctx.ModuleDir(),
				"pkg":        bm.properties.Pkg,
			},
		})
	}

	ctx.Build(pctx, blueprint.BuildParams{
		Description: fmt.Sprintf("Test my module"),
		Rule:        goTest,
		Outputs:     []string{outputFile},
		Implicits:   inputsTest,
		Optional:    bm.properties.Optional,
		Args: map[string]string{
			"outputFile":     outputFile,
			"workDir":        ctx.ModuleDir(),
			"pkgTest":        bm.properties.TestPkg,
		},
	})
} 
 
func (bm *BinaryModule) Outputs() []string {
  	return []string{path.Join(config.BaseOutputDir, "bin", bm.Name())}
}

func SimpleBinFactory() (blueprint.Module, []interface{}) {
	mType := &BinaryModule{}
	return mType, []interface{}{&mType.SimpleName.Properties, &mType.properties}
}
