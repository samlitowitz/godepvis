package primitives_test

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/samlitowitz/godepvis/internal/modfile"
	"github.com/samlitowitz/godepvis/internal/primitives"
	"github.com/samlitowitz/godepvis/internal/test"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestBuildForModule(t *testing.T) {
	opts := cmp.Options{
		cmpopts.SortSlices(
			func(x, y string) bool {
				return strings.Compare(x, y) <= 0
			},
		),
	}

	testCases := map[string]struct {
		dir           string
		expectedNames []string
	}{
		"direct-circular-dependency": {
			dir:           "direct-circular-dependency",
			expectedNames: []string{"a", "b", "log", "main"},
		},
		"direct-circular-dependency-blank-identifiers": {
			dir:           "direct-circular-dependency-blank-identifiers",
			expectedNames: []string{"a", "b", "log", "main"},
		},
		"direct-circular-dependency-with-blank-identifier": {
			dir:           "direct-circular-dependency-with-blank-identifier",
			expectedNames: []string{"a", "b", "log", "main"},
		},
		"multiple-independent-direct-circular-dependencies": {
			dir:           "multiple-independent-direct-circular-dependencies",
			expectedNames: []string{"a", "b", "c", "log", "main"},
		},
		"multiple-interlinked-direct-circular-dependencies": {
			dir:           "multiple-interlinked-direct-circular-dependencies",
			expectedNames: []string{"a", "b", "c", "log", "main"},
		},
		"multiple-interlinked-direct-circular-dependencies-with-blank-identifier": {
			dir:           "multiple-interlinked-direct-circular-dependencies-with-blank-identifier",
			expectedNames: []string{"a", "b", "c", "log", "main"},
		},
		"no-circular-dependencies": {
			dir:           "no-circular-dependencies",
			expectedNames: []string{"a", "b", "c", "log", "main"},
		},
		"no-circular-dependencies-with-blank-identifier": {
			dir:           "no-circular-dependencies-with-blank-identifier",
			expectedNames: []string{"a", "b", "c", "log", "main"},
		},
		"transitive-circular-dependency": {
			dir:           "transitive-circular-dependency",
			expectedNames: []string{"a", "b", "c", "log", "main"},
		},
		"with-test-packages": {
			dir:           "with-test-packages",
			expectedNames: []string{"a", "b", "c", "log", "main"},
		},
	}

	for desc, testCase := range testCases {
		func() {
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

			err = os.CopyFS(tmpDir, os.DirFS(filepath.Join(origDir, "testdata", "build-for-module", testCase.dir)))
			if err != nil {
				t.Fatal("copy test data:", err)
			}

			moduleDir := tmpDir

			goModFile, err := modfile.FindGoModFile(moduleDir)
			if err != nil {
				t.Fatal(desc, ": failed to find go.mod: ", err)
			}
			modulePath, err := modfile.GetModulePath(goModFile)
			if err != nil {
				t.Fatal(desc, ": failed to get module path: ", err)
			}

			actualPkgs, err := primitives.BuildForModule(modulePath, moduleDir)
			if err != nil {
				t.Fatal(desc, ": BuildForModule: ", err)
			}

			var actualNames []string
			for _, pkg := range actualPkgs {
				actualNames = append(actualNames, pkg.Name)
			}

			if diff := cmp.Diff(testCase.expectedNames, actualNames, opts); diff != "" {
				t.Error(desc, test.Mismatch(": expected packages: ", diff))
			}
		}()
	}
}

func TestBuildForModule_WithCorrectBlankImports(t *testing.T) {
	opts := cmp.Options{
		cmpopts.SortSlices(
			func(x, y string) bool {
				return strings.Compare(x, y) <= 0
			},
		),
	}

	testCases := map[string]struct {
		dir                                  string
		expectedBlankImportsByPackage        []string
		expectedBlankImportsInCycleByPackage []string
	}{
		"direct-circular-dependency": {
			dir: "direct-circular-dependency",
		},
		"direct-circular-dependency-blank-identifiers": {
			dir:                                  "direct-circular-dependency-blank-identifiers",
			expectedBlankImportsByPackage:        []string{"a", "b"},
			expectedBlankImportsInCycleByPackage: []string{"a", "b"},
		},
		"direct-circular-dependency-with-blank-identifier": {
			dir:                                  "direct-circular-dependency-with-blank-identifier",
			expectedBlankImportsByPackage:        []string{"b"},
			expectedBlankImportsInCycleByPackage: []string{"b"},
		},
		"multiple-independent-direct-circular-dependencies": {
			dir: "multiple-independent-direct-circular-dependencies",
		},
		"multiple-interlinked-direct-circular-dependencies": {
			dir: "multiple-interlinked-direct-circular-dependencies",
		},
		"multiple-interlinked-direct-circular-dependencies-with-blank-identifier": {
			dir:                                  "multiple-interlinked-direct-circular-dependencies-with-blank-identifier",
			expectedBlankImportsByPackage:        []string{"a", "b", "c"},
			expectedBlankImportsInCycleByPackage: []string{"a", "b", "c"},
		},
		"no-circular-dependencies": {
			dir: "no-circular-dependencies",
		},
		"no-circular-dependencies-with-blank-identifier": {
			dir:                           "no-circular-dependencies-with-blank-identifier",
			expectedBlankImportsByPackage: []string{"b"},
		},
		"transitive-circular-dependency": {
			dir: "transitive-circular-dependency",
		},
		"with-test-packages": {
			dir: "with-test-packages",
		},
	}

	for desc, testCase := range testCases {
		func() {
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

			err = os.CopyFS(tmpDir, os.DirFS(filepath.Join(origDir, "testdata", "build-for-module", testCase.dir)))
			if err != nil {
				t.Fatal("copy test data:", err)
			}

			moduleDir := tmpDir

			goModFile, err := modfile.FindGoModFile(moduleDir)
			if err != nil {
				t.Fatal(desc, ": failed to find go.mod: ", err)
			}
			modulePath, err := modfile.GetModulePath(goModFile)
			if err != nil {
				t.Fatal(desc, ": failed to get module path: ", err)
			}

			actualPkgs, err := primitives.BuildForModule(modulePath, moduleDir)
			if err != nil {
				t.Fatal(desc, ": BuildForModule: ", err)
			}

			var actualBlankImportsByPackage []string
			var actualBlankImportsInCycleByPackage []string
			for _, pkg := range actualPkgs {
				if pkg.BlankImportFile == nil {
					continue
				}
				actualBlankImportsByPackage = append(actualBlankImportsByPackage, pkg.Name)
				if !pkg.BlankImportFile.InImportCycle {
					continue
				}
				actualBlankImportsInCycleByPackage = append(actualBlankImportsInCycleByPackage, pkg.Name)
			}

			if diff := cmp.Diff(testCase.expectedBlankImportsByPackage, actualBlankImportsByPackage, opts); diff != "" {
				t.Error(desc, test.Mismatch(": expected blank imports by package: ", diff))
			}
			if diff := cmp.Diff(testCase.expectedBlankImportsInCycleByPackage, actualBlankImportsInCycleByPackage, opts); diff != "" {
				t.Error(desc, test.Mismatch(": expected blank imports in cycle by package: ", diff))
			}
		}()
	}
}
