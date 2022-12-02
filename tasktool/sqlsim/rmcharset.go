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

// rmCharset: remove charset
func rmCharset(bug *task.BugReport, conn *connector.Connector) error {
	sql1, err := rmCharsetUnit(bug.OriginalSql)
	if err != nil {
		return err
	}
	sql2, err := rmCharsetUnit(bug.MutatedSql)
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

type rmCharsetVisitor struct {
}

func (v *rmCharsetVisitor) Enter(in ast.Node) (ast.Node, bool) {
	switch in.(type) {
	case *test_driver.ValueExpr:
		valueExpr := in.(*test_driver.ValueExpr)
		valueExpr.Type.Charset = ""
	}
	return in, false
}

func (v *rmCharsetVisitor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

func rmCharsetUnit(sql string) (string, error) {
	p := parser.New()
	stmtNodes, _, err := p.Parse(sql, "", "")
	if err != nil {
		return "", errors.Wrap(err, "[rmCharsetUnit]parse error")
	}
	if stmtNodes == nil || len(stmtNodes) == 0 {
		return "", errors.New("[rmCharsetUnit]stmtNodes == nil || len(stmtNodes) == 0 ")
	}
	rootNode := &stmtNodes[0]

	v := &rmCharsetVisitor{}
	(*rootNode).Accept(v)

	simplifiedSql, err := restore(*rootNode)
	if err != nil {
		return "", errors.Wrap(err, "[rmCharsetUnit]restore error")
	}
	return string(simplifiedSql), nil
}
