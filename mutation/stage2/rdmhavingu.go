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

// addRdMHavingU: RdMHavingU, *ast.SelectStmt: HAVING xxx -> HAVING TRUE | HAVING (xxx) OR 1
func (v *MutateVisitor) addRdMHavingU(in *ast.SelectStmt, flag int) {
	if in.Having != nil && in.Having.Expr != nil {
		v.addCandidate(RdMHavingU, 1, in, flag)
	}
}

// doRdMHavingU: RdMHavingU, *ast.SelectStmt: HAVING xxx -> HAVING TRUE | HAVING (xxx) OR 1
func doRdMHavingU(rootNode ast.Node, in ast.Node, seed int64) ([]byte, error) {
	rander := rand.New(rand.NewSource(seed))
	switch in.(type) {
	case *ast.SelectStmt:
		sel := in.(*ast.SelectStmt)
		// check
		if sel.Having == nil || sel.Having.Expr == nil {
			return nil, errors.New("doRdMHavingU: sel.Having == nil || sel.Having.Expr == nil")
		}
		// mutate
		old := sel.Having.Expr

		rd := rander.Intn(2)
		if rd == 0 {
			// HAVING xxx -> HAVING TRUE
			sel.Having.Expr = &test_driver.ValueExpr{
				Datum: test_driver.NewDatum(1),
			}
		} else {
			// HAVING xxx -> HAVING (xxx) OR 1
			sel.Having.Expr = &ast.BinaryOperationExpr{
				Op: opcode.LogicOr,
				L: old,
				R: &test_driver.ValueExpr{
					Datum: test_driver.NewDatum(1),
				},
			}
		}

		sql, err := restore(rootNode)
		if err != nil {
			return nil, errors.New("doRdMHavingU: " +  err.Error())
		}
		// recover
		sel.Having.Expr = old
		return sql, nil
	case nil:
		return nil, errors.New("doRdMHavingU: type error: nil")
	default:
		return nil, errors.New("doRdMHavingU: type error: " + reflect.TypeOf(in).String())
	}
}