package stage2

import (
	"github.com/pkg/errors"
	_ "github.com/pingcap/tidb/parser/test_driver"
	"github.com/pingcap/tidb/parser/ast"
	"reflect"
)

// addFixMDistinctU: FixMDistinctU, *ast.SelectStmt: Distinct true -> false
func (v *MutateVisitor) addFixMDistinctU(in *ast.SelectStmt, flag int) {
	if in.Distinct == true {
		v.addCandidate(FixMDistinctU, 1, in, flag)
	}
}

// doFixMDistinctU: FixMDistinctU, *ast.SelectStmt: Distinct true -> false
func doFixMDistinctU(rootNode ast.Node, in ast.Node) ([]byte, error) {
	switch in.(type) {
	case *ast.SelectStmt:
		sel := in.(*ast.SelectStmt)
		// check
		if sel.Distinct != true {
			return nil, errors.New("[doFixMDistinctU]in.Distinct != true")
		}
		// mutate
		sel.Distinct = false
		sql, err := restore(rootNode)
		if err != nil {
			return nil, errors.Wrap(err, "[doFixMDistinctU]restore error")
		}
		// recover
		sel.Distinct = true
		return sql, nil
	case nil:
		return nil, errors.New("[doFixMDistinctU]type nil")
	default:
		return nil, errors.New("[doFixMDistinctU]type default " + reflect.TypeOf(in).String())
	}
}