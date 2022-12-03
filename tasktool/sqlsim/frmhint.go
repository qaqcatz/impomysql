package sqlsim

import (
	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	_ "github.com/pingcap/tidb/parser/test_driver"
	"github.com/pkg/errors"
	"github.com/qaqcatz/impomysql/connector"
	"github.com/qaqcatz/impomysql/task"
)

// frmHint: remove optimization hint
func frmHint(bug *task.BugReport, conn *connector.Connector) error {
	sql2 := []*string {
		&(bug.OriginalSql),
		&(bug.MutatedSql),
	}
	res2 := []**connector.Result {
		&(bug.OriginalResult),
		&(bug.MutatedResult),
	}
	for i := 0; i < 2; i++ {
		tempSql, err := frmHintUnit(*sql2[i])
		if err != nil {
			return err
		}

		tempResult := conn.ExecSQL(tempSql)
		if tempResult.Err == nil {
			cmp, err := (*res2[i]).CMP(tempResult)
			if err == nil && cmp == 0 {
				*sql2[i] = tempSql
				*res2[i] = tempResult
			}
		}
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
