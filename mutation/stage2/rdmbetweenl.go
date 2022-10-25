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

// addRdMBetweenL: RdMBetweenL, *ast.BetweenExpr:
//   expr between l and r
//   ->
//   (expr) >= l and (expr) <= r
//   -> FixMCmpOpL / FixMCmpL )
func (v *MutateVisitor) addRdMBetweenL(in *ast.BetweenExpr, flag int) {
	v.addCandidate(RdMBetweenL, 0, in, flag)
}

// doRdMBetweenL: RdMBetweenL, *ast.BetweenExpr:
//   expr between l and r
//   ->
//   (expr) >= l and (expr) <= r
//   -> FixMCmpOpL / FixMCmpL )
func doRdMBetweenL(rootNode ast.Node, in ast.Node, seed int64) ([]byte, error) {
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
		// -> FixMCmpOpL / FixMCmpL
		rd := rander.Intn(6)
		switch rd {
		case 0:
			// FixMCmpOpL, newRight.L
			sql, err = doFixMCmpOpL(rootNode, newRight.L)
		case 1:
			// FixMCmpOpL, newRight.R
			sql, err = doFixMCmpOpL(rootNode, newRight.R)
		case 2:
			// FixMCmpL, newRight.L
			sql, err = doFixMCmpL(rootNode, newRight.L)
		case 3:
			// FixMCmpL, newRight.R
			sql, err = doFixMCmpL(rootNode, newRight.R)
		}
		if err != nil {
			return nil, errors.New("doRdMBetweenL: -> FixMCmpOpL / FixMCmpL: " + err.Error())
		}
		// recover
		btn.Expr = oldExpr
		btn.Left = oldLeft
		btn.Right = oldRight
		return sql, nil
	case nil:
		return nil, errors.New("doRdMBetweenL: type error: nil")
	default:
		return nil, errors.New("doRdMBetweenL: type error: " + reflect.TypeOf(in).String())
	}
}