package lib

import (
	"go/ast"
	"go/format"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-toolsmith/astcopy"

	"github.com/pkg/errors"

	"golang.org/x/tools/go/loader"

	goofyast "github.com/mpppk/mustify/ast"
)

func GenerateErrorWrappersFromFile(filePath string) (*ast.File, []ast.Decl, error) {
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to get abs file path")
	}

	prog, err := goofyast.NewProgram(filePath)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to load program file")
	}

	pkg, file, ok := findPkgAndFileFromProgram(prog, absFilePath)
	if !ok {
		return nil, nil, errors.New("file not found: " + filePath)
	}

	var newDecls []ast.Decl
	importDecls := goofyast.ExtractImportDeclsFromDecls(file.Decls)
	newDecls = append(newDecls, goofyast.ImportDeclsToDecls(importDecls)...)
	exportedFuncDecls := extractExportedFuncDeclsFromDecls(file.Decls)
	errorWrappers := funcDeclsToErrorFuncWrappers(exportedFuncDecls, pkg)
	newDecls = append(newDecls, errorWrappers...)
	return file, newDecls, nil
}

func GenerateErrorWrappersFromPackage(filePath, pkgName, ignorePrefix string) (map[string]*ast.File, error) {
	prog, err := goofyast.NewProgram(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load program file")
	}

	pkg := prog.Package(pkgName)
	m := map[string]*ast.File{}
	for _, file := range pkg.Files {
		filePath := prog.Fset.File(file.Pos()).Name()
		if pathHasPrefix(filePath, ignorePrefix) {
			continue
		}

		newFile := astcopy.File(file)
		var newDecls []ast.Decl
		importDecls := goofyast.ExtractImportDeclsFromDecls(newFile.Decls)
		newDecls = append(newDecls, goofyast.ImportDeclsToDecls(importDecls)...)
		exportedFuncDecls := extractExportedFuncDeclsFromDecls(newFile.Decls)
		errorWrappers := funcDeclsToErrorFuncWrappers(exportedFuncDecls, pkg)
		if len(errorWrappers) <= 0 {
			continue
		}
		newDecls = append(newDecls, errorWrappers...)
		newFile.Decls = newDecls
		m[filePath] = newFile
	}
	return m, nil
}

func extractExportedFuncDeclsFromDecls(decls []ast.Decl) (funcDecls []*ast.FuncDecl) {
	for _, decl := range decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			if ast.IsExported(funcDecl.Name.Name) {
				funcDecls = append(funcDecls, funcDecl)
			}
		}
	}
	return
}

func funcDeclsToErrorFuncWrappers(funcDecls []*ast.FuncDecl, pkg *loader.PackageInfo) (newDecls []ast.Decl) {
	for _, funcDecl := range funcDecls {
		//newDecl, ok := goofyast.ConvertErrorFuncToMustFunc(prog, pkg, funcDecl)
		if newDecl, ok := goofyast.GenerateErrorFuncWrapper(pkg, funcDecl); ok {
			newDecls = append(newDecls, newDecl)
		}
	}
	return
}

func findPkgAndFileFromProgram(prog *loader.Program, targetAbsFilePath string) (*loader.PackageInfo, *ast.File, bool) {
	for _, pkg := range prog.Created {
		for _, file := range pkg.Files {
			currentFilePath := prog.Fset.File(file.Pos()).Name()
			if absCurrentFilePath, err := filepath.Abs(currentFilePath); err == nil {
				if targetAbsFilePath == absCurrentFilePath {
					return pkg, file, true
				}
			}

		}
	}
	return nil, nil, false
}

func pathHasPrefix(path, prefix string) bool {
	fileName := filepath.Base(path)
	return strings.HasPrefix(fileName, prefix)
}

func WriteAstFile(filePath string, file *ast.File) error {
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return errors.Wrap(err, "failed to get abs file path: "+absFilePath)
	}

	f, err := os.Create(absFilePath)
	if err != nil {
		return errors.Wrap(err, "failed to create file: "+absFilePath)
	}

	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()
	if err := format.Node(f, token.NewFileSet(), file); err != nil {
		return errors.Wrap(err, "failed to write ast file to  "+absFilePath)
	}
	return nil
}
