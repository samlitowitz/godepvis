package primitives_test

import (
	"github.com/google/go-cmp/cmp"
	"github.com/samlitowitz/godepvis/internal/primitives"
	"github.com/samlitowitz/godepvis/internal/test"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"testing"
)

var depVisTestdataDir = filepath.Join("testdata", "dependency-visitor")

func TestDependencyVisitor_Visit_ImportSpecs(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal("finding working dir:", err)
	}

	absFilepath := filepath.Join(cwd, depVisTestdataDir, "imports.go")

	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, absFilepath, nil, 0)
	if err != nil {
		t.Fatal("parse file: ", err)
	}
	depVis := primitives.NewDependencyVisitor()
	ast.Walk(depVis, astFile)

	var actualImportPaths []string
	for _, node := range depVis.InOrderNodes() {
		imp, ok := node.(*primitives.ImportSpec)
		if !ok {
			continue
		}
		actualImportPaths = append(actualImportPaths, imp.Name.String())
	}

	expectedImportPaths := []string{
		"log",
		"log",
		"os",
	}

	if diff := cmp.Diff(expectedImportPaths, actualImportPaths); diff != "" {
		t.Error(test.Mismatch("", diff))
	}
}

func TestDependencyVisitor_Visit_Constants(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal("finding working dir:", err)
	}

	absFilepath := filepath.Join(cwd, depVisTestdataDir, "constants.go")

	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, absFilepath, nil, 0)
	if err != nil {
		t.Fatal("parse file: ", err)
	}
	depVis := primitives.NewDependencyVisitor()
	ast.Walk(depVis, astFile)

	var actualConstants []string
	for _, node := range depVis.InOrderNodes() {
		declNode, ok := node.(*ast.GenDecl)
		if !ok {
			continue
		}
		if declNode.Tok != token.CONST {
			continue
		}
		for _, spec := range declNode.Specs {
			spec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for _, name := range spec.Names {
				actualConstants = append(actualConstants, name.String())
			}
		}
	}

	expectedConstants := []string{
		"CA",
		"CB",
		"CC",
		"CD",
		"ce",
	}

	if diff := cmp.Diff(expectedConstants, actualConstants); diff != "" {
		t.Error(test.Mismatch("", diff))
	}
}

func TestDependencyVisitor_Visit_Types(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal("finding working dir:", err)
	}

	absFilepath := filepath.Join(cwd, depVisTestdataDir, "types.go")

	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, absFilepath, nil, 0)
	if err != nil {
		t.Fatal("parse file: ", err)
	}
	depVis := primitives.NewDependencyVisitor()
	ast.Walk(depVis, astFile)

	var actualConstants []string
	for _, node := range depVis.InOrderNodes() {
		declNode, ok := node.(*ast.GenDecl)
		if !ok {
			continue
		}
		if declNode.Tok != token.TYPE {
			continue
		}
		for _, spec := range declNode.Specs {
			spec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			actualConstants = append(actualConstants, spec.Name.String())
		}
	}

	expectedConstants := []string{
		"TA",
		"TB",
		"tc",
	}

	if diff := cmp.Diff(expectedConstants, actualConstants); diff != "" {
		t.Error(test.Mismatch("", diff))
	}
}

func TestDependencyVisitor_Visit_Vars(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal("finding working dir:", err)
	}

	absFilepath := filepath.Join(cwd, depVisTestdataDir, "vars.go")

	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, absFilepath, nil, 0)
	if err != nil {
		t.Fatal("parse file: ", err)
	}
	depVis := primitives.NewDependencyVisitor()
	ast.Walk(depVis, astFile)

	var actualConstants []string
	for _, node := range depVis.InOrderNodes() {
		declNode, ok := node.(*ast.GenDecl)
		if !ok {
			continue
		}
		if declNode.Tok != token.VAR {
			continue
		}
		for _, spec := range declNode.Specs {
			spec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for _, name := range spec.Names {
				actualConstants = append(actualConstants, name.String())
			}
		}
	}

	expectedConstants := []string{
		"VA",
		"VB",
		"VC",
		"VD",
		"ve",
	}

	if diff := cmp.Diff(expectedConstants, actualConstants); diff != "" {
		t.Error(test.Mismatch("", diff))
	}
}

func TestDependencyVisitor_Visit_Functions(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal("finding working dir:", err)
	}

	absFilepath := filepath.Join(cwd, depVisTestdataDir, "functions.go")

	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, absFilepath, nil, 0)
	if err != nil {
		t.Fatal("parse file: ", err)
	}
	depVis := primitives.NewDependencyVisitor()
	ast.Walk(depVis, astFile)

	var actualConstants []string
	for _, node := range depVis.InOrderNodes() {
		declNode, ok := node.(*primitives.FuncDecl)
		if !ok {
			continue
		}
		actualConstants = append(actualConstants, declNode.Name.String())
	}

	expectedConstants := []string{
		"FTA1",
		"FTA2",
		"FC",
		"F1",
	}

	if diff := cmp.Diff(expectedConstants, actualConstants); diff != "" {
		t.Error(test.Mismatch("", diff))
	}
}

func TestDependencyVisitor_Visit_SelectorExprs(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal("finding working dir:", err)
	}

	absFilepath := filepath.Join(cwd, depVisTestdataDir, "selectorexprs.go")

	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, absFilepath, nil, 0)
	if err != nil {
		t.Fatal("parse file: ", err)
	}
	depVis := primitives.NewDependencyVisitor()
	ast.Walk(depVis, astFile)

	var actualConstants []string
	for _, node := range depVis.InOrderNodes() {
		selectorExpr, ok := node.(*primitives.SelectorExpr)
		if !ok {
			continue
		}
		x, ok := selectorExpr.X.(*ast.Ident)
		if !ok {
			t.Error("expected selector node to have X")
		}

		actualConstants = append(actualConstants, x.String()+"."+selectorExpr.Sel.String())
	}

	expectedConstants := []string{
		"fmt.Println",
		"fmt.Println",
	}

	if diff := cmp.Diff(expectedConstants, actualConstants); diff != "" {
		t.Error(test.Mismatch("", diff))
	}
}
