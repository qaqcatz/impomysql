package stage2

import (
	"errors"
	"github.com/pingcap/tidb/parser/opcode"
	"github.com/pingcap/tidb/parser/test_driver"
	_ "github.com/pingcap/tidb/parser/test_driver"
	"github.com/pingcap/tidb/parser/ast"
	"math/rand"
	"reflect"
)

// addRdMBetweenU: RdMBetweenU, *ast.BetweenExpr:
//   expr between l and r
//   ->
//   (expr) >= l and (expr) <= r
//   -> FixMCmpU, 1 and and (expr) <= r, (expr) >= l and 1 )
func (v *MutateVisitor) addRdMBetweenU(in *ast.BetweenExpr, flag int) {
	v.addCandidate(RdMBetweenU, 1, in, flag)
}

// doRdMBetweenU: RdMBetweenU, *ast.BetweenExpr:
//   expr between l and r
//   ->
//   (expr) >= l and (expr) <= r
//   -> FixMCmpU, 1 and and (expr) <= r, (expr) >= l and 1 )
func doRdMBetweenU(rootNode ast.Node, in ast.Node, seed int64) ([]byte, error) {
	rander := rand.New(rand.NewSource(seed))
	switch in.(type) {
	case *ast.BetweenExpr:
		btn := in.(*ast.BetweenExpr)
		// mutate
		var sql []byte = nil
		var err error = nil
		// trick:
		// X BETWEEN A AND B <=> 1 BETWEEN 0 AND (X >= A && X <= B);
		oldExpr := btn.Expr
		oldLeft := btn.Left
		oldRight := btn.Right
		newExpr := &test_driver.ValueExpr{
			Datum: test_driver.NewDatum(1),
		}
		newLeft := &test_driver.ValueExpr{
			Datum: test_driver.NewDatum(0),
		}
		newRight := &ast.BinaryOperationExpr{
			Op: opcode.LogicAnd,
			L: &ast.BinaryOperationExpr{
				Op: opcode.GE,
				L: oldExpr,
				R: oldLeft,
			},
			R: &ast.BinaryOperationExpr{
				Op: opcode.LE,
				L: oldExpr,
				R: oldRight,
			},
		}
		btn.Expr = newExpr
		btn.Left = newLeft
		btn.Right = &ast.ParenthesesExpr{
			Expr: newRight,
		}
		// -> FixMCmpU
		rd := rander.Intn(4)
		switch rd {
		case 0:
			// FixMCmpOpU, newRight.L
			sql, err = doFixMCmpU(rootNode, newRight.L)
		case 1:
			// FixMCmpOpU, newRight.R
			sql, err = doFixMCmpU(rootNode, newRight.R)
		case 2:
			// 1 and newRight.R
			newRight.L = &test_driver.ValueExpr{
				Datum: test_driver.NewDatum(1),
			}
			sql, err = restore(rootNode)
		case 3:
			// newRight.L and 1
			newRight.R = &test_driver.ValueExpr{
				Datum: test_driver.NewDatum(1),
			}
			sql, err = restore(rootNode)
		}
		if err != nil {
			return nil, errors.New("doRdMBetweenU: -> FixMCmpU: " + err.Error())
		}
		// recover
		btn.Expr = oldExpr
		btn.Left = oldLeft
		btn.Right = oldRight
		return sql, nil
	case nil:
		return nil, errors.New("doRdMBetweenU: type error: nil")
	default:
		return nil, errors.New("doRdMBetweenU: type error: " + reflect.TypeOf(in).String())
	}
}