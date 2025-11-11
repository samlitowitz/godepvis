package ast

import (
	"context"
	"errors"
	"fmt"
	"github.com/samlitowitz/godepvis/internal"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
)

func BuildPrimitivesForModule(modulePath string, moduleRootDir string) ([]*internal.Package, error) {
	var filesToParse []string
	err := filepath.WalkDir(
		moduleRootDir,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				return nil
			}
			if strings.HasPrefix(d.Name(), ".") {
				return fs.SkipDir
			}
			if strings.HasPrefix(d.Name(), "_") {
				return fs.SkipDir
			}
			path, err = filepath.Abs(path)
			if err != nil {
				return err
			}
			filesToParse = append(filesToParse, path)
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("enumerate files to parse: %w", err)
	}

	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(nil)

	depVis, nodeOut := NewDependencyVisitor()
	defer depVis.Close()

	builder := NewPrimitiveBuilder(modulePath, moduleRootDir)

	go func() {
		err := buildDependencyGraph(builder, nodeOut, ctx.Done())
		if err != nil {
			cancel(fmt.Errorf("build depedency graph: %w", err))
		}
	}()

	go func() {
		for _, file := range filesToParse {
			fset := token.NewFileSet()
			pkgs, err := parser.ParseDir(fset, file, nil, 0)
			if err != nil {
				cancel(fmt.Errorf("parse files: %s: %w", file, err))
			}

			for _, pkg := range pkgs {
				ast.Walk(depVis, pkg)
			}
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
