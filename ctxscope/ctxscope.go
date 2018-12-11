package ctxscope

import (
	"fmt"
	"go/ast"
	"go/types"
	"os"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:             "ctxscope",
	Doc:              "report passing outer scope context",
	Requires:         []*analysis.Analyzer{inspect.Analyzer},
	RunDespiteErrors: true,
	Run:              run,
}

type contexts []string

func (v *contexts) Set(s string) error {
	*v = append(*v, s)
	return nil
}

func (v *contexts) Get() interface{} { return *v }
func (v *contexts) String() string   { return "<contexts>" }

var (
	exitNonZero bool
	ctxs        contexts
)

func init() {
	Analyzer.Flags.BoolVar(&exitNonZero, "exit-non-zero", true, "exit non-zero if any problems were found")
	Analyzer.Flags.Var(&ctxs, "target-context", "additional target context types other than the standard library's context")
}

func isContext(s string) bool {
	for _, v := range append([]string{"context.Context"}, ctxs...) {
		if strings.HasSuffix(s, v) {
			return true
		}
	}
	return false
}

func run(pass *analysis.Pass) (interface{}, error) {
	if !exitNonZero {
		pass.Report = func(diag analysis.Diagnostic) {
			posn := pass.Fset.Position(diag.Pos)
			fmt.Fprintf(os.Stderr, "%s: %s\n", posn, diag.Message)
		}
	}

	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}
	inspect.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)
		if len(call.Args) == 0 {
			return
		}

		sig, ok := pass.TypesInfo.TypeOf(call.Fun).(*types.Signature)
		if !ok {
			return
		}
		if sig.Params().Len() == 0 {
			return
		}
		if !isContext(types.TypeString(sig.Params().At(0).Type(), nil)) {
			return
		}

		// Check ctx declaration exists in the function that contains current scope.
		scope := pass.Pkg.Scope().Innermost(call.Pos())
		for !strings.HasPrefix(scope.String(), "function scope") {
			scope = scope.Parent()
		}
		if obj := scope.Lookup(types.ExprString(call.Args[0])); obj != nil {
			return
		}

		// Allow deferred closure without arguments.
		var dc bool
		ast.Inspect(getFile(pass, scope), func(n ast.Node) bool {
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
			return
		}

		if allowedCtx(call.Args[0]) {
			return
		}

		pass.Reportf(call.Args[0].Pos(), "passing outer scope context %q to %s()", types.ExprString(call.Args[0]), types.ExprString(call.Fun))
	})

	return nil, nil
}

// allowedCtx checks whether arg which returns context is whitelisted.
//   - "context" or "google.golang.org/appengine" package
//   - "net/http".Request.Context
//   - func that returns above context
func allowedCtx(arg ast.Expr) bool {
	if c, ok := arg.(*ast.CallExpr); ok {
		switch t := c.Fun.(type) {
		case *ast.SelectorExpr:
			if i, ok := t.X.(*ast.Ident); ok {
				if i.Obj == nil {
					if i.Name == "context" || i.Name == "appengine" {
						return true
					}
				} else {
					if f, ok := i.Obj.Decl.(*ast.Field); ok {
						if types.ExprString(f.Type) == "*http.Request" {
							return true
						}
					}
				}
			}
		case *ast.Ident:
			return allowedCtx(c.Args[0])
		}
	}
	return false
}

func getFile(pass *analysis.Pass, scope *types.Scope) *ast.File {
	name := pass.Fset.Position(scope.Pos()).Filename
	for _, v := range pass.Files {
		if name == pass.Fset.Position(v.Pos()).Filename {
			return v
		}
	}
	return nil
}
