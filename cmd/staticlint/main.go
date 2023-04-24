package main

// Package main provides a set of analyzers:
//
// 1. ExitCheckAnalyzer checks for os.Exit call in main function
// 2. printf analyzer: checks consistency of Printf format strings and arguments.
// 3. shadow analyzer checks for shadowed variables
// 4. structtag analyzer check that struct field tags conform to reflect.StructTag.Get.
//    Also report certain struct tags (json, xml) used with unexported fields
//
// To run this analyzer execute staticlint with a list of go files or wildcards:
//   staticlint ./...
//

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
)

func main() {
	mychecks := []*analysis.Analyzer{
		ExitCheckAnalyzer,
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
	}
	classes := make(map[string]bool)
	for _, v := range staticcheck.Analyzers {
		if strings.HasPrefix(v.Analyzer.Name, "SA") {
			mychecks = append(mychecks, v.Analyzer)
		} else if !classes[v.Analyzer.Name[0:1]] {
			mychecks = append(mychecks, v.Analyzer)
			classes[v.Analyzer.Name[0:1]] = true
		}
	}
	multichecker.Main(mychecks...)
}

var ExitCheckAnalyzer = &analysis.Analyzer{
	Name: "exitanalyzer",
	Doc:  "check for os.Exit call in the main function of main package",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if file.Name.Name == "main" {
			ast.Inspect(file, func(node ast.Node) bool {
				funcDecl, ok := node.(*ast.FuncDecl)
				if !ok {
					return true
				}
				if funcDecl.Recv != nil || funcDecl.Name.Name != "main" { // we need only 'main' functions, not methods
					return false
				}

				checkFunc(pass, funcDecl)

				return false
			})
		}
	}
	return nil, nil
}

func checkFunc(pass *analysis.Pass, fn *ast.FuncDecl) {
	for _, stmt := range fn.Body.List {
		if expr, ok := stmt.(*ast.ExprStmt); ok {
			if call, ok := expr.X.(*ast.CallExpr); ok {
				if selExpr, ok := call.Fun.(*ast.SelectorExpr); ok {
					if x, ok1 := selExpr.X.(*ast.Ident); ok1 && x.Name == "os" && selExpr.Sel.Name == "Exit" {
						pass.Reportf(call.Pos(), "call of the 'os.Exit' in a main function of a main package")
					}
				}
			}
		}
	}
}
