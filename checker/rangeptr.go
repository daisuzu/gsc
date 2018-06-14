package checker

import (
	"go/ast"
	"go/token"
	"strings"

	"honnef.co/go/tools/lint"
)

func (c *checker) CheckRangePtr(j *lint.Job) {
	fn := func(node ast.Node) bool {
		assign, ok := node.(*ast.AssignStmt)
		if !ok {
			return true
		}

		scope := j.NodePackage(node).Pkg.Scope().Innermost(assign.Pos()).Parent()
		if !strings.HasPrefix(scope.String(), "for scope") {
			return true
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
				j.Errorf(ident, "using pointer to the loop iteration variable %q", obj.Name())
			}
		}

		return true
	}
	for _, f := range j.Program.Files {
		ast.Inspect(f, fn)
	}
}
