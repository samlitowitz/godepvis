package cmd

import (
	"fmt"
	"github.com/samlitowitz/godepvis/internal"
	"github.com/samlitowitz/godepvis/internal/color"
	"github.com/samlitowitz/godepvis/internal/dot"
	"github.com/samlitowitz/godepvis/internal/modfile"
	"github.com/samlitowitz/godepvis/internal/primitives"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	PaletteFlag    = "palette"
	DotFlag        = "dot"
	PathFlag       = "path"
	ResolutionFlag = "resolution"
)

func Root() *cobra.Command {
	var resolution resolutionFlag
	rootCmd := &cobra.Command{
		Use:          "godepvis",
		Short:        "Go Dependency Visualizer",
		Long:         "Go Dependency Visualizer",
		SilenceUsage: true,
		RunE: func(self *cobra.Command, args []string) error {
			if len(args) != 0 {
				return self.Help()
			}

			paletteFile, err := self.Flags().GetString(PaletteFlag)
			if err != nil {
				return err
			}
			dotFile, err := self.Flags().GetString(DotFlag)
			if err != nil {
				return err
			}
			path, err := self.Flags().GetString(PathFlag)
			if err != nil {
				return nil
			}

			palette := color.DefaultPalette
			if paletteFile != "" {
				palette, err = color.GetPaletteFromFile(paletteFile)
				if err != nil {
					return err
				}
				if palette == nil {
					palette = color.DefaultPalette
				}
			}

			absPath, err := filepath.Abs(path)
			if err != nil {
				return err
			}

			goModFile, err := modfile.FindGoModFile(absPath)
			if err != nil {
				log.Fatal(fmt.Errorf("failed to find go.mod: %w", err))
			}

			modulePath, err := modfile.GetModulePath(goModFile)
			if err != nil {
				return err
			}
			moduleDir := filepath.Dir(goModFile)

			pkgs, err := primitives.BuildForModule(modulePath, moduleDir)
			if err != nil {
				log.Fatal(err)
			}

			output, err := dot.Marshal(
				modulePath,
				pkgs,
				dot.WithResolution(internal.Resolution(resolution.String())),
				dot.WithPalette(*palette),
			)
			if err != nil {
				log.Fatal(fmt.Errorf("marshal dependency graph: %w", err))
			}
			if dotFile == "" {
				_, err := os.Stdout.Write(output)
				if err != nil {
					return err
				}
				return nil
			}
			err = os.WriteFile(dotFile, output, 0644)
			if err != nil {
				return err
			}
			return nil
		},
	}

	rootCmd.Flags().String(PaletteFlag, "", "palette file")
	rootCmd.Flags().String(DotFlag, "", "DOT file to output")
	rootCmd.Flags().String(PathFlag, "", "files to process")
	rootCmd.Flags().Var(&resolution, ResolutionFlag, "resolution at which to visualize dependencies")

	err := rootCmd.MarkFlagRequired(DotFlag)
	if err != nil {
		panic(err)
	}

	return rootCmd
}

type resolutionFlag []byte

func (rf *resolutionFlag) String() string {
	return string(*rf)
}

func (rf *resolutionFlag) Set(v string) error {
	if internal.IsValidResolution(internal.Resolution(v)) {
		*rf = resolutionFlag(v)
		return nil
	}
	return fmt.Errorf("must be one of: %s", strings.Join(internal.ValidResolutions(), ", "))
}

func (rf *resolutionFlag) Type() string {
	return "resolutionFlag"
}
