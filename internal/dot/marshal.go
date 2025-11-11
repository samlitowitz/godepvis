package dot

import (
	"bytes"
	"cmp"
	"fmt"
	"github.com/samlitowitz/godepvis/internal/color"
	"slices"
	"strings"

	"github.com/samlitowitz/godepvis/internal"
)

func Marshal(cfg *color.Config, modulePath string, pkgs []*internal.Package) ([]byte, error) {
	slices.SortFunc(pkgs, pkgCmpFn)

	buf := &bytes.Buffer{}

	writeHeader(buf, modulePath)
	switch cfg.Resolution {
	case color.FileResolution:
		writeNodeDefsForFileResolution(buf, cfg, pkgs)
		writeRelationshipsForFileResolution(buf, cfg, pkgs)
	case color.PackageResolution:
		writeNodeDefsForPackageResolution(buf, cfg, pkgs)
		writeRelationshipsForPackageResolution(buf, cfg, pkgs)
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
