package sqlsimx

import (
	"github.com/pingcap/tidb/parser/opcode"
	"github.com/pkg/errors"
	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/test_driver"
	"github.com/qaqcatz/impomysql/connector"
)

// rmBinOp01: l binop r -> l/0/1 binop r/0/1
func rmBinOp01(sql string, result *connector.Result, conn *connector.Connector) (string, error) {
	p := parser.New()
	stmtNodes, _, err := p.Parse(sql, "", "")
	if err != nil {
		return "", errors.Wrap(err, "[rmBinOp01]parse error")
	}
	if stmtNodes == nil || len(stmtNodes) == 0 {
		return "", errors.New("[rmBinOp01]stmtNodes == nil || len(stmtNodes) == 0 ")
	}
	rootNode := &stmtNodes[0]

	v := &rmBinOp01Visitor{
		rootNode: *rootNode,
		result: result,
		conn: conn,
	}
	(*rootNode).Accept(v)

	simplifiedSql, err := restore(*rootNode)
	if err != nil {
		return "", errors.Wrap(err, "[rmBinOp01]restore error")
	}

	tempSql := string(simplifiedSql)
	tempResult := conn.ExecSQL(tempSql)
	cmp, err := tempResult.CMP(result)
	if err == nil && cmp == 0 {
		return tempSql, nil
	}

	return sql, nil
}

type rmBinOp01Visitor struct {
	rootNode ast.Node
	result *connector.Result
	conn *connector.Connector
}

func (v *rmBinOp01Visitor) Enter(in ast.Node) (ast.Node, bool) {
	switch in.(type) {
	case *ast.BinaryOperationExpr:
		binOpExpr := in.(*ast.BinaryOperationExpr)

		// only logical op can change to 0
		lw := 1
		if binOpExpr.Op == opcode.LogicOr || binOpExpr.Op == opcode.LogicAnd {
			lw = 0
		}

		oldL := binOpExpr.L
		oldR := binOpExpr.R

		for i := 1; i >= lw; i -= 1 {
			// l -> i
			binOpExpr.L = &test_driver.ValueExpr{
				Datum: test_driver.NewDatum(i),
			}
			simplifiedSql, err := restore(v.rootNode)
			if err == nil {
				tempSql := string(simplifiedSql)
				tempResult := v.conn.ExecSQL(tempSql)
				cmp, err := tempResult.CMP(v.result)
				if err == nil && cmp == 0 {
					break
				}
			}
			binOpExpr.L = oldL
		}

		for i := 1; i >= lw; i -= 1 {
			// r -> i
			binOpExpr.R = &test_driver.ValueExpr{
				Datum: test_driver.NewDatum(i),
			}
			simplifiedSql, err := restore(v.rootNode)
			if err == nil {
				tempSql := string(simplifiedSql)
				tempResult := v.conn.ExecSQL(tempSql)
				cmp, err := tempResult.CMP(v.result)
				if err == nil && cmp == 0 {
					break
				}
			}
			binOpExpr.R = oldR
		}
	}
	return in, false
}

func (v *rmBinOp01Visitor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}