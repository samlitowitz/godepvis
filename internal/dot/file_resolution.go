package dot

import (
	"bytes"
	"fmt"
	"github.com/samlitowitz/godepvis/internal"
	"github.com/samlitowitz/godepvis/internal/color"
)

const (
	fileResolutionEdgeDef = `
		"%s" -> "%s" [color="%s"];`
)

func writeNodeDefsForFileResolution(buf *bytes.Buffer, palette *color.Palette, pkgs []*internal.Package) {
	var err error
	clusterDefHeader := `
	subgraph "cluster_%s" {
		label="%s";
		style="filled";
		fontcolor="%s";
		fillcolor="%s";
`
	clusterDefFooter := `
	};
`
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
			clusterDefHeader,
			pkgNodeName(pkg),
			pkg.ModuleRelativePath(),
			pkgText.Hex(),
			pkgBackground.Hex(),
		)
		if err != nil {
			panic(err)
		}
		for _, file := range pkg.Files {
			if file.IsStub && !file.IsBlankImport {
				continue
			}
			if len(file.Decls) == 0 {
				continue
			}
			fileText := palette.Base.FileName
			fileBackground := palette.Base.FileBackground
			if file.InImportCycle {
				fileText = palette.Cycle.FileName
				fileBackground = palette.Cycle.FileBackground
			}
			_, err = fmt.Fprintf(
				buf,
				nodeDef,
				fileNodeName(file),
				file.FileName,
				fileText.Hex(),
				fileBackground.Hex(),
			)
			if err != nil {
				panic(err)
			}
		}
		buf.WriteString(clusterDefFooter)
	}
}

func writeRelationshipsForFileResolution(showMultipleReferences bool, buf *bytes.Buffer, palette *color.Palette, pkgs []*internal.Package) {
	writeEdgesFn := showOneReferencePerFileImportForFileResolution
	if showMultipleReferences {
		writeEdgesFn = showMultipleReferencesPerFileImportForFileResolution
	}

	for _, pkg := range pkgs {
		if pkg.IsStub {
			continue
		}
		for _, file := range pkg.Files {
			if file.IsStub {
				continue
			}
			writeEdgesFn(buf, palette, file)
		}
	}
}

func showMultipleReferencesPerFileImportForFileResolution(buf *bytes.Buffer, palette *color.Palette, file *internal.File) {
	var err error
	for _, imp := range file.Imports {
		if imp.Package == nil {
			continue
		}
		if imp.Package.IsStub {
			continue
		}
		for _, refTyp := range imp.ReferencedTypes {
			arrowColor := palette.Base.ImportArrow
			if _, ok := imp.ReferencedFilesInCycle[refTyp.File.UID()]; ok {
				arrowColor = palette.Cycle.ImportArrow
			}
			_, err = fmt.Fprintf(
				buf,
				fileResolutionEdgeDef,
				fileNodeName(file),
				fileNodeName(refTyp.File),
				arrowColor.Hex(),
			)
			if err != nil {
				panic(err)
			}
		}
	}
}

func showOneReferencePerFileImportForFileResolution(buf *bytes.Buffer, palette *color.Palette, file *internal.File) {
	var err error
	for _, imp := range file.Imports {
		if imp.Package == nil {
			continue
		}
		if imp.Package.IsStub {
			continue
		}

		for _, refTyp := range imp.ReferencedTypes {
			arrowColor := palette.Base.ImportArrow
			if _, ok := imp.ReferencedFilesInCycle[refTyp.File.UID()]; ok {
				arrowColor = palette.Cycle.ImportArrow
			}
			_, err = fmt.Fprintf(
				buf,
				fileResolutionEdgeDef,
				fileNodeName(file),
				fileNodeName(refTyp.File),
				arrowColor.Hex(),
			)
			if err != nil {
				panic(err)
			}
			break
		}
	}
}
