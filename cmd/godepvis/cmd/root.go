package cmd

import (
	"fmt"
	internalAST "github.com/samlitowitz/godepvis/internal/ast"
	"github.com/samlitowitz/godepvis/internal/color"
	"github.com/samlitowitz/godepvis/internal/dot"
	"github.com/samlitowitz/godepvis/internal/modfile"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
)

const (
	ConfigFlag     = "config"
	DebugFlag      = "debug"
	DotFlag        = "dot"
	PathFlag       = "path"
	ResolutionFlag = "resolution"
)

func Root() *cobra.Command {
	var resolution resolutionFlag = resolutionFlagFile
	rootCmd := &cobra.Command{
		Use:          "godepvis",
		Short:        "Go Dependency Visualizer",
		Long:         "Go Dependency Visualizer",
		SilenceUsage: true,
		RunE: func(self *cobra.Command, args []string) error {
			if len(args) != 0 {
				return self.Help()
			}

			configFile, err := self.Flags().GetString(ConfigFlag)
			if err != nil {
				return err
			}
			debug, err := self.Flags().GetBool(DebugFlag)
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

			cfg := color.Default()
			if configFile != "" {
				cfg, err = color.FromYamlFile(configFile)
				if err != nil {
					return err
				}
				if cfg == nil {
					cfg = color.Default()
				}
			}

			if debug {
				cfg.Debug.SetOutput(os.Stdout)
			}

			absPath, err := filepath.Abs(path)
			if err != nil {
				return err
			}

			goModFile, err := modfile.FindGoModFile(absPath)
			if err != nil {
				log.Fatal(fmt.Errorf("failed to find go.mod: %w", err))
			}
			cfg.Debug.Printf("go.mod file: %s", goModFile)

			modulePath, err := modfile.GetModulePath(goModFile)
			if err != nil {
				return err
			}
			moduleRootDir := filepath.Dir(goModFile)
			cfg.Debug.Printf("Module Path: %s", modulePath)
			cfg.Debug.Printf("Module Root Directory: %s", moduleRootDir)

			switch resolution {
			case "file":
				cfg.Resolution = color.FileResolution
				if err != nil {
					return err
				}
			case "package":
				cfg.Resolution = color.PackageResolution
				if err != nil {
					return err
				}
			default:
				log.Fatal("resolution must be 'file' or 'package'")
			}

			primitivePkgs, err := internalAST.BuildPrimitivesForModule(modulePath, moduleRootDir)
			if err != nil {
				log.Fatal(fmt.Errorf("build primitives for module: %w", err))
			}

			output, err := dot.Marshal(cfg, modulePath, primitivePkgs)
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

	rootCmd.Flags().String(ConfigFlag, "", "configuration file")
	rootCmd.Flags().Bool(DebugFlag, false, "emit debug output")
	rootCmd.Flags().String(DotFlag, "", "DOT file to output")
	rootCmd.Flags().String(PathFlag, "", "files to process")
	rootCmd.Flags().Var(&resolution, ResolutionFlag, "resolution at which to visualize dependencies")

	err := rootCmd.MarkFlagRequired(DotFlag)
	if err != nil {
		panic(err)
	}

	return rootCmd
}

const (
	resolutionFlagFile    = "file"
	resolutionFlagPackage = "package"
)

type resolutionFlag string

func (rf *resolutionFlag) String() string {
	return string(*rf)
}

func (rf *resolutionFlag) Set(v string) error {
	switch v {
	case resolutionFlagFile, resolutionFlagPackage:
		*rf = resolutionFlag(v)
		return nil
	default:
		return fmt.Errorf(
			`must be "%s" or "%s"`,
			resolutionFlagFile,
			resolutionFlagPackage,
		)
	}
}

func (rf *resolutionFlag) Type() string {
	return "resolutionFlag"
}
