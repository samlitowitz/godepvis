package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/samlitowitz/godepvis/internal/config"

	"github.com/samlitowitz/godepvis/internal/dot"

	internalAST "github.com/samlitowitz/godepvis/internal/ast"

	"github.com/samlitowitz/godepvis/internal/modfile"
)

func main() {
	var configFile, dotFile, path, resolution string
	var debug bool
	flag.StringVar(&configFile, "config", "", "Config file")
	flag.StringVar(&dotFile, "dot", "", "DOT file for output")
	flag.StringVar(&path, "path", "./", "Files to process")
	flag.StringVar(&resolution, "resolution", "file", "Resolution, 'file' or 'package'")
	flag.BoolVar(&debug, "debug", false, "Emit debug output")
	flag.Parse()

	var err error
	cfg := config.Default()
	if configFile != "" {
		cfg, err = config.FromYamlFile(configFile)
		if err != nil {
			log.Fatal(err)
		}
		if cfg == nil {
			cfg = config.Default()
		}
	}

	if debug {
		cfg.Debug.SetOutput(os.Stdout)
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}

	goModFile, err := modfile.FindGoModFile(absPath)
	if err != nil {
		log.Fatal(err)
	}
	cfg.Debug.Printf("go.mod file: %s", goModFile)

	modulePath, err := modfile.GetModulePath(goModFile)
	if err != nil {
		log.Fatal(err)
	}
	moduleRootDir := filepath.Dir(goModFile)
	cfg.Debug.Printf("Module Path: %s", modulePath)
	cfg.Debug.Printf("Module Root Directory: %s", moduleRootDir)

	switch resolution {
	case "file":
		cfg.Resolution = config.FileResolution
		if err != nil {
			log.Fatal(err)
		}
	case "package":
		cfg.Resolution = config.PackageResolution
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("resolution must be 'file' or 'package'")
	}

	packages, err := internalAST.BuildPrimitivesForModule(modulePath, moduleRootDir)
	if err != nil {
		log.Fatal(fmt.Errorf("build primitives for module: %w", err))
	}

	output, err := dot.Marshal(cfg, modulePath, packages)
	if err != nil {
		log.Fatal(fmt.Errorf("marshal dependency graph: %w", err))
	}
	if dotFile == "" {
		_, err := os.Stdout.Write(output)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	err = os.WriteFile(dotFile, output, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
