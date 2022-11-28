package sqlsim

import (
	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	_ "github.com/pingcap/tidb/parser/test_driver"
	"github.com/pkg/errors"
	"github.com/qaqcatz/impomysql/connector"
	"github.com/qaqcatz/impomysql/task"
)

// rmWith: remove top WITH(non-recursive) and SELECT
func rmWith(bug *task.BugReport, conn *connector.Connector) error {
	sql2 := []*string {
		&(bug.OriginalSql),
		&(bug.MutatedSql),
	}
	res2 := []**connector.Result {
		&(bug.OriginalResult),
		&(bug.MutatedResult),
	}
	for i := 0; i < 2; i++ {
		tempSql, err := rmWithUnit(*sql2[i])
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

func rmWithUnit(sql string) (string, error) {
	p := parser.New()
	stmtNodes, _, err := p.Parse(sql, "", "")
	if err != nil {
		return "", errors.Wrap(err, "[rmWithUnit]parse error")
	}
	if stmtNodes == nil || len(stmtNodes) == 0 {
		return "", errors.New("[rmWithUnit]stmtNodes == nil || len(stmtNodes) == 0 ")
	}
	rootNode := &stmtNodes[0]

	switch (*rootNode).(type) {
	case *ast.SelectStmt:
		sel := (*rootNode).(*ast.SelectStmt)
		with := sel.With
		if with == nil {
			break
		}
		if with.IsRecursive {
			break
		}
		if len(with.CTEs) != 1 {
			break
		}
		if with.CTEs[0].Query == nil {
			break
		}
		switch (with.CTEs[0].Query.Query).(type) {
		case *ast.SetOprStmt, *ast.SelectStmt:
			simplifiedSql, err := restore(with.CTEs[0].Query.Query)
			if err != nil {
				return "", errors.Wrap(err, "[rmWithUnit]restore error")
			}
			return string(simplifiedSql), nil
		}
	}
	return sql, nil
}
