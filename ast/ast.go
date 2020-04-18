package ast

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/go-toolsmith/astcopy"
)

func GenerateErrorFuncWrapper(orgFuncDecl *ast.FuncDecl) (*ast.FuncDecl, bool) {
	funcDecl := astcopy.FuncDecl(orgFuncDecl)
	if !IsErrorFunc(funcDecl) {
		return nil, false
	}

	results := getFuncDeclResults(funcDecl)
	funcDecl.Type.Results.List = funcDecl.Type.Results.List[:len(funcDecl.Type.Results.List)-1]

	wrappedCallExpr := generateCallExpr(extractRecvName(funcDecl), funcDecl.Name.Name, funcDecl.Type.Params.List)
	var lhs []string
	for _, result := range results {
		for _, name := range result.Names {
			lhs = append(lhs, name.Name)
		}
	}

	if len(lhs) == 0 {
		for i := 0; i < len(results); i++ {
			tempValueName := fmt.Sprintf("_v%d", i)
			lhs = append(lhs, tempValueName)
		}

		tempErrValueName := "_err"
		lhs[len(lhs)-1] = tempErrValueName
	}

	funcDecl.Body = &ast.BlockStmt{
		List: []ast.Stmt{
			generateAssignStmt(lhs, wrappedCallExpr),
			generatePanicIfErrorExistStmtAst(lhs[len(lhs)-1]),
			&ast.ReturnStmt{Results: identsToExprs(newIdents(lhs[:len(lhs)-1]))},
		},
	}
	addPrefixToFunc(funcDecl, "Must")

	funcDecl.Doc = orgFuncDecl.Doc
	return funcDecl, true
}

func identsToExprs(idents []*ast.Ident) (exprs []ast.Expr) {
	for _, ident := range idents {
		exprs = append(exprs, ast.Expr(ident))
	}
	return
}

func newIdents(identNames []string) (idents []*ast.Ident) {
	for _, identName := range identNames {
		idents = append(idents, &ast.Ident{
			Name: identName,
		})
	}
	return
}

func extractRecvName(funcDecl *ast.FuncDecl) string {
	if funcDecl.Recv == nil || len(funcDecl.Recv.List) <= 0 {
		return ""
	}
	names := funcDecl.Recv.List[0].Names
	if len(names) <= 0 {
		panic(fmt.Sprintf("unexpected recv names: %v from %v", names, funcDecl.Name.Name))
	}
	return names[0].Name
}

func ExtractImportDeclsFromDecls(decls []ast.Decl) (importDecls []*ast.GenDecl) {
	for _, decl := range decls {
		if importDecl, ok := declToImportDecl(decl); ok {
			importDecls = append(importDecls, importDecl)
		}
	}
	return
}

func declToImportDecl(decl ast.Decl) (*ast.GenDecl, bool) {
	if genDecl, ok := decl.(*ast.GenDecl); ok {
		if genDecl.Tok == token.IMPORT {
			return genDecl, true
		}
	}
	return nil, false
}

func ImportDeclsToDecls(importDecls []*ast.GenDecl) (decls []ast.Decl) {
	for _, decl := range importDecls {
		decls = append(decls, decl)
	}
	return
}
