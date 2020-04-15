package lib

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"

	"github.com/go-toolsmith/astcopy"

	"github.com/pkg/errors"

	ast2 "github.com/mpppk/mustify/ast"
)

func GenerateErrorWrappersFromFilePath(filePath string) (*ast.File, bool, error) {
	dirPath := filepath.Dir(filePath)
	if !strings.HasPrefix(dirPath, ".") && !strings.HasSuffix(dirPath, "/") {
		dirPath = "./" + dirPath
	}
	file, err := parser.ParseFile(token.NewFileSet(), filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, false, err
	}
	return generateErrorWrappersFromFile(file)
}

func NewPackage(path, packageName string) (*packages.Package, error) {
	config := &packages.Config{
		Mode: packages.NeedCompiledGoFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo | packages.LoadAllSyntax,
	}
	pkgs, err := packages.Load(config, path)
	if err != nil {
		return nil, err
	}
	if packages.PrintErrors(pkgs) > 0 {
		return nil, errors.New("error occurred in NewProgramFromPackages")
	}
	for _, pkg := range pkgs {
		if pkg.Name == packageName {
			return pkg, nil
		}
	}
	return nil, errors.New("pkg not found: " + packageName)
}

func generateErrorWrappersFromFile(file *ast.File) (*ast.File, bool, error) {
	newFile := astcopy.File(file)
	var newDecls []ast.Decl
	importDecls := ast2.ExtractImportDeclsFromDecls(newFile.Decls)
	newDecls = append(newDecls, ast2.ImportDeclsToDecls(importDecls)...)
	exportedFuncDecls := extractExportedFuncDeclsFromDecls(newFile.Decls)
	errorWrappers := funcDeclsToErrorFuncWrappers(exportedFuncDecls)
	if len(errorWrappers) <= 0 {
		return nil, false, nil
	}
	newDecls = append(newDecls, errorWrappers...)
	newFile.Decls = newDecls
	return newFile, true, nil
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

func funcDeclsToErrorFuncWrappers(funcDecls []*ast.FuncDecl) (newDecls []ast.Decl) {
	for _, funcDecl := range funcDecls {
		if newDecl, ok := ast2.GenerateErrorFuncWrapper(funcDecl); ok {
			newDecls = append(newDecls, newDecl)
		}
	}
	return
}
