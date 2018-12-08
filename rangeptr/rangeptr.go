package rangeptr

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:             "rangeptr",
	Doc:              "report using pointer to the loop iteration variable",
	Requires:         []*analysis.Analyzer{inspect.Analyzer},
	RunDespiteErrors: true,
	Run:              run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.AssignStmt)(nil),
	}
	inspect.Preorder(nodeFilter, func(n ast.Node) {
		assign := n.(*ast.AssignStmt)

		scope := pass.Pkg.Scope().Innermost(assign.Pos()).Parent()
		if !strings.HasPrefix(scope.String(), "for scope") {
			return
		}

		var exprs []*ast.UnaryExpr
		for _, v := range assign.Rhs {
			switch t := v.(type) {
			case *ast.UnaryExpr:
				exprs = append(exprs, t)
			case *ast.CallExpr:
				for _, vv := range t.Args {
					if ue, ok := vv.(*ast.UnaryExpr); ok {
						exprs = append(exprs, ue)
					}
				}
			}
		}

		for _, v := range exprs {
			if v.Op != token.AND {
				continue
			}

			ident, ok := v.X.(*ast.Ident)
			if !ok {
				continue
			}

			if obj := scope.Lookup(ident.Name); obj != nil {
				pass.Reportf(ident.Pos(), "using pointer to the loop iteration variable %q", obj.Name())
			}
		}
	})

	return nil, nil
}
