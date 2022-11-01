package stage2

import (
	"errors"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/test_driver"
	_ "github.com/pingcap/tidb/parser/test_driver"
	"reflect"
)

// addFixMWhere1U: FixMWhere1U, *ast.SelectStmt: WHERE xxx -> WHERE 1
func (v *MutateVisitor) addFixMWhere1U(in *ast.SelectStmt, flag int) {
	if in.Where != nil {
		v.addCandidate(FixMWhere1U, 1, in, flag)
	}
}

// doFixMWhere1U: FixMWhere1U, *ast.SelectStmt: WHERE xxx -> WHERE 1
func doFixMWhere1U(rootNode ast.Node, in ast.Node) ([]byte, error) {
	switch in.(type) {
	case *ast.SelectStmt:
		sel := in.(*ast.SelectStmt)
		// check
		if sel.Where == nil {
			return nil, errors.New("FixMWhere1U: sel.Where == nil")
		}
		// mutate
		old := sel.Where

		// WHERE xxx -> WHERE 1
		sel.Where = &test_driver.ValueExpr{
			Datum: test_driver.NewDatum(1),
		}

		sql, err := restore(rootNode)
		if err != nil {
			return nil, errors.New("FixMWhere1U: " +  err.Error())
		}
		// recover
		sel.Where = old
		return sql, nil
	case nil:
		return nil, errors.New("FixMWhere1U: type error: nil")
	default:
		return nil, errors.New("FixMWhere1U: type error: " + reflect.TypeOf(in).String())
	}
}