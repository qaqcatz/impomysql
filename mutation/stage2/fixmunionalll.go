package stage2

import (
	"errors"
	"github.com/pingcap/tidb/parser/ast"
	_ "github.com/pingcap/tidb/parser/test_driver"
	"reflect"
)
// addFixMUnionAllL: FixMUnionAllL: *ast.SelectStmt: AfterSetOperator UNION ALL -> UNION
func (v *MutateVisitor) addFixMUnionAllL(in *ast.SelectStmt, flag int) {
	if in.AfterSetOperator != nil && *in.AfterSetOperator == ast.UnionAll {
		v.addCandidate(FixMUnionAllL, 0, in, flag)
	}
}

// doFixMUnionAllL: FixMUnionAllL, *ast.SelectStmt: AfterSetOperator UNION ALL -> UNION
func doFixMUnionAllL(rootNode ast.Node, in ast.Node) ([]byte, error) {
	switch in.(type) {
	case *ast.SelectStmt:
		sel := in.(*ast.SelectStmt)
		// check
		if sel.AfterSetOperator == nil || *sel.AfterSetOperator != ast.UnionAll {
			return nil, errors.New("doFixMUnionAllL: sel.AfterSetOperator == nil || *sel.AfterSetOperator != ast.UnionAll")
		}
		// mutate
		*sel.AfterSetOperator = ast.Union
		sql, err := restore(rootNode)
		if err != nil {
			return nil, errors.New("doFixMUnionAllL: " +  err.Error())
		}
		// recover
		*sel.AfterSetOperator = ast.UnionAll
		return sql, nil
	case nil:
		return nil, errors.New("doFixMUnionAllL: type error: nil")
	default:
		return nil, errors.New("doFixMUnionAllL: type error: " + reflect.TypeOf(in).String())
	}
}