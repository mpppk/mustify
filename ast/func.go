package ast

import (
	"go/ast"
	"strings"
)

func IsErrorFunc(funcDecl *ast.FuncDecl) bool {
	lastResultIdent, ok := extractFuncLastResultIdent(funcDecl)
	if !ok {
		return false
	}
	return lastResultIdent.Name == "error"
}

func getFuncDeclResults(funcDecl *ast.FuncDecl) (newResults []*ast.Field) {
	results := funcDecl.Type.Results.List
	for _, result := range results {
		newResults = append(newResults, result)
	}
	return
}

func addPrefixToFunc(funcDecl *ast.FuncDecl, prefix string) {
	funcNameRunes := []rune(funcDecl.Name.Name)
	funcDecl.Name.Name = prefix + strings.ToUpper(string(funcNameRunes[0])) + string(funcNameRunes[1:])
}
