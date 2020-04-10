package lib

import (
	"go/ast"
	"path/filepath"
	"strings"

	"github.com/go-toolsmith/astcopy"

	"github.com/pkg/errors"

	"golang.org/x/tools/go/loader"

	goofyast "github.com/mpppk/mustify/ast"
)

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

func pathHasPrefix(path, prefix string) bool {
	fileName := filepath.Base(path)
	return strings.HasPrefix(fileName, prefix)
}
