package color_test

import (
	"github.com/samlitowitz/godepvis/internal/color"
	"gopkg.in/yaml.v3"
	"os"
	"runtime"
	"testing"
)

func TestGetPaletteFromFile(t *testing.T) {
	// REFURL: https://github.com/golang/go/blob/988b718f4130ab5b3ce5a5774e1a58e83c92a163/src/path/filepath/path_test.go#L600
	// -- START -- //
	if runtime.GOOS == "ios" {
		restore := chtmpdir(t)
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
	if expected.PackageName.Hex() != actual.PackageName.Hex() {
		t.Errorf("PackageName: expected %s got %s", expected.PackageName.Hex(), actual.PackageName.Hex())
	}
	if expected.PackageBackground.Hex() != actual.PackageBackground.Hex() {
		t.Errorf("PackageBackground: expected %s got %s", expected.PackageBackground.Hex(), actual.PackageBackground.Hex())
	}
	if expected.FileName.Hex() != actual.FileName.Hex() {
		t.Errorf("FileName: expected %s got %s", expected.FileName.Hex(), actual.FileName.Hex())
	}
	if expected.FileBackground.Hex() != actual.FileBackground.Hex() {
		t.Errorf("FileBackground: expected %s got %s", expected.FileBackground.Hex(), actual.FileBackground.Hex())
	}
	if expected.ImportArrow.Hex() != actual.ImportArrow.Hex() {
		t.Errorf("ImportArrow: expected %s got %s", expected.ImportArrow.Hex(), actual.ImportArrow.Hex())
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

// REFURL: https://github.com/golang/go/blob/988b718f4130ab5b3ce5a5774e1a58e83c92a163/src/path/filepath/path_test.go#L553
func chtmpdir(t *testing.T) (restore func()) {
	oldwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("chtmpdir: %v", err)
	}
	d, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("chtmpdir: %v", err)
	}
	if err := os.Chdir(d); err != nil {
		t.Fatalf("chtmpdir: %v", err)
	}
	return func() {
		if err := os.Chdir(oldwd); err != nil {
			t.Fatalf("chtmpdir: %v", err)
		}
		_ = os.RemoveAll(d)
	}
}
