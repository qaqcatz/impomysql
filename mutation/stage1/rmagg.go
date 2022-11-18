package stage1

import (
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/test_driver"
)

// rmAgg: remove aggregate functions and GROUP BY.
//
// For example:
//
// SELECT C1, SUM(C2) FROM T GROUP BY C1 -> SELECT C1, (1) FROM T
func rmAgg(in ast.Node) bool {
	change := false
	if selectStmt, ok := in.(*ast.SelectStmt); ok {
		if selectStmt.GroupBy != nil {
			change = true
			selectStmt.GroupBy = nil
		}
	}
	if aggFunExpr, ok := in.(*ast.AggregateFuncExpr); ok {
		change = true
		aggFunExpr.F = ""
		aggFunExpr.Distinct = false
		aggFunExpr.Order = nil
		aggFunExpr.Args = make([]ast.ExprNode, 0)
		aggFunExpr.Args = append(aggFunExpr.Args, &test_driver.ValueExpr{
			Datum: test_driver.NewDatum(1),
		})
	}
	return change
}
