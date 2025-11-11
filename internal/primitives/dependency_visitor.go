package primitives

import (
	"go/ast"
	"go/token"
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

// DependencyVisitor is only for use against individual files
type DependencyVisitor struct {
	inOrderNodes []ast.Node

	fileImports map[string]struct{}
}

func NewDependencyVisitor() *DependencyVisitor {
	v := &DependencyVisitor{}

	return v
}

func (v *DependencyVisitor) Reset() {
	v.inOrderNodes = nil
	v.fileImports = nil
}

func (v *DependencyVisitor) InOrderNodes() []ast.Node {
	return v.inOrderNodes
}

func (v *DependencyVisitor) Visit(node ast.Node) ast.Visitor {
	switch node := node.(type) {

	case *ast.File:
		v.fileImports = make(map[string]struct{})

	case *ast.ImportSpec:
		v.addImportSpec(node)

	case *ast.FuncDecl:
		v.addFuncDecl(node)

	case *ast.GenDecl:
		switch node.Tok {
		case token.CONST:
			fallthrough
		case token.TYPE:
			fallthrough
		case token.VAR:
			v.inOrderNodes = append(v.inOrderNodes, node)
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

		v.inOrderNodes = append(
			v.inOrderNodes,
			&SelectorExpr{
				SelectorExpr: node,
				ImportName:   impName,
			},
		)
	}
	return v
}

func (v *DependencyVisitor) addImportSpec(node *ast.ImportSpec) {
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

	v.inOrderNodes = append(
		v.inOrderNodes,
		&ImportSpec{
			ImportSpec: node,
			IsAliased:  isAliased,
			Alias:      alias,
		},
	)
}

func (v *DependencyVisitor) addFuncDecl(node *ast.FuncDecl) {
	receiverName := ""
	qualifiedName := node.Name.String()

	if node.Recv != nil {
		// TODO: don't emit receiver functions/methods? we don't need them
		var typName string
		switch expr := node.Recv.List[0].Type.(type) {
		case *ast.Ident:
			typName = expr.String()
		case *ast.StarExpr:
			if expr.X == nil {
				// panic error, invalid receiver method
				panic("invalid receiver method")
			}
			ident, ok := expr.X.(*ast.Ident)
			if !ok {
				// panic error, invalid receiver method
				panic("invalid receiver method")
			}
			typName = ident.String()
		default:
			// panic error, invalid receiver method
			panic("invalid receiver method")
		}
		receiverName = typName
		qualifiedName = typName + "." + node.Name.String()
	}

	v.inOrderNodes = append(
		v.inOrderNodes,
		&FuncDecl{
			FuncDecl:      node,
			ReceiverName:  receiverName,
			QualifiedName: qualifiedName,
		},
	)
}
