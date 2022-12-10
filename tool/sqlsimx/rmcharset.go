package sqlsimx

import (
	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/test_driver"
	"github.com/pkg/errors"
	"github.com/qaqcatz/impomysql/connector"
)

// rmCharset: remove charset
func rmCharset(sql string, result *connector.Result, conn *connector.Connector) (string, error) {
	p := parser.New()
	stmtNodes, _, err := p.Parse(sql, "", "")
	if err != nil {
		return "", errors.Wrap(err, "[rmCharset]parse error")
	}
	if stmtNodes == nil || len(stmtNodes) == 0 {
		return "", errors.New("[rmCharset]stmtNodes == nil || len(stmtNodes) == 0 ")
	}
	rootNode := &stmtNodes[0]

	v := &frmCharsetVisitor{}
	(*rootNode).Accept(v)

	simplifiedSql, err := restore(*rootNode)
	if err != nil {
		return "", errors.Wrap(err, "[rmCharset]restore error")
	}

	tempSql := string(simplifiedSql)
	tempResult := conn.ExecSQL(tempSql)
	cmp, err := tempResult.CMP(result)
	if err == nil && cmp == 0 {
		return tempSql, nil
	}

	return sql, nil
}

type frmCharsetVisitor struct {
}

func (v *frmCharsetVisitor) Enter(in ast.Node) (ast.Node, bool) {
	switch in.(type) {
	case *test_driver.ValueExpr:
		valueExpr := in.(*test_driver.ValueExpr)
		valueExpr.Type.Charset = ""
	}
	return in, false
}

func (v *frmCharsetVisitor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}
