package stage2

import (
	"errors"
	_ "github.com/pingcap/tidb/parser/test_driver"
	"github.com/pingcap/tidb/parser/ast"
	"reflect"
)

// addFixMCmpSubU: FixMCmpSubU, *ast.CompareSubqueryExpr: ALL true -> false
func (v *MutateVisitor) addFixMCmpSubU(in *ast.CompareSubqueryExpr, flag int) {
	if in.All == true {
		v.addCandidate(FixMCmpSubU, 1, in, flag)
	}
}

// doFixMCmpSubU: FixMCmpSubU, *ast.CompareSubqueryExpr: ALL true -> false
func doFixMCmpSubU(rootNode ast.Node, in ast.Node) ([]byte, error) {
	switch in.(type) {
	case *ast.CompareSubqueryExpr:
		cmp := in.(*ast.CompareSubqueryExpr)
		// check
		if cmp.All != true {
			return nil, errors.New("doFixMCmpSubU: cmp.All != true")
		}
		// mutate
		cmp.All = false
		sql, err := restore(rootNode)
		if err != nil {
			return nil, errors.New("doFixMCmpSubU: " +  err.Error())
		}
		// recover
		cmp.All = true
		return sql, nil
	case nil:
		return nil, errors.New("doFixMCmpSubU: type error: nil")
	default:
		return nil, errors.New("doFixMCmpSubU: type error: " + reflect.TypeOf(in).String())
	}
}

