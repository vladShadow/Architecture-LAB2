package archmodule

import (
	"bytes"
	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
	bin "github.com/Tornado9966/Lab2_Go/build/binmodule"
	"strings"
	"testing"
)

func TestSimpleArchiveFactory(t *testing.T) {
	ctx := blueprint.NewContext()

	ctx.MockFileSystem(map[string][]byte{
		"Blueprints": []byte(`
			go_binary {
			  name: "test-bin",
			  srcs: ["test-src.go"],
			  pkg: ".",
	                  vendorFirst: true
			}

			archive_bin {
			  name: "test-out",
			  binary: ["test-bin"],
			}
		`),
		"test-src.go": nil,
	})

	ctx.RegisterModuleType("go_binary", bin.SimpleBinFactory)
	ctx.RegisterModuleType("archive_bin", SimpleArchiveFactory)

	cfg := bood.NewConfig()

	_, errs := ctx.ParseBlueprintsFiles(".", cfg)
	if len(errs) != 0 {
		t.Fatalf("Syntax errors in the test blueprint file: %s", errs)
	}
 
	_, errs = ctx.PrepareBuildActions(cfg)
	if len(errs) != 0 {
		t.Errorf("Unexpected errors while preparing build actions: %s", errs)
	}
	buffer := new(bytes.Buffer)
	if err := ctx.WriteBuildFile(buffer); err != nil {
		t.Errorf("Error writing ninja file: %s", err)
	} else {
		text := buffer.String()
		t.Logf("Gennerated ninja build file:\n%s", text)
		if !strings.Contains(text, "out/archives/test-out.zip: ") {
			t.Errorf("Generated ninja file does not have build of the test module")
		}
	}
}
