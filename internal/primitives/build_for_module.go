package primitives

import (
	"context"
	"errors"
	"fmt"
	"github.com/samlitowitz/godepvis/internal"
	"go/ast"
	"go/parser"
	"go/token"
	"golang.org/x/tools/go/packages"
)

func BuildForModule(
	modulePath,
	moduleDir string,
	buildFlags []string,
) ([]*internal.Package, error) {
	filesToParse, err := getFilesForModule(moduleDir, buildFlags)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(nil)

	depVis, nodeOut := NewDependencyVisitor()
	defer depVis.Close()

	builder := NewPrimitiveBuilder(modulePath, moduleDir)

	go func() {
		err := buildDependencyGraph(builder, nodeOut, ctx.Done())
		if err != nil {
			cancel(fmt.Errorf("build depedency graph: %w", err))
		}
	}()

	go func() {
		for _, file := range filesToParse {
			fset := token.NewFileSet()
			pkg, err := parser.ParseFile(fset, file, nil, 0)
			if err != nil {
				cancel(fmt.Errorf("parse files: %s: %w", file, err))
			}
			ast.Walk(depVis, pkg)
		}
		cancel(nil)
	}()

	<-ctx.Done()

	err = context.Cause(ctx)
	if !errors.Is(err, context.Canceled) {
		return nil, err
	}
	err = builder.MarkupImportCycles()
	if err != nil {
		return nil, err
	}
	return builder.Packages(), nil
}

func getFilesForModule(
	moduleDir string,
	buildFlags []string,
) ([]string, error) {
	pkgCfg := &packages.Config{
		Mode:       packages.LoadFiles | packages.LoadImports,
		Dir:        moduleDir,
		BuildFlags: buildFlags,
	}

	pkgs, err := packages.Load(pkgCfg, "./...")
	if err != nil {
		return nil, fmt.Errorf("failed to load module: %w", err)
	}

	var filesToParse []string
	iter := packages.Postorder(pkgs)
	for pkg := range iter {
		if len(pkg.Errors) > 0 {
			return nil, fmt.Errorf("errors loading module: %w", pkg.Errors[0])
		}
		filesToParse = append(filesToParse, pkg.CompiledGoFiles...)
	}
	return filesToParse, nil
}

func buildDependencyGraph(
	builder *PrimitiveBuilder,
	nodeOut <-chan ast.Node,
	done <-chan struct{},
) error {
	for {
		select {
		case node, ok := <-nodeOut:
			if !ok {
				return errors.New("failed to get next node")
			}
			err := builder.AddNode(node)
			if err != nil {
				return err
			}

		case <-done:
			return nil
		}
	}
}
