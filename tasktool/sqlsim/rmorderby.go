package sqlsim

import (
	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	_ "github.com/pingcap/tidb/parser/test_driver"
	"github.com/pkg/errors"
	"github.com/qaqcatz/impomysql/connector"
	"github.com/qaqcatz/impomysql/task"
)

// rmOrderBy: remove ORDER BY
func rmOrderBy(bug *task.BugReport, conn *connector.Connector) error {
	sql2 := []*string {
		&(bug.OriginalSql),
		&(bug.MutatedSql),
	}
	res2 := []**connector.Result {
		&(bug.OriginalResult),
		&(bug.MutatedResult),
	}
	for i := 0; i < 2; i++ {
		tempSql, err := rmOrderByUnit(*sql2[i])
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

type rmOrderByVisitor struct {
}

func (v *rmOrderByVisitor) Enter(in ast.Node) (ast.Node, bool) {
	switch in.(type) {
	case *ast.SelectStmt:
		sel := in.(*ast.SelectStmt)
		sel.OrderBy = nil
	}
	return in, false
}

func (v *rmOrderByVisitor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

func rmOrderByUnit(sql string) (string, error) {
	p := parser.New()
	stmtNodes, _, err := p.Parse(sql, "", "")
	if err != nil {
		return "", errors.Wrap(err, "[rmOrderByUnit]parse error")
	}
	if stmtNodes == nil || len(stmtNodes) == 0 {
		return "", errors.New("[rmOrderByUnit]stmtNodes == nil || len(stmtNodes) == 0 ")
	}
	rootNode := &stmtNodes[0]

	v := &rmOrderByVisitor{}
	(*rootNode).Accept(v)

	simplifiedSql, err := restore(*rootNode)
	if err != nil {
		return "", errors.Wrap(err, "[rmOrderByUnit]restore error")
	}
	return string(simplifiedSql), nil
}