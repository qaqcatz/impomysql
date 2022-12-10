package sqlsimx

import (
	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/opcode"
	"github.com/pingcap/tidb/parser/test_driver"
	"github.com/pkg/errors"
	"github.com/qaqcatz/impomysql/connector"
)

// rmBinOpFalse: xxx AND/OR xxx -> 0 AND/OR xxx | xxx AND/OR 0
func rmBinOpFalse(sql string, result *connector.Result, conn *connector.Connector) (string, error) {

	// init rmBinOpVisitor, the first goal of traversal is to get binOpExprValueNum
	v := &rmBinOpFalseVisitor{binOpExprValueNum: 0, isChangedBinOpExprValue: false,
		changedBinOpExprValueNum: 0, isFirstEnter: true, cursorBinOpExprValue: 0}
	_, err := rmBinOpFalseUnit(sql, v)
	if err != nil {
		return sql, err
	}
	v.isFirstEnter = false

	// rmBinOpExprValue
	for i := 0; i < v.binOpExprValueNum; i++ {
		v.cursorBinOpExprValue = 0
		v.isChangedBinOpExprValue = false
		tempSql, err := rmBinOpFalseUnit(sql, v)
		if err != nil {
			return sql, err
		}

		tempResult := conn.ExecSQL(tempSql)
		if tempResult.Err != nil {
			continue
		}
		cmp, err := tempResult.CMP(result)
		if err == nil && cmp == 0 {
			sql = tempSql
			result = tempResult
		}
	}
	return sql, nil
}

type rmBinOpFalseVisitor struct {
	binOpExprValueNum        int
	isChangedBinOpExprValue  bool
	changedBinOpExprValueNum int
	isFirstEnter             bool
	cursorBinOpExprValue     int
}

func (v *rmBinOpFalseVisitor) Enter(in ast.Node) (ast.Node, bool) {
	if v.isFirstEnter == true {
		switch in.(type) {
		case *ast.BinaryOperationExpr:
			binOpExpr := in.(*ast.BinaryOperationExpr)
			if binOpExpr.Op == opcode.LogicOr || binOpExpr.Op == opcode.LogicAnd {
				v.binOpExprValueNum += 2
			}
		}
		return in, false
	} else {
		if v.isChangedBinOpExprValue == true {
			return in, false
		}
		switch in.(type) {
		case *ast.BinaryOperationExpr:
			binOpExpr := in.(*ast.BinaryOperationExpr)
			if binOpExpr.Op == opcode.LogicOr || binOpExpr.Op == opcode.LogicAnd {
				v.cursorBinOpExprValue += 2
				if v.cursorBinOpExprValue <= v.changedBinOpExprValueNum {
					return in, false
				}

				if v.changedBinOpExprValueNum%2 == 0 {
					binOpExpr.L = &test_driver.ValueExpr{
						Datum: test_driver.NewDatum(0),
					}
				} else {
					binOpExpr.R = &test_driver.ValueExpr{
						Datum: test_driver.NewDatum(0),
					}
				}
				v.changedBinOpExprValueNum++
				v.isChangedBinOpExprValue = true
			}
		}
		return in, false
	}
}

func (v *rmBinOpFalseVisitor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

func rmBinOpFalseUnit(sql string, v *rmBinOpFalseVisitor) (string, error) {
	p := parser.New()
	stmtNodes, _, err := p.Parse(sql, "", "")
	if err != nil {
		return "", errors.Wrap(err, "[rmBinOpFalseUnit]parse error")
	}
	if stmtNodes == nil || len(stmtNodes) == 0 {
		return "", errors.New("[rmBinOpFalseUnit]stmtNodes == nil || len(stmtNodes) == 0 ")
	}
	rootNode := &stmtNodes[0]

	(*rootNode).Accept(v)

	simplifiedSql, err := restore(*rootNode)
	if err != nil {
		return "", errors.Wrap(err, "[rmBinOpFalseUnit]restore error")
	}
	return string(simplifiedSql), nil
}
