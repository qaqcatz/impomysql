package sqlsim

import (
	"fmt"
	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/opcode"
	"github.com/pingcap/tidb/parser/test_driver"
	"github.com/pkg/errors"
)

type rmBinOpVisitor struct {
}

func (v *rmBinOpVisitor) Enter(in ast.Node) (ast.Node, bool) {
	switch in.(type) {
	case *ast.BinaryOperationExpr:
		binOpExpr := in.(*ast.BinaryOperationExpr)
		if binOpExpr.Op == opcode.LogicOr {
			binOpExpr.L = &test_driver.ValueExpr{
				Datum: test_driver.NewDatum(1),
			}
		} else if binOpExpr.Op == opcode.LogicAnd {
			fmt.Println(binOpExpr.Op)
		}
	}
	return in, false
}

func (v *rmBinOpVisitor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

func rmBinOp(sql string) (string, error) {
	p := parser.New()
	stmtNodes, _, err := p.Parse(sql, "", "")
	if err != nil {
		return "", errors.Wrap(err, "[rmBinOp]parse error")
	}
	if stmtNodes == nil || len(stmtNodes) == 0 {
		return "", errors.New("[rmBinOp]stmtNodes == nil || len(stmtNodes) == 0 ")
	}
	rootNode := &stmtNodes[0]

	v := &rmBinOpVisitor{}
	(*rootNode).Accept(v)

	simplifiedSql, err := restore(*rootNode)
	if err != nil {
		return "", errors.Wrap(err, "[rmBinOp]restore error")
	}
	fmt.Println(string(simplifiedSql))
	return sql, nil
}