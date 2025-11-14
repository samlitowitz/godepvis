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

	actualConstantsByFuncScopeName := make(map[string][]string)
	for _, node := range depVis.InOrderNodes() {
		declNode, ok := node.(*primitives.GenDecl)
		if !ok {
			continue
		}
		if declNode.Tok != token.CONST {
			continue
		}
		if _, ok := actualConstantsByFuncScopeName[declNode.FuncScopeName]; !ok {
			actualConstantsByFuncScopeName[declNode.FuncScopeName] = make([]string, 0, len(declNode.Specs))
		}
		for _, spec := range declNode.Specs {
			spec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for _, name := range spec.Names {
				actualConstantsByFuncScopeName[declNode.FuncScopeName] = append(actualConstantsByFuncScopeName[declNode.FuncScopeName], name.String())
			}
		}
	}

	expectedConstantsByFuncScopeName := map[string][]string{
		"":        {"CA", "CB", "CC", "CD", "ce"},
		"2":       {"ce"},
		"3":       {"ce"},
		"CFn1":    {"ce"},
		"CFn1.14": {"ce"},
		"CFn1.17": {"ce"},
		"CFn2":    {"ce"},
	}

	if diff := cmp.Diff(expectedConstantsByFuncScopeName, actualConstantsByFuncScopeName); diff != "" {
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

	actualTypesByFuncScopeName := make(map[string][]string)
	for _, node := range depVis.InOrderNodes() {
		declNode, ok := node.(*primitives.GenDecl)
		if !ok {
			continue
		}
		if declNode.Tok != token.TYPE {
			continue
		}
		if _, ok := actualTypesByFuncScopeName[declNode.FuncScopeName]; !ok {
			actualTypesByFuncScopeName[declNode.FuncScopeName] = make([]string, 0, len(declNode.Specs))
		}
		for _, spec := range declNode.Specs {
			spec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			actualTypesByFuncScopeName[declNode.FuncScopeName] = append(actualTypesByFuncScopeName[declNode.FuncScopeName], spec.Name.String())
		}
	}

	expectedTypesByFuncScopeName := map[string][]string{
		"":     {"TA", "TB", "tc"},
		"0":    {"tc"},
		"1":    {"tc"},
		"TFn1": {"tc"},
		"TFn2": {"tc"},
	}

	if diff := cmp.Diff(expectedTypesByFuncScopeName, actualTypesByFuncScopeName); diff != "" {
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

	actualVarsByFuncScopeName := make(map[string][]string)
	for _, node := range depVis.InOrderNodes() {
		declNode, ok := node.(*primitives.GenDecl)
		if !ok {
			continue
		}
		if declNode.Tok != token.VAR {
			continue
		}
		if _, ok := actualVarsByFuncScopeName[declNode.FuncScopeName]; !ok {
			actualVarsByFuncScopeName[declNode.FuncScopeName] = make([]string, 0, len(declNode.Specs))
		}
		for _, spec := range declNode.Specs {
			spec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for _, name := range spec.Names {
				actualVarsByFuncScopeName[declNode.FuncScopeName] = append(actualVarsByFuncScopeName[declNode.FuncScopeName], name.String())
			}
		}
	}

	expectedVarsByFuncScopeName := map[string][]string{
		"":     {"VA", "VB", "VC", "VD", "ve", "_", "_"},
		"0":    {"ve"},
		"1":    {"ve"},
		"VFn1": {"ve"},
		"VFn2": {"ve"},
	}

	if diff := cmp.Diff(expectedVarsByFuncScopeName, actualVarsByFuncScopeName); diff != "" {
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
