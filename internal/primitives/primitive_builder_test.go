package primitives_test

import (
	"github.com/google/go-cmp/cmp"
	"github.com/samlitowitz/godepvis/internal"
	"github.com/samlitowitz/godepvis/internal/modfile"
	"github.com/samlitowitz/godepvis/internal/primitives"
	"github.com/samlitowitz/godepvis/internal/test"
	"go/ast"
	"os"
	"path/filepath"
	"testing"
)

var primitiveBuilderTestdataDir = filepath.Join("testdata", "primitive-builder")

func TestPrimitiveBuilder_AddNode_Package(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal("finding working dir:", err)
	}

	moduleRoot := filepath.Join(cwd, primitiveBuilderTestdataDir)

	goModFile, err := modfile.FindGoModFile(cwd)
	if err != nil {
		t.Fatal("failed to find go.mod: ", err)
	}
	modulePath, err := modfile.GetModulePath(goModFile)
	if err != nil {
		t.Fatal("failed to get module path: ", err)
	}

	testCases := map[string]struct {
		expected *internal.Package
	}{
		"main": {
			expected: &internal.Package{
				DirName:    moduleRoot,
				ModulePath: modulePath,
				ModuleDir:  moduleRoot,
				Name:       "main",
				Files:      map[string]*internal.File{},
			},
		},
		"a": {
			expected: &internal.Package{
				DirName:    moduleRoot,
				ModulePath: modulePath,
				ModuleDir:  moduleRoot,
				Name:       "a",
				Files:      map[string]*internal.File{},
			},
		},
	}

	for desc, testCase := range testCases {
		builder := primitives.NewPrimitiveBuilder(modulePath, moduleRoot)
		expected := testCase.expected
		err = builder.AddNode(&primitives.Package{
			Package: &ast.Package{
				Name: expected.Name,
			},
			DirName: expected.DirName,
		})

		actualPackages := builder.Packages()

		if diff := cmp.Diff([]*internal.Package{expected}, actualPackages); diff != "" {
			t.Error(test.Mismatch(desc+": ", diff))
		}
	}
}

func TestPrimitiveBuilder_AddNode_File(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal("finding working dir:", err)
	}

	moduleRoot := filepath.Join(cwd, primitiveBuilderTestdataDir)

	goModFile, err := modfile.FindGoModFile(cwd)
	if err != nil {
		t.Fatal("failed to find go.mod: ", err)
	}
	modulePath, err := modfile.GetModulePath(goModFile)
	if err != nil {
		t.Fatal("failed to get module path: ", err)
	}

	testCases := map[string]struct {
		pkg   *internal.Package
		files []*internal.File
	}{
		"main": {
			pkg: &internal.Package{
				DirName:    moduleRoot,
				ModulePath: modulePath,
				ModuleDir:  moduleRoot,
				Name:       "main",
				Files:      map[string]*internal.File{},
			},
			files: []*internal.File{
				{
					FileName: "main.go",
					Imports:  map[string]*internal.Import{},
					Decls:    map[string]*internal.Decl{},
				},
			},
		},
		"a": {
			pkg: &internal.Package{
				DirName:    moduleRoot,
				ModulePath: modulePath,
				ModuleDir:  moduleRoot,
				Name:       "a",
				Files:      map[string]*internal.File{},
			},
			files: []*internal.File{
				{
					FileName: "a_1.go",
					Imports:  map[string]*internal.Import{},
					Decls:    map[string]*internal.Decl{},
				},
				{
					FileName: "a_2.go",
					Imports:  map[string]*internal.Import{},
					Decls:    map[string]*internal.Decl{},
				},
			},
		},
	}

	for desc, testCase := range testCases {
		builder := primitives.NewPrimitiveBuilder(modulePath, moduleRoot)
		expected := testCase.pkg
		err = builder.AddNode(&primitives.Package{
			Package: &ast.Package{
				Name: expected.Name,
			},
			DirName: expected.DirName,
		})
		if err != nil {
			t.Fatal("add package: ", err)
		}

		for _, file := range testCase.files {
			absFilepath := filepath.Join(expected.DirName, file.FileName)
			err = builder.AddNode(&primitives.File{
				File:    &ast.File{},
				AbsPath: absFilepath,
				DirName: expected.DirName,
			})
			if err != nil {
				t.Fatal("add file: ", err)
			}
			file.AbsPath = absFilepath
			file.Package = expected
			expected.Files[file.UID()] = file
		}

		actualPackages := builder.Packages()

		if diff := cmp.Diff([]*internal.Package{expected}, actualPackages); diff != "" {
			t.Error(test.Mismatch(desc+": ", diff))
		}
	}
}
