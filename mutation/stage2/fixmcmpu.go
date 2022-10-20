package stage2

import (
	"errors"
	"github.com/pingcap/tidb/parser/opcode"
	"github.com/pingcap/tidb/parser/test_driver"
	_ "github.com/pingcap/tidb/parser/test_driver"
	"github.com/pingcap/tidb/parser/ast"
	"reflect"
)

// addFixMCmpU: FixMCmpU, *ast.BinaryOperationExpr, *ast.CompareSubqueryExpr:
//
// a {>|>=} b -> (a) + 1 {>|>=} (b) + 0
//
// a {<|<=} b -> (a) + 0 {<|<=} (b) + 1
//
// may false positive
func (v *MutateVisitor) addFixMCmpU(in ast.Node, flag int) {
	var myOp *opcode.Op = nil
	switch in.(type) {
	case *ast.BinaryOperationExpr:
		bin := in.(*ast.BinaryOperationExpr)
		myOp = &bin.Op
	case *ast.CompareSubqueryExpr:
		cmp := in.(*ast.CompareSubqueryExpr)
		myOp = &cmp.Op
	default:
		return
	}
	switch *myOp {
	case opcode.LT:
	case opcode.GT:
	case opcode.LE:
	case opcode.GE:
	default:
		return
	}
	v.addCandidate(FixMCmpU, 1, in, flag)
}

// doFixMCmpU: FixMCmpU, *ast.BinaryOperationExpr, *ast.CompareSubqueryExpr:
//
// a {>|>=} b -> (a) + 1 {>|>=} (b) + 0
//
// a {<|<=} b -> (a) + 0 {<|<=} (b) + 1
//
// may false positive
func doFixMCmpU(rootNode ast.Node, in ast.Node) ([]byte, error) {
	// check
	var myOp *opcode.Op = nil
	var myL *ast.ExprNode = nil
	var myR *ast.ExprNode = nil
	switch in.(type) {
	case *ast.BinaryOperationExpr:
		bin := in.(*ast.BinaryOperationExpr)
		myOp = &bin.Op
		myL = &bin.L
		myR = &bin.R
	case *ast.CompareSubqueryExpr:
		cmp := in.(*ast.CompareSubqueryExpr)
		myOp = &cmp.Op
		myL = &cmp.L
		myR = &cmp.R
	case nil:
		return nil, errors.New("doFixMCmpU: type error: nil")
	default:
		return nil, errors.New("doFixMCmpU: type error: " + reflect.TypeOf(in).String())
	}

	oldL := *myL
	oldR := *myR
	var newL ast.ExprNode
	var newR ast.ExprNode
	switch *myOp {
	case opcode.LT, opcode.LE:
		// a {<|<=} b -> (a) + 0 {<|<=} (b) + 1
		newL = &ast.BinaryOperationExpr {
			Op: opcode.Plus,
			L: oldL,
			R: &test_driver.ValueExpr{
				Datum: test_driver.NewDatum(0),
			},
		}
		newR = &ast.BinaryOperationExpr{
			Op: opcode.Plus,
			L: oldR,
			R: &test_driver.ValueExpr{
				Datum: test_driver.NewDatum(1),
			},
		}
	case opcode.GT | opcode.GE:
		// a {>|>=} b -> (a) + 1 {>|>=} (b) + 0
		newL = &ast.BinaryOperationExpr {
			Op: opcode.Plus,
			L: oldL,
			R: &test_driver.ValueExpr{
				Datum: test_driver.NewDatum(1),
			},
		}
		newR = &ast.BinaryOperationExpr{
			Op: opcode.Plus,
			L: oldR,
			R: &test_driver.ValueExpr{
				Datum: test_driver.NewDatum(0),
			},
		}
	default:
		return nil, errors.New("doFixMCmpU: Op default")
	}
	// mutate
	*myL = newL
	*myR = newR
	sql, err := restore(rootNode)
	if err != nil {
		return nil, errors.New("doFixMCmpU: " +  err.Error())
	}
	// recover
	*myL = oldL
	*myR = oldR
	return sql, nil
}