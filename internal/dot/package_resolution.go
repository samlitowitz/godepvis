package dot

import (
	"bytes"
	"fmt"

	"github.com/samlitowitz/goimportcycle/internal"
	"github.com/samlitowitz/goimportcycle/internal/config"
)

func writeNodeDefsForPackageResolution(buf *bytes.Buffer, cfg *config.Config, pkgs []*internal.Package) {
	nodeDef := `
	"%s" [label="%s", style="filled", fontcolor="%s", fillcolor="%s"];`

	for _, pkg := range pkgs {
		if pkg.IsStub {
			continue
		}
		if len(pkg.Files) == 0 {
			continue
		}
		pkgText := cfg.Palette.Base.PackageName
		pkgBackground := cfg.Palette.Base.PackageBackground
		if pkg.InImportCycle {
			pkgText = cfg.Palette.Cycle.PackageName
			pkgBackground = cfg.Palette.Cycle.PackageBackground
		}

		buf.WriteString(
			fmt.Sprintf(
				nodeDef,
				pkgNodeName(pkg),
				pkg.ModuleRelativePath(),
				pkgText.Hex(),
				pkgBackground.Hex(),
			),
		)
	}
}

func writeRelationshipsForPackageResolution(buf *bytes.Buffer, cfg *config.Config, pkgs []*internal.Package) {
	edgeDef := `
	"%s" -> "%s" [color="%s"];`

	// Track package-to-package edges to avoid duplicates
	type edge struct {
		from  string
		to    string
		color string
	}
	edges := make(map[string]edge)

	for _, pkg := range pkgs {
		if pkg.IsStub {
			continue
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

				fromNode := pkgNodeName(pkg)
				toNode := pkgNodeName(imp.Package)
				edgeKey := fromNode + "->" + toNode

				arrowColor := cfg.Palette.Base.ImportArrow
				if imp.InImportCycle {
					arrowColor = cfg.Palette.Cycle.ImportArrow
				}

				// Only add edge if not already seen, or if this one is in a cycle (higher priority)
				if existing, exists := edges[edgeKey]; !exists || (imp.InImportCycle && existing.color != arrowColor.Hex()) {
					edges[edgeKey] = edge{
						from:  fromNode,
						to:    toNode,
						color: arrowColor.Hex(),
					}
				}
			}
		}
	}

	// Write all unique edges
	for _, e := range edges {
		buf.WriteString(
			fmt.Sprintf(
				edgeDef,
				e.from,
				e.to,
				e.color,
			),
		)
	}
}
