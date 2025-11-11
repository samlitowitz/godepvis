package dot

import (
	"bytes"
	"cmp"
	"fmt"
	"github.com/samlitowitz/godepvis/internal"
	"github.com/samlitowitz/godepvis/internal/color"
	"slices"
	"strings"
)

func Marshal(modulePath string, pkgs []*internal.Package, opts ...Option) ([]byte, error) {
	options := options{
		palette:    *color.DefaultPalette,
		resolution: internal.FileResolution,
	}
	for _, opt := range opts {
		opt.apply(&options)
	}

	slices.SortFunc(pkgs, pkgCmpFn)

	buf := &bytes.Buffer{}

	writeHeader(buf, modulePath)
	switch options.resolution {
	case internal.FileResolution:
		writeNodeDefsForFileResolution(buf, &options.palette, pkgs)
		writeRelationshipsForFileResolution(buf, &options.palette, pkgs)
	case internal.PackageResolution:
		writeNodeDefsForPackageResolution(buf, &options.palette, pkgs)
		writeRelationshipsForPackageResolution(buf, &options.palette, pkgs)
	}
	writeFooter(buf)

	return buf.Bytes(), nil
}

func pkgCmpFn(a, b *internal.Package) int {
	return cmp.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
}

func writeHeader(buf *bytes.Buffer, modulePath string) {
	_, err := fmt.Fprintf(
		buf,
		`digraph {
	labelloc="t";
	label="%s";
	rankdir="TB";
	node [shape="rect"];
`,
		modulePath,
	)
	if err != nil {
		panic(err)
	}
}

func writeFooter(buf *bytes.Buffer) {
	buf.WriteString(`
}
`,
	)
}

func pkgNodeName(pkg *internal.Package) string {
	return fmt.Sprintf(
		"pkg_%s",
		pkg.Name,
	)
}

func fileNodeName(file *internal.File) string {
	if file.Package == nil {
		return fmt.Sprintf(
			"file_%s",
			file.FileName,
		)
	}
	return fmt.Sprintf(
		"pkg_%s_file_%s",
		file.Package.Name,
		strings.TrimSuffix(file.FileName, ".go"),
	)
}
