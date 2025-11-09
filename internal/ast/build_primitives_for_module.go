package ast

import (
	"context"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"log"
	"path/filepath"
	"strings"

	"github.com/samlitowitz/goimportcycle/internal"
)

func BuildPrimitivesForModule(modulePath string, moduleRootDir string) ([]*internal.Package, error) {
	builder := NewPrimitiveBuilder(modulePath, moduleRootDir)
	ctx, cancel := context.WithCancel(context.Background())
	errChan := make(chan error)

	go func() {
		for {
			select {
			case err := <-errChan:
				cancel()
				log.Fatal(err)
			case <-ctx.Done():
				return
			}
		}
	}()

	dirOut := walkDirectories(moduleRootDir, errChan)
	nodeOut := parseFiles(dirOut, errChan, ctx.Done())
	err := detectInputCycles(builder, cancel, nodeOut, errChan, ctx.Done())
	close(errChan)
	if err != nil {
		log.Fatal(err)
	}
	return builder.Packages(), nil
}

func walkDirectories(path string, errChan chan<- error) <-chan string {
	dirOut := make(chan string)

	go func() {
		defer close(dirOut)
		err := filepath.WalkDir(
			path,
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
				dirOut <- path
				return nil
			},
		)
		if err != nil {
			errChan <- err
		}
	}()

	return dirOut
}

func parseFiles(
	dirOut <-chan string,
	errChan chan<- error,
	done <-chan struct{},
) <-chan ast.Node {
	depVis, nodeOut := NewDependencyVisitor()

	go func() {
		for {
			select {
			case dirPath, ok := <-dirOut:
				if !ok {
					depVis.Close()
					return
				}
				fset := token.NewFileSet()
				pkgs, err := parser.ParseDir(fset, dirPath, nil, 0)
				if err != nil {
					errChan <- err
				}

				for _, pkg := range pkgs {
					ast.Walk(depVis, pkg)
				}

			case <-done:
				depVis.Close()
				return
			}
		}
	}()
	return nodeOut
}

func detectInputCycles(
	builder *PrimitiveBuilder,
	cancel context.CancelFunc,
	nodeOut <-chan ast.Node,
	errChan chan<- error,
	done <-chan struct{},
) error {
	go func() {
		for {
			select {
			case node, ok := <-nodeOut:
				if !ok {
					cancel()
					return
				}
				err := builder.AddNode(node)
				if err != nil {
					errChan <- err
					return
				}

			case <-done:
				return
			}
		}
	}()
	<-done
	return builder.MarkupImportCycles()
}
