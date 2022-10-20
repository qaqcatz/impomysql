package stage2

import (
	"errors"
	"github.com/pingcap/tidb/parser/opcode"
	_ "github.com/pingcap/tidb/parser/test_driver"
	"github.com/pingcap/tidb/parser/ast"
	"reflect"
)

// addFixMCmpOpL: FixMCmpOpL, *ast.BinaryOperationExpr, *ast.CompareSubqueryExpr: a {>=|<=} b -> a {>|<} b
func (v *MutateVisitor) addFixMCmpOpL(in ast.Node, flag int) {
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
	case opcode.LE:
	case opcode.GE:
	default:
		return
	}
	v.addCandidate(FixMCmpOpL, 0, in, flag)
}

// doFixMCmpOpL: FixMCmpOpL, *ast.BinaryOperationExpr, *ast.CompareSubqueryExpr: a {>=|<=} b -> a {>|<} b
func doFixMCmpOpL(rootNode ast.Node, in ast.Node) ([]byte, error) {
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
		return nil, errors.New("doFixMCmpOpL: type error: nil")
	default:
		return nil, errors.New("doFixMCmpOpL: type error: " + reflect.TypeOf(in).String())
	}

	oldOp := *myOp
	var newOp opcode.Op
	switch oldOp {
	case opcode.LE:
		newOp = opcode.LT
	case opcode.GE:
		newOp = opcode.GT
	default:
		return nil, errors.New("doFixMCmpOpL: Op default")
	}
	// mutate
	*myOp = newOp
	sql, err := restore(rootNode)
	if err != nil {
		return nil, errors.New("doFixMCmpOpL: " +  err.Error())
	}
	// recover
	*myOp = oldOp
	return sql, nil
}