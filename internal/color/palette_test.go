package color_test

import (
	"github.com/google/go-cmp/cmp"
	"github.com/samlitowitz/godepvis/internal/color"
	"github.com/samlitowitz/godepvis/internal/test"
	"gopkg.in/yaml.v3"
	"os"
	"runtime"
	"testing"
)

func TestGetPaletteFromFile(t *testing.T) {
	// REFURL: https://github.com/golang/go/blob/988b718f4130ab5b3ce5a5774e1a58e83c92a163/src/path/filepath/path_test.go#L600
	// -- START -- //
	if runtime.GOOS == "ios" {
		restore := test.Chtmpdir(t)
		defer restore()
	}

	tmpDir := t.TempDir()

	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal("finding working dir:", err)
	}
	if err = os.Chdir(tmpDir); err != nil {
		t.Fatal("entering temp dir:", err)
	}
	defer os.Chdir(origDir)
	// -- END -- //

	expectedPalette := color.InvertedDefaultPalette
	palettePath := tmpDir + string(os.PathSeparator) + "palette.yaml"
	writePalette(t, palettePath, expectedPalette)
	actualPalette, err := color.GetPaletteFromFile(palettePath)
	if err != nil {
		t.Fatalf("failed to load palette: %v", err)
	}
	compareHalfPalette(t, expectedPalette.Base, actualPalette.Base)
	compareHalfPalette(t, expectedPalette.Cycle, actualPalette.Cycle)
}

func compareHalfPalette(t *testing.T, expected, actual *color.HalfPalette) {
	if diff := cmp.Diff(expected.PackageName.Hex(), actual.PackageName.Hex()); diff != "" {
		t.Fatal(test.Mismatch("PackageName: ", diff))
	}
	if diff := cmp.Diff(expected.PackageBackground.Hex(), actual.PackageBackground.Hex()); diff != "" {
		t.Fatal(test.Mismatch("PackageBackground: ", diff))
	}
	if diff := cmp.Diff(expected.FileName.Hex(), actual.FileName.Hex()); diff != "" {
		t.Fatal(test.Mismatch("FileName: ", diff))
	}
	if diff := cmp.Diff(expected.FileBackground.Hex(), actual.FileBackground.Hex()); diff != "" {
		t.Fatal(test.Mismatch("FileBackground: ", diff))
	}
	if diff := cmp.Diff(expected.ImportArrow.Hex(), actual.ImportArrow.Hex()); diff != "" {
		t.Fatal(test.Mismatch("ImportArrow: ", diff))
	}
}

func writePalette(t *testing.T, filePath string, p *color.Palette) {
	fd, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("writePalette: %v", err)
		return
	}
	defer fd.Close()

	data, err := yaml.Marshal(p)
	if err != nil {
		t.Fatalf("writePalette: %v", err)
	}
	_, err = fd.Write(data)
	if err != nil {
		t.Fatalf("writePalette: %v", err)

	}
}
