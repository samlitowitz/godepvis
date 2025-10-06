package ast

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"
)

type Package struct {
	*ast.Package

	DirName string
}

type File struct {
	*ast.File

	AbsPath string
	DirName string
}

type ImportSpec struct {
	*ast.ImportSpec

	IsAliased bool
	Alias     string
}

type FuncDecl struct {
	*ast.FuncDecl

	ReceiverName  string
	QualifiedName string
}

type SelectorExpr struct {
	*ast.SelectorExpr

	ImportName string
}

func (decl FuncDecl) IsReceiver() bool {
	return decl.Recv != nil
}

type DependencyVisitor struct {
	out chan<- ast.Node

	fileImports map[string]struct{}

	// Track current package context for filename lookup
	currentPackage *ast.Package
	currentDirName string
}

func NewDependencyVisitor() (*DependencyVisitor, <-chan ast.Node) {
	out := make(chan ast.Node)
	v := &DependencyVisitor{
		out: out,
	}

	return v, out
}

func (v *DependencyVisitor) Visit(node ast.Node) ast.Visitor {
	switch node := node.(type) {
	case *ast.Package:
		v.emitPackageAndFiles(node)

	case *ast.File:
		v.fileImports = make(map[string]struct{})
		v.emitFile(node)

	case *ast.ImportSpec:
		v.emitImportSpec(node)

	case *ast.FuncDecl:
		v.emitFuncDecl(node)
		// Don't descend into function bodies to avoid collecting function-scoped declarations
		return nil

	case *ast.GenDecl:
		switch node.Tok {
		case token.CONST:
			fallthrough
		case token.TYPE:
			fallthrough
		case token.VAR:
			v.out <- node
		}

	case *ast.SelectorExpr:
		// only references to external packages
		if node.X == nil {
			return v
		}

		impName := ""
		switch x := node.X.(type) {
		case *ast.Ident:
			impName = x.String()
		}

		// if the "import name" is actually a variable and not a package, skip it
		if _, ok := v.fileImports[impName]; !ok {
			return v
		}

		v.out <- &SelectorExpr{
			SelectorExpr: node,
			ImportName:   impName,
		}
	}
	return v
}

func (v *DependencyVisitor) emitPackageAndFiles(node *ast.Package) {
	// Get directory name from first file
	var dirName string
	for filename := range node.Files {
		absPath, err := filepath.Abs(filename)
		if err != nil {
			continue
		}
		dirName, _ = filepath.Split(absPath)
		dirName = strings.TrimRight(dirName, "/")
		break
	}

	// Store package context for later file emission
	v.currentPackage = node
	v.currentDirName = dirName

	v.out <- &Package{
		Package: node,
		DirName: dirName,
	}
}

func (v *DependencyVisitor) emitFile(node *ast.File) {
	if v.currentPackage == nil {
		return
	}

	// Find the filename for this ast.File in the package's Files map
	var filename string
	for fn, astFile := range v.currentPackage.Files {
		if astFile == node {
			filename = fn
			break
		}
	}

	if filename == "" {
		return
	}

	absPath, err := filepath.Abs(filename)
	if err != nil {
		return
	}

	v.out <- &File{
		File:    node,
		AbsPath: absPath,
		DirName: v.currentDirName,
	}
}

func (v *DependencyVisitor) emitImportSpec(node *ast.ImportSpec) {
	node.Path.Value = strings.Trim(node.Path.Value, "\"")
	pieces := strings.Split(node.Path.Value, "/")
	name := pieces[len(pieces)-1]

	isAliased := node.Name != nil
	alias := ""

	if isAliased {
		alias = node.Name.String()
		node.Name.Name = name
		v.fileImports[alias] = struct{}{}
	}

	if !isAliased {
		node.Name = &ast.Ident{
			Name: name,
		}
		v.fileImports[name] = struct{}{}
	}

	v.out <- &ImportSpec{
		ImportSpec: node,
		IsAliased:  isAliased,
		Alias:      alias,
	}
}

func (v *DependencyVisitor) emitFuncDecl(node *ast.FuncDecl) {
	receiverName := ""
	qualifiedName := node.Name.String()

	if node.Recv != nil {
		// TODO: don't emit receiver functions/methods? we don't need them
		var typName string
		recvType := node.Recv.List[0].Type
		switch expr := recvType.(type) {
		case *ast.Ident:
			typName = expr.String()
		case *ast.StarExpr:
			typName = extractTypeName(expr.X)
		case *ast.IndexExpr:
			// Generic type with one parameter: Foo[T]
			typName = extractTypeName(expr)
		case *ast.IndexListExpr:
			// Generic type with multiple parameters: Foo[T, U]
			typName = extractTypeName(expr)
		default:
			// panic error, invalid receiver method
		}
		receiverName = typName
		qualifiedName = typName + "." + node.Name.String()
	}

	v.out <- &FuncDecl{
		FuncDecl:      node,
		ReceiverName:  receiverName,
		QualifiedName: qualifiedName,
	}
}

func (v *DependencyVisitor) Close() {
	close(v.out)
}

// extractTypeName extracts the type name from an expression, handling generics
func extractTypeName(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		// Simple type: Foo
		return e.String()
	case *ast.IndexExpr:
		// Generic type with one parameter: Foo[T]
		return extractTypeName(e.X)
	case *ast.IndexListExpr:
		// Generic type with multiple parameters: Foo[T, U]
		return extractTypeName(e.X)
	default:
		return ""
	}
}
