package sqlsim

import (
	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	_ "github.com/pingcap/tidb/parser/test_driver"
	"github.com/pkg/errors"
	"github.com/qaqcatz/impomysql/connector"
	"github.com/qaqcatz/impomysql/task"
)

// rmHint: remove optimization hint
func rmHint(bug *task.BugReport, conn *connector.Connector) error {
	sql2 := []*string {
		&(bug.OriginalSql),
		&(bug.MutatedSql),
	}
	res2 := []**connector.Result {
		&(bug.OriginalResult),
		&(bug.MutatedResult),
	}
	for i := 0; i < 2; i++ {
		tempSql, err := rmHintUnit(*sql2[i])
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

type rmHintVisitor struct {
}

func (v *rmHintVisitor) Enter(in ast.Node) (ast.Node, bool) {
	switch in.(type) {
	case *ast.TableName:
		tb := in.(*ast.TableName)
		tb.IndexHints = make([]*ast.IndexHint, 0)
	}
	return in, false
}

func (v *rmHintVisitor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

func rmHintUnit(sql string) (string, error) {
	p := parser.New()
	stmtNodes, _, err := p.Parse(sql, "", "")
	if err != nil {
		return "", errors.Wrap(err, "[rmHintUnit]parse error")
	}
	if stmtNodes == nil || len(stmtNodes) == 0 {
		return "", errors.New("[rmHintUnit]stmtNodes == nil || len(stmtNodes) == 0 ")
	}
	rootNode := &stmtNodes[0]

	v := &rmHintVisitor{}
	(*rootNode).Accept(v)

	simplifiedSql, err := restore(*rootNode)
	if err != nil {
		return "", errors.Wrap(err, "[rmHintUnit]restore error")
	}
	return string(simplifiedSql), nil
}
