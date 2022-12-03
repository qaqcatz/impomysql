package sqlsim

import (
	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/test_driver"
	"github.com/pkg/errors"
	"github.com/qaqcatz/impomysql/connector"
	"github.com/qaqcatz/impomysql/mutation/oracle"
	"github.com/qaqcatz/impomysql/task"
)

// frmCharset: remove charset
func frmCharset(bug *task.BugReport, conn *connector.Connector) error {
	sql1, err := frmCharsetUnit(bug.OriginalSql)
	if err != nil {
		return err
	}
	sql2, err := frmCharsetUnit(bug.MutatedSql)
	if err != nil {
		return err
	}
	res1 := conn.ExecSQL(sql1)
	res2 := conn.ExecSQL(sql2)
	check, err := oracle.Check(res1, res2, bug.IsUpper)
	if !check {
		bug.OriginalSql = sql1
		bug.OriginalResult = res1
		bug.MutatedSql = sql2
		bug.MutatedResult = res2
	}
	return nil
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

func frmCharsetUnit(sql string) (string, error) {
	p := parser.New()
	stmtNodes, _, err := p.Parse(sql, "", "")
	if err != nil {
		return "", errors.Wrap(err, "[frmCharsetUnit]parse error")
	}
	if stmtNodes == nil || len(stmtNodes) == 0 {
		return "", errors.New("[frmCharsetUnit]stmtNodes == nil || len(stmtNodes) == 0 ")
	}
	rootNode := &stmtNodes[0]

	v := &frmCharsetVisitor{}
	(*rootNode).Accept(v)

	simplifiedSql, err := restore(*rootNode)
	if err != nil {
		return "", errors.Wrap(err, "[frmCharsetUnit]restore error")
	}
	return string(simplifiedSql), nil
}
