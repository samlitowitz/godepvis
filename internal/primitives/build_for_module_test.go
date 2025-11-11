package primitives_test

import (
	"github.com/google/go-cmp/cmp"
	"github.com/samlitowitz/godepvis/internal/modfile"
	"github.com/samlitowitz/godepvis/internal/primitives"
	"github.com/samlitowitz/godepvis/internal/test"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestBuildForModule_NoCircularDependency(t *testing.T) {
	// REFURL: https://github.com/golang/go/blob/988b718f4130ab5b3ce5a5774e1a58e83c92a163/src/path/filepath/path_test.go#L600
	// -- START -- //
	if runtime.GOOS == "ios" {
		restore := test.Chtmpdir(t)
		defer restore()
	}

	tmpDir := t.TempDir()

	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal("finding working dir:", err)
	}
	if err = os.Chdir(tmpDir); err != nil {
		t.Fatal("entering temp dir:", err)
	}
	defer os.Chdir(origDir)
	// -- END -- //

	err = os.CopyFS(tmpDir, os.DirFS(filepath.Join(origDir, "testdata/no-circular-dependencies")))
	if err != nil {
		t.Fatal("copy test data: ", err)
	}

	goModFile, err := modfile.FindGoModFile(tmpDir)
	if err != nil {
		t.Fatal("failed to find go.mod: ", err)
	}
	modulePath, err := modfile.GetModulePath(goModFile)
	if err != nil {
		t.Fatal("failed to get module path: ", err)
	}

	actualPkgs, err := primitives.BuildForModule(modulePath, tmpDir, nil)
	if err != nil {
		t.Fatal("BuildForModule: ", err)
	}
	var expectedPkgs []*primitives.Package

	if diff := cmp.Diff(expectedPkgs, actualPkgs); diff != "" {
		t.Fatal(test.Mismatch("BuildForModule: ", diff))
	}
}
