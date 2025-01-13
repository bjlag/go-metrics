// Package exitcheck проверяет пакеты main на наличие прямого вызова os.Exit в функции main.
package exitcheck

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var ExitCheckAnalyzer = &analysis.Analyzer{
	Name: "exitcheck",
	Doc:  "check for use os.Exit in main package",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if file.Name.Name != "main" {
			continue
		}

		ast.Inspect(file, func(node ast.Node) bool {
			x, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}

			fun, ok := x.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			if ident, ok := fun.X.(*ast.Ident); ok {
				if ident.Name == "os" && fun.Sel.Name == "Exit" {
					pass.Reportf(x.Pos(), "os.Exit in main package")
				}
			}

			return true
		})
	}

	return nil, nil
}
