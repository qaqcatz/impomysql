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

// frmWith: remove top WITH(non-recursive) and SELECT
func frmWith(bug *task.BugReport, conn *connector.Connector) error {
	sql1, err := frmWithUnit(bug.OriginalSql)
	if err != nil {
		return err
	}
	sql2, err := frmWithUnit(bug.MutatedSql)
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

func frmWithUnit(sql string) (string, error) {
	p := parser.New()
	stmtNodes, _, err := p.Parse(sql, "", "")
	if err != nil {
		return "", errors.Wrap(err, "[frmWithUnit]parse error")
	}
	if stmtNodes == nil || len(stmtNodes) == 0 {
		return "", errors.New("[frmWithUnit]stmtNodes == nil || len(stmtNodes) == 0 ")
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
				return "", errors.Wrap(err, "[frmWithUnit]restore error")
			}
			return string(simplifiedSql), nil
		}
	}
	return sql, nil
}
