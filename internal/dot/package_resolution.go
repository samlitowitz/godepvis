package dot

import (
	"bytes"
	"fmt"
	"github.com/samlitowitz/godepvis/internal"
	"github.com/samlitowitz/godepvis/internal/color"
)

func writeNodeDefsForPackageResolution(buf *bytes.Buffer, palette *color.Palette, pkgs []*internal.Package) {
	var err error
	nodeDef := `
	"%s" [label="%s", style="filled", fontcolor="%s", fillcolor="%s"];`

	for _, pkg := range pkgs {
		if pkg.IsStub {
			continue
		}
		if len(pkg.Files) == 0 {
			continue
		}
		pkgText := palette.Base.PackageName
		pkgBackground := palette.Base.PackageBackground
		if pkg.InImportCycle {
			pkgText = palette.Cycle.PackageName
			pkgBackground = palette.Cycle.PackageBackground
		}

		_, err = fmt.Fprintf(
			buf,
			nodeDef,
			pkgNodeName(pkg),
			pkg.ModuleRelativePath(),
			pkgText.Hex(),
			pkgBackground.Hex(),
		)
		if err != nil {
			panic(err)
		}
	}
}

func writeRelationshipsForPackageResolution(buf *bytes.Buffer, palette *color.Palette, pkgs []*internal.Package) {
	var err error
	edgeDef := `
	"%s" -> "%s" [color="%s"];`

	pkgRelationships := make(map[string]map[string]bool)
	for _, pkg := range pkgs {
		if pkg.IsStub {
			continue
		}
		pkgName := pkgNodeName(pkg)
		if _, ok := pkgRelationships[pkgName]; !ok {
			pkgRelationships[pkgName] = make(map[string]bool)
		}
		for _, file := range pkg.Files {
			if file.IsStub {
				continue
			}
			for _, imp := range file.Imports {
				if imp.Package == nil {
					continue
				}
				if imp.Package.IsStub {
					continue
				}
				impPkgName := pkgNodeName(imp.Package)
				// don't write a relationship multiple times
				// this could happen when multiple files in a package import the same package
				if _, ok := pkgRelationships[pkgName][impPkgName]; ok {
					continue
				}
				pkgRelationships[pkgName][impPkgName] = true

				arrowColor := palette.Base.ImportArrow
				if imp.InImportCycle {
					arrowColor = palette.Cycle.ImportArrow
				}
				_, err = fmt.Fprintf(
					buf,
					edgeDef,
					pkgName,
					impPkgName,
					arrowColor.Hex(),
				)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}
