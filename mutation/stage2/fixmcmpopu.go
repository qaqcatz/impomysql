package stage2

import (
	"errors"
	"github.com/pingcap/tidb/parser/opcode"
	_ "github.com/pingcap/tidb/parser/test_driver"
	"github.com/pingcap/tidb/parser/ast"
	"reflect"
)

// addFixMCmpOpU: FixMCmpOpU, *ast.BinaryOperationExpr, *ast.CompareSubqueryExpr: a {>|<|=} b -> a {>=|<=|>=} b
func (v *MutateVisitor) addFixMCmpOpU(in ast.Node, flag int) {
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
	case opcode.EQ:
	case opcode.LT:
	case opcode.GT:
	default:
		return
	}
	v.addCandidate(FixMCmpOpU, 1, in, flag)
}

// doFixMCmpOpU: FixMCmpOpU, *ast.BinaryOperationExpr, *ast.CompareSubqueryExpr: a {>|<|=} b -> a {>=|<=|>=} b
func doFixMCmpOpU(rootNode ast.Node, in ast.Node) ([]byte, error) {
	// check
	var myOp *opcode.Op = nil
	switch in.(type) {
	case *ast.BinaryOperationExpr:
		bin := in.(*ast.BinaryOperationExpr)
		myOp = &bin.Op
	case *ast.CompareSubqueryExpr:
		cmp := in.(*ast.CompareSubqueryExpr)
		myOp = &cmp.Op
	case nil:
		return nil, errors.New("doFixMCmpOpU: type error: nil")
	default:
		return nil, errors.New("doFixMCmpOpU: type error: " + reflect.TypeOf(in).String())
	}

	oldOp := *myOp
	var newOp opcode.Op
	switch oldOp {
	case opcode.EQ:
		newOp = opcode.GE
	case opcode.LT:
		newOp = opcode.LE
	case opcode.GT:
		newOp = opcode.GE
	default:
		return nil, errors.New("doFixMCmpOpU: Op default")
	}
	// mutate
	*myOp = newOp
	sql, err := restore(rootNode)
	if err != nil {
		return nil, errors.New("doFixMCmpOpU: " +  err.Error())
	}
	// recover
	*myOp = oldOp
	return sql, nil
}