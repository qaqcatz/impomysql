package sqlsim

import (
	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	_ "github.com/pingcap/tidb/parser/test_driver"
	"github.com/pkg/errors"
)

// rmWith: remove top WITH(non-recursive) and SELECT
func rmWith(sql string) (string, error) {
	p := parser.New()
	stmtNodes, _, err := p.Parse(sql, "", "")
	if err != nil {
		return "", errors.Wrap(err, "[rmWith]parse error")
	}
	if stmtNodes == nil || len(stmtNodes) == 0 {
		return "", errors.New("[rmWith]stmtNodes == nil || len(stmtNodes) == 0 ")
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
				return "", errors.Wrap(err, "[rmWith]restore error")
			}
			return string(simplifiedSql), nil
		}
	}
	return sql, nil
}
