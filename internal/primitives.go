package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Package struct {
	DirName string

	ModulePath string
	ModuleDir  string
	Name       string

	BlankImportFile *File
	Files           map[string]*File

	IsStub        bool
	InImportCycle bool
}

func (pkg Package) ImportPath() string {
	if pkg.Name == "main" {
		return ""
	}
	moduleRoot := pkg.ModuleDir
	if strings.LastIndex(moduleRoot, string(os.PathSeparator)) != len(pkg.ModuleDir)-1 {
		moduleRoot += string(os.PathSeparator)
	}
	if strings.HasPrefix(pkg.DirName, moduleRoot) {
		return fmt.Sprintf(
			"%s/%s",
			pkg.ModulePath,
			strings.TrimPrefix(
				pkg.DirName,
				moduleRoot,
			),
		)
	}
	if strings.HasPrefix(pkg.DirName, pkg.ModulePath) {
		return pkg.DirName
	}

	return pkg.Name
}

func (pkg Package) ModuleRelativePath() string {
	if strings.HasPrefix(pkg.DirName, pkg.ModuleDir) {
		path := strings.TrimPrefix(
			pkg.DirName,
			pkg.ModuleDir,
		)
		path = strings.TrimPrefix(path, string(filepath.Separator))
		if pkg.Name != "main" {
			return path
		}
		if path == "" {
			return pkg.Name
		}
		return path + ":" + pkg.Name
	}
	if strings.HasPrefix(pkg.DirName, pkg.ModulePath) {
		path := strings.TrimPrefix(
			pkg.DirName,
			pkg.ModulePath,
		)
		path = strings.TrimPrefix(path, string(filepath.Separator))
		if pkg.Name != "main" {
			return path
		}
		if path == "" {
			return pkg.Name
		}
		return path + ":" + pkg.Name
	}
	return pkg.Name
}

func (pkg Package) UID() string {
	uid := pkg.ImportPath()
	if uid != "" {
		return uid
	}
	return pkg.DirName
}

func (pkg Package) HasBlankImports() bool {
	if pkg.BlankImportFile == nil {
		return false
	}
	_, ok := pkg.Files[pkg.BlankImportFile.UID()]
	return ok
}

type File struct {
	Package *Package

	FileName string
	AbsPath  string

	Imports map[string]*Import
	Decls   map[string]*Decl

	IsStub        bool
	IsBlankImport bool
	InImportCycle bool
}

func (f File) HasDecl(decl *Decl) bool {
	for _, fDecl := range f.Decls {
		if decl.UID() != fDecl.UID() {
			continue
		}
		return true
	}
	return false
}

func (f File) ReferencedFiles() []*File {
	alreadyReferenced := make(map[string]struct{})
	referencedFiles := make([]*File, 0, len(f.Imports))

	for _, imp := range f.Imports {
		for _, typ := range imp.ReferencedTypes {
			if _, ok := alreadyReferenced[typ.File.AbsPath]; ok {
				continue
			}
			alreadyReferenced[typ.File.AbsPath] = struct{}{}
			referencedFiles = append(referencedFiles, typ.File)
		}
	}
	return referencedFiles

}

func (f File) UID() string {
	return f.AbsPath
}

type Decl struct {
	File *File

	Name string
}

func (decl Decl) UID() string {
	return decl.QualifiedName()
}

func (decl Decl) QualifiedName() string {
	return decl.Name
}

type Import struct {
	Package *Package

	Name  string
	Alias string
	Path  string

	IsAliased bool
	IsBlank   bool

	ReferencedTypes map[string]*Decl

	InImportCycle          bool
	ReferencedFilesInCycle map[string]*File
}

func (i Import) UID() string {
	if i.IsBlank {
		return i.Alias + i.Name
	}
	if i.IsAliased {
		return i.Alias
	}
	return i.Name
}
