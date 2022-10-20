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

// addRdMHavingL: RdMHavingL, *ast.SelectStmt: HAVING xxx -> HAVING FALSE | HAVING (xxx) AND 0
func (v *MutateVisitor) addRdMHavingL(in *ast.SelectStmt, flag int) {
	if in.Having != nil && in.Having.Expr != nil {
		v.addCandidate(RdMHavingL, 0, in, flag)
	}
}

// doRdMHavingL: RdMHavingL, *ast.SelectStmt: HAVING xxx -> HAVING FALSE | HAVING (xxx) AND 0
func doRdMHavingL(rootNode ast.Node, in ast.Node, seed int64) ([]byte, error) {
	rander := rand.New(rand.NewSource(seed))
	switch in.(type) {
	case *ast.SelectStmt:
		sel := in.(*ast.SelectStmt)
		// check
		if sel.Having == nil || sel.Having.Expr == nil {
			return nil, errors.New("doRdMHavingL: sel.Having == nil || sel.Having.Expr == nil")
		}
		// mutate
		old := sel.Having.Expr

		rd := rander.Intn(2)
		if rd == 0 {
			// HAVING xxx -> HAVING FALSE
			sel.Having.Expr = &test_driver.ValueExpr{
				Datum: test_driver.NewDatum(0),
			}
		} else {
			// HAVING xxx -> HAVING (xxx) AND 0
			sel.Having.Expr = &ast.BinaryOperationExpr{
				Op: opcode.LogicAnd,
				L: old,
				R: &test_driver.ValueExpr{
					Datum: test_driver.NewDatum(0),
				},
			}
		}

		sql, err := restore(rootNode)
		if err != nil {
			return nil, errors.New("doRdMHavingL: " +  err.Error())
		}
		// recover
		sel.Having.Expr = old
		return sql, nil
	case nil:
		return nil, errors.New("doRdMHavingL: type error: nil")
	default:
		return nil, errors.New("doRdMHavingL: type error: " + reflect.TypeOf(in).String())
	}
}