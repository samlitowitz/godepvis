package primitives

import (
	"fmt"
	"github.com/samlitowitz/godepvis/internal"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func BuildForModule(
	modulePath,
	moduleDir string,
) ([]*internal.Package, error) {
	var dirsToParse []string
	err := filepath.WalkDir(
		moduleDir,
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
			dirsToParse = append(dirsToParse, path)
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("get packages: %w", err)
	}

	depVis := NewDependencyVisitor()
	builder := NewPrimitiveBuilder(modulePath, moduleDir)
	fset := token.NewFileSet()
	// copied parser.ParseDir because we want to handle files and packages manually
	for _, dirToParse := range dirsToParse {
		list, err := os.ReadDir(dirToParse)
		if err != nil {
			return nil, err
		}

		pkgsSeen := map[string]bool{}
		for _, d := range list {
			if d.IsDir() ||
				strings.HasPrefix(d.Name(), ".") ||
				!strings.HasSuffix(d.Name(), ".go") ||
				strings.HasSuffix(d.Name(), "_test.go") {
				continue
			}

			filename := filepath.Join(dirToParse, d.Name())
			src, err := parser.ParseFile(fset, filename, nil, 0)
			if err != nil {
				return nil, fmt.Errorf("parse error: %s: %w", filename, err)
			}
			name := src.Name.Name
			if _, seen := pkgsSeen[name]; !seen {
				err = builder.AddNode(&Package{
					Package: &ast.Package{
						Name: name,
					},
					DirName: dirToParse,
				})
				if err != nil {
					return nil, fmt.Errorf("add package: %s: %w", filename, err)
				}
				pkgsSeen[name] = true
			}

			err = builder.AddNode(&File{
				File:    &ast.File{},
				AbsPath: filename,
				DirName: dirToParse,
			})
			if err != nil {
				return nil, fmt.Errorf("add file: %s: %w", filename, err)
			}
			depVis.Reset()
			ast.Walk(depVis, src)
			for _, node := range depVis.InOrderNodes() {
				err = builder.AddNode(node)
				if err != nil {
					return nil, fmt.Errorf("add node: %s: %w", filename, err)
				}
			}
		}
	}
	err = builder.MarkupImportCycles()
	if err != nil {
		return nil, err
	}
	return builder.Packages(), nil
}
