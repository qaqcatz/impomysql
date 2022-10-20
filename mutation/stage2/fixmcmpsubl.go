package stage2

import (
	"errors"
	_ "github.com/pingcap/tidb/parser/test_driver"
	"github.com/pingcap/tidb/parser/ast"
	"reflect"
)

// addFixMCmpSubL: FixMCmpSubL: *ast.CompareSubqueryExpr: ALL false -> true
func (v *MutateVisitor) addFixMCmpSubL(in *ast.CompareSubqueryExpr, flag int) {
	if in.All == false {
		v.addCandidate(FixMCmpSubL, 0, in, flag)
	}
}

// doFixMCmpSubL: FixMCmpSubL: *ast.CompareSubqueryExpr: ALL false -> true
func doFixMCmpSubL(rootNode ast.Node, in ast.Node) ([]byte, error) {
	switch in.(type) {
	case *ast.CompareSubqueryExpr:
		cmp := in.(*ast.CompareSubqueryExpr)
		// check
		if cmp.All != false {
			return nil, errors.New("doFixMCmpSubL: cmp.All != false")
		}
		// mutate
		cmp.All = true
		sql, err := restore(rootNode)
		if err != nil {
			return nil, errors.New("doFixMCmpSubL: " +  err.Error())
		}
		// recover
		cmp.All = false
		return sql, nil
	case nil:
		return nil, errors.New("doFixMCmpSubL: type error: nil")
	default:
		return nil, errors.New("doFixMCmpSubL: type error: " + reflect.TypeOf(in).String())
	}
}