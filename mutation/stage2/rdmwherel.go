package stage2

import (
	"errors"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/opcode"
	"github.com/pingcap/tidb/parser/test_driver"
	_ "github.com/pingcap/tidb/parser/test_driver"
	"math/rand"
	"reflect"
)
// addRdMWhereL: RdMWhereL: *ast.SelectStmt: WHERE xxx -> WHERE FALSE | WHERE (xxx) AND 0
func (v *MutateVisitor) addRdMWhereL(in *ast.SelectStmt, flag int) {
	if in.Where != nil {
		v.addCandidate(RdMWhereL, 0, in, flag)
	}
}

// doRdMWhereL: RdMWhereL: *ast.SelectStmt: WHERE xxx -> WHERE FALSE | WHERE (xxx) AND 0
func doRdMWhereL(rootNode ast.Node, in ast.Node, seed int64) ([]byte, error) {
	rander := rand.New(rand.NewSource(seed))
	switch in.(type) {
	case *ast.SelectStmt:
		sel := in.(*ast.SelectStmt)
		// check
		if sel.Where == nil {
			return nil, errors.New("doRdMWhereL: sel.Where == nil")
		}
		// mutate
		old := sel.Where

		rd := rander.Intn(2)
		if rd == 0 {
			// WHERE xxx -> WHERE FALSE
			sel.Where = &test_driver.ValueExpr{
				Datum: test_driver.NewDatum(0),
			}
		} else {
			// WHERE xxx -> WHERE (xxx) AND 0
			sel.Where = &ast.BinaryOperationExpr{
				Op: opcode.LogicAnd,
				L: old,
				R: &test_driver.ValueExpr{
					Datum: test_driver.NewDatum(0),
				},
			}
		}

		sql, err := restore(rootNode)
		if err != nil {
			return nil, errors.New("doRdMWhereL: " +  err.Error())
		}
		// recover
		sel.Where = old
		return sql, nil
	case nil:
		return nil, errors.New("doRdMWhereL: type error: nil")
	default:
		return nil, errors.New("doRdMWhereL: type error: " + reflect.TypeOf(in).String())
	}
}
