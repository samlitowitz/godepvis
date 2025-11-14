package primitives

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
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

func (decl FuncDecl) IsReceiver() bool {
	return decl.Recv != nil
}

type GenDecl struct {
	*ast.GenDecl

	FuncScopeName string
}

type SelectorExpr struct {
	*ast.SelectorExpr

	ImportName string
}

// DependencyVisitor is only for use against individual files
type DependencyVisitor struct {
	inOrderNodes []ast.Node

	fileImports map[string]struct{}

	curFuncScope             *funcScope
	funcScopeStack           funcScopeStack
	topLevelFuncLitNodeCount int

	logfFn func(...any)

	tmp []ast.Node
}

func NewDependencyVisitor() *DependencyVisitor {
	v := &DependencyVisitor{
		logfFn: func(...any) {},
	}

	return v
}

func (v *DependencyVisitor) Reset() {
	v.inOrderNodes = nil
	v.fileImports = nil
	v.funcScopeStack = nil

	v.tmp = nil
}

func (v *DependencyVisitor) InOrderNodes() []ast.Node {
	return v.inOrderNodes
}

func (v *DependencyVisitor) Visit(node ast.Node) ast.Visitor {
	//if v.curFuncScope != nil {
	//	v.tmp = append(v.tmp, node)
	//}

	if node == nil {
		v.exitNode()
	} else {
		v.enterNode()
	}

	switch node := node.(type) {

	case *ast.File:
		v.fileImports = make(map[string]struct{})

	case *ast.ImportSpec:
		v.addImportSpec(node)

	case *ast.FuncDecl:
		v.addFuncDecl(node)
		v.enterFuncDecl(node)

	case *ast.FuncLit:
		v.enterFuncLit(node)

	case *ast.GenDecl:
		switch node.Tok {
		case token.CONST:
			fallthrough
		case token.TYPE:
			fallthrough
		case token.VAR:
			genDecl := &GenDecl{GenDecl: node}
			if top := v.funcScopeStack.Top(); top != nil {
				genDecl.FuncScopeName = strings.Join(v.funcScopeStack.GetInOrderScopes(), ".")
			}

			v.inOrderNodes = append(v.inOrderNodes, genDecl)
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
		return
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

func (v *DependencyVisitor) enterFuncScope(name string) {
	v.curFuncScope = &funcScope{
		Name:         name,
		CurrentCount: 0,
	}
	v.funcScopeStack = v.funcScopeStack.Push(v.curFuncScope)
	v.logfFn("Enter Func: %s: %d : %d\n", v.curFuncScope.Name, v.curFuncScope.CurrentCount, len(v.tmp))
}

func (v *DependencyVisitor) enterFuncDecl(node *ast.FuncDecl) {
	v.tmp = append(v.tmp, node)
	v.enterFuncScope(node.Name.String())
}

func (v *DependencyVisitor) enterFuncLit(node *ast.FuncLit) {
	v.tmp = append(v.tmp, node)
	name := strconv.Itoa(v.topLevelFuncLitNodeCount)
	if v.curFuncScope != nil {
		name = strconv.Itoa(v.curFuncScope.TotalCount)
	}
	v.enterFuncScope(name)
	v.topLevelFuncLitNodeCount++
}

func (v *DependencyVisitor) exitFuncScope() {
	if v.curFuncScope != nil {
		v.logfFn("Exit Func: %s: %d : %d\n", v.curFuncScope.Name, v.curFuncScope.CurrentCount, len(v.tmp))
	}

	v.funcScopeStack, _ = v.funcScopeStack.Pop()
	v.curFuncScope = v.funcScopeStack.Top()
}

func (v *DependencyVisitor) enterNode() {
	if v.curFuncScope == nil {
		return
	}
	v.curFuncScope.CurrentCount++
	v.curFuncScope.TotalCount++
	v.logfFn("Enter Node: %d : %d\n", v.curFuncScope.CurrentCount, len(v.tmp))
}

func (v *DependencyVisitor) exitNode() {
	if v.curFuncScope == nil {
		return
	}
	v.logfFn("Exit Node: %d : %d\n", v.curFuncScope.CurrentCount, len(v.tmp))
	if v.curFuncScope.CurrentCount == 0 {
		v.exitFuncScope()
	}
	if v.curFuncScope == nil {
		return
	}
	v.curFuncScope.CurrentCount--
}

func getTypeName(typ ast.Expr) []string {
	switch expr := typ.(type) {
	case *ast.Ident:
		return []string{expr.String()}
	case *ast.StarExpr:
		if expr.X == nil {
			panic("invalid star expression")
		}
		return getTypeName(expr.X)

	case *ast.IndexExpr:
		return getTypeName(expr.Index)

	case *ast.IndexListExpr:
		typNames := make([]string, 0, len(expr.Indices))
		for _, lstExpr := range expr.Indices {
			typNames = append(typNames, getTypeName(lstExpr)...)
		}
		return typNames

	default:
		panic(fmt.Sprintf("unsupported type expression: %T", typ))
	}
}

type funcScope struct {
	Name string
	// TODO: document this and consider better nomenclature
	CurrentCount int
	TotalCount   int
}

type funcScopeStack []*funcScope

func (s funcScopeStack) Push(n *funcScope) funcScopeStack {
	return append(s, n)
}

func (s funcScopeStack) Pop() (funcScopeStack, *funcScope) {
	l := len(s)
	if l == 0 {
		return s, nil
	}
	return s[:l-1], s[l-1]
}

func (s funcScopeStack) Top() *funcScope {
	if len(s) == 0 {
		return nil
	}
	return s[len(s)-1]
}

func (s funcScopeStack) GetInOrderScopes() []string {
	scopes := make([]string, 0, len(s))
	for _, fnScope := range s {
		scopes = append(scopes, fnScope.Name)
	}
	return scopes
}
