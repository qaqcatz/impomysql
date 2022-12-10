package sqlsim

import (
	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	_ "github.com/pingcap/tidb/parser/test_driver"
	"github.com/pkg/errors"
	"github.com/qaqcatz/impomysql/connector"
	"github.com/qaqcatz/impomysql/mutation/oracle"
	"github.com/qaqcatz/impomysql/task"
)

// frmHint: remove optimization hint
func frmHint(bug *task.BugReport, conn *connector.Connector) error {
	sql1, err := frmHintUnit(bug.OriginalSql)
	if err != nil {
		return err
	}
	sql2, err := frmHintUnit(bug.MutatedSql)
	if err != nil {
		return err
	}
	res1 := conn.ExecSQL(sql1)
	res2 := conn.ExecSQL(sql2)
	check, err := oracle.Check(res1, res2, bug.IsUpper)
	if err == nil && !check {
		bug.OriginalSql = sql1
		bug.OriginalResult = res1
		bug.MutatedSql = sql2
		bug.MutatedResult = res2
	}
	return nil
}

type frmHintVisitor struct {
}

func (v *frmHintVisitor) Enter(in ast.Node) (ast.Node, bool) {
	switch in.(type) {
	case *ast.TableName:
		tb := in.(*ast.TableName)
		tb.IndexHints = make([]*ast.IndexHint, 0)
	}
	return in, false
}

func (v *frmHintVisitor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

func frmHintUnit(sql string) (string, error) {
	p := parser.New()
	stmtNodes, _, err := p.Parse(sql, "", "")
	if err != nil {
		return "", errors.Wrap(err, "[frmHintUnit]parse error")
	}
	if stmtNodes == nil || len(stmtNodes) == 0 {
		return "", errors.New("[frmHintUnit]stmtNodes == nil || len(stmtNodes) == 0 ")
	}
	rootNode := &stmtNodes[0]

	v := &frmHintVisitor{}
	(*rootNode).Accept(v)

	simplifiedSql, err := restore(*rootNode)
	if err != nil {
		return "", errors.Wrap(err, "[frmHintUnit]restore error")
	}
	return string(simplifiedSql), nil
}
