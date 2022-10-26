package stage2

import (
	"errors"
	_ "github.com/pingcap/tidb/parser/test_driver"
	"github.com/pingcap/tidb/parser/ast"
	"reflect"
)

// addFixMDistinctL: FixMDistinctL: *ast.SelectStmt: Distinct false -> true
func (v *MutateVisitor) addFixMDistinctL(in *ast.SelectStmt, flag int) {
	// ERROR 3065 (HY000): Expression #1 of ORDER BY clause is not in SELECT list,
	// references column xxx which is not in SELECT list; this is incompatible with DISTINCT
	// order by + distinct may error
	if in.Distinct == false && in.OrderBy == nil {
		v.addCandidate(FixMDistinctL, 0, in, flag)
	}
}

// doFixMDistinctL: FixMDistinctL: *ast.SelectStmt: Distinct false -> true
func doFixMDistinctL(rootNode ast.Node, in ast.Node) ([]byte, error) {
	switch in.(type) {
	case *ast.SelectStmt:
		sel := in.(*ast.SelectStmt)
		// check
		if sel.Distinct != false {
			return nil, errors.New("doFixMDistinctL: in.Distinct != false")
		}
		// mutate
		sel.Distinct = true
		sql, err := restore(rootNode)
		if err != nil {
			return nil, errors.New("doFixMDistinctL: " +  err.Error())
		}
		// recover
		sel.Distinct = false
		return sql, nil
	case nil:
		return nil, errors.New("doFixMDistinctU: type error: nil")
	default:
		return nil, errors.New("doFixMDistinctL: type error: " + reflect.TypeOf(in).String())
	}
}