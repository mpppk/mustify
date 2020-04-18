package lib

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io"

	"github.com/go-toolsmith/astcopy"

	ast2 "github.com/mpppk/mustify/ast"
)

func GenerateErrorWrappersFromReaderOrFile(fset *token.FileSet, filePath string, reader io.Reader) (*ast.File, bool, error) {
	if reader == nil {
		file, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
		if err != nil {
			return nil, false, err
		}
		return generateErrorWrappersFromFile(fset, file)
	}
	file, err := parser.ParseFile(fset, filePath, reader, parser.ParseComments)
	if err != nil {
		return nil, false, err
	}
	return generateErrorWrappersFromFile(fset, file)
}

func generateErrorWrappersFromFile(fset *token.FileSet, file *ast.File) (*ast.File, bool, error) {
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
	newFile.Comments = nil
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
