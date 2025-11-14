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
		"direct-circular-dependency-with-fn-receivers": {
			dir:           "direct-circular-dependency-with-fn-receivers",
			expectedNames: []string{"a", "b", "log", "main"},
		},
		"direct-circular-dependency-with-type-assert": {
			dir:           "direct-circular-dependency-with-type-assert",
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
		"with-fn-scoped": {
			dir:           "with-fn-scoped",
			expectedNames: []string{"a", "b", "main"},
		},
		"with-generics": {
			dir:           "with-generics",
			expectedNames: []string{"a", "b", "c", "log", "main"},
		},
		"with-types": {
			dir:           "with-types",
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
		"direct-circular-dependency-with-fn-receivers": {
			dir: "direct-circular-dependency-with-fn-receivers",
		},
		"direct-circular-dependency-with-type-assert": {
			dir: "direct-circular-dependency-with-type-assert",
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
		"with-fn-scoped": {
			dir: "with-fn-scoped",
		},
		"with-generics": {
			dir: "with-generics",
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

func TestBuildForModule_WithCorrectDecls(t *testing.T) {
	opts := cmp.Options{
		cmpopts.SortSlices(
			func(x, y string) bool {
				return strings.Compare(x, y) <= 0
			},
		),
	}

	testCases := map[string]struct {
		dir           string
		expectedDecls map[string][]string
	}{
		"direct-circular-dependency": {
			dir: "direct-circular-dependency",
			expectedDecls: map[string][]string{
				"github.com/fake/fake/a": {"Fn"},
				"github.com/fake/fake/b": {"Fn"},
				"log":                    {"Println"},
				"main":                   {"main"},
			},
		},
		"direct-circular-dependency-blank-identifiers": {
			dir: "direct-circular-dependency-blank-identifiers",
			expectedDecls: map[string][]string{
				"github.com/fake/fake/a": {"_", "Fn"},
				"github.com/fake/fake/b": {"_", "Fn"},
				"log":                    {"Println"},
				"main":                   {"main"},
			},
		},
		"direct-circular-dependency-with-blank-identifier": {
			dir: "direct-circular-dependency-with-blank-identifier",
			expectedDecls: map[string][]string{
				"github.com/fake/fake/a": {"Fn"},
				"github.com/fake/fake/b": {"_", "Fn"},
				"log":                    {"Println"},
				"main":                   {"main"},
			},
		},
		"direct-circular-dependency-with-fn-receivers": {
			dir: "direct-circular-dependency-with-fn-receivers",
			expectedDecls: map[string][]string{
				"github.com/fake/fake/a": {"A"},
				"github.com/fake/fake/b": {"B"},
				"log":                    {"Println"},
				"main":                   {"main"},
			},
		},
		"multiple-independent-direct-circular-dependencies": {
			dir: "multiple-independent-direct-circular-dependencies",
			expectedDecls: map[string][]string{
				"github.com/fake/fake/a": {"Fn"},
				"github.com/fake/fake/b": {"Fn"},
				"github.com/fake/fake/c": {"Fn"},
				"log":                    {"Println"},
				"main":                   {"main"},
			},
		},
		"multiple-interlinked-direct-circular-dependencies": {
			dir: "multiple-interlinked-direct-circular-dependencies",
			expectedDecls: map[string][]string{
				"github.com/fake/fake/a": {"Fn"},
				"github.com/fake/fake/b": {"Fn"},
				"github.com/fake/fake/c": {"Fn"},
				"log":                    {"Println"},
				"main":                   {"main"},
			},
		},
		"multiple-interlinked-direct-circular-dependencies-with-blank-identifier": {
			dir: "multiple-interlinked-direct-circular-dependencies-with-blank-identifier",
			expectedDecls: map[string][]string{
				"github.com/fake/fake/a": {"_", "Fn"},
				"github.com/fake/fake/b": {"_", "Fn"},
				"github.com/fake/fake/c": {"_", "Fn"},
				"log":                    {"Println"},
				"main":                   {"main"},
			},
		},
		"no-circular-dependencies": {
			dir: "no-circular-dependencies",
			expectedDecls: map[string][]string{
				"github.com/fake/fake/a": {"Fn"},
				"github.com/fake/fake/b": {"Fn"},
				"github.com/fake/fake/c": {"Fn1", "Fn2", "Fn3"},
				"log":                    {"Println"},
				"main":                   {"main"},
			},
		},
		"no-circular-dependencies-with-blank-identifier": {
			dir: "no-circular-dependencies-with-blank-identifier",
			expectedDecls: map[string][]string{
				"github.com/fake/fake/a": {"Fn"},
				"github.com/fake/fake/b": {"Fn", "_"},
				"github.com/fake/fake/c": {"Fn1", "Fn2", "Fn3"},
				"log":                    {"Println"},
				"main":                   {"main"},
			},
		},
		"transitive-circular-dependency": {
			dir: "transitive-circular-dependency",
			expectedDecls: map[string][]string{
				"github.com/fake/fake/a": {"Fn"},
				"github.com/fake/fake/b": {"Fn"},
				"github.com/fake/fake/c": {"Fn"},
				"log":                    {"Println"},
				"main":                   {"main"},
			},
		},
		"with-fn-scoped": {
			dir: "with-fn-scoped",
			expectedDecls: map[string][]string{
				"github.com/fake/fake/a": {"A"},
				"github.com/fake/fake/b": {
					"1.aS",
					"1.b",
					"1.c",
					"2.aS",
					"2.b",
					"2.c",
					"Fn1",
					"Fn1.aS",
					"Fn1.b",
					"Fn1.c",
					"Fn2",
					"Fn2.aS",
					"Fn2.b",
					"Fn2.c",
					"Fn3",
					"Fn3.7.aS",
					"Fn3.7.b",
					"Fn3.7.c",
				},

				"main": {"main"},
			},
		},
		"with-generics": {
			dir: "with-generics",
			expectedDecls: map[string][]string{
				"github.com/fake/fake/a": {"IsGreater", "Popper", "Stack", "Number", "Slice", "Clip"},
				"github.com/fake/fake/b": {"Fn", "gtFn", "gtV", "st", "sl", "c", "Sum", "Sum.s", "Product", "Product.s"},
				"github.com/fake/fake/c": {"Fn1", "Fn2", "Fn3"},
				"log":                    {"Println"},
				"main":                   {"main"},
			},
		},
		"with-types": {
			dir: "with-types",
			expectedDecls: map[string][]string{
				"github.com/fake/fake/a": {"AStruct", "ATStruct", "AString", "Fn"},
				"github.com/fake/fake/b": {"Fn", "aStruct", "aTStruct", "aString"},
				"github.com/fake/fake/c": {"Fn2", "Fn3", "Fn1"},
				"log":                    {"Println"},
				"main":                   {"main"},
			},
		},
		"with-test-packages": {
			dir: "with-test-packages",
			expectedDecls: map[string][]string{
				"github.com/fake/fake/a": {"Fn"},
				"github.com/fake/fake/b": {"Fn"},
				"github.com/fake/fake/c": {"Fn1", "Fn2", "Fn3"},
				"log":                    {"Println"},
				"main":                   {"main"},
			},
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

			actualDecls := make(map[string][]string)
			for _, pkg := range actualPkgs {
				pkgUID := pkg.UID()
				if pkg.Name == "main" {
					pkgUID = "main"
				}
				if _, ok := actualDecls[pkgUID]; !ok {
					actualDecls[pkgUID] = []string{}
				}
				for _, file := range pkg.Files {
					for _, decl := range file.Decls {
						actualDecls[pkgUID] = append(actualDecls[pkgUID], decl.UID())
					}
				}
			}

			if diff := cmp.Diff(testCase.expectedDecls, actualDecls, opts); diff != "" {
				t.Error(desc, test.Mismatch(": expected packages: ", diff))
			}
		}()
	}
}
