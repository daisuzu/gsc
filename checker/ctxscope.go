package checker

import (
	"go/ast"
	"go/types"
	"strings"

	"honnef.co/go/tools/lint"
)

func (c *checker) isContext(s string) bool {
	for _, v := range c.contextNames {
		if strings.HasSuffix(s, v) {
			return true
		}
	}
	return false
}

func (c *checker) CheckCtxScope(j *lint.Job) {
	fn := func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		if len(call.Args) == 0 {
			return true
		}
		sig, ok := j.Program.Info.TypeOf(call.Fun).(*types.Signature)
		if !ok {
			return true
		}
		if sig.Params().Len() == 0 {
			return true
		}
		if !c.isContext(types.TypeString(sig.Params().At(0).Type(), nil)) {
			return true
		}

		// Check ctx declaration exists in the function that contains current scope.
		scope := j.NodePackage(node).Pkg.Scope().Innermost(call.Pos())
		for !strings.HasPrefix(scope.String(), "function scope") {
			scope = scope.Parent()
		}
		if obj := scope.Lookup(types.ExprString(call.Args[0])); obj != nil {
			return true
		}

		// Allow deferred closure without arguments.
		var dc bool
		ast.Inspect(j.File(scope), func(n ast.Node) bool {
			d, ok := n.(*ast.DeferStmt)
			if !ok {
				return true
			}
			f, ok := d.Call.Fun.(*ast.FuncLit)
			if !ok {
				return true
			}
			if f.Body == nil {
				return true
			}
			if f.Body.Pos() != scope.Pos() {
				return true
			}
			dc = len(f.Type.Params.List) == 0
			return false
		})
		if dc {
			return true
		}

		j.Errorf(call.Args[0], "passing outer scope context %q to %s()", types.ExprString(call.Args[0]), types.ExprString(call.Fun))
		return true
	}
	for _, f := range j.Program.Files {
		ast.Inspect(f, fn)
	}
}
