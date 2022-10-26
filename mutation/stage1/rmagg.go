package stage1

import (
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/test_driver"
)

// rmAgg: remove aggregate function and group by.
// For a node *ast.AggregateFuncExpr, we:
//   - set node.F = "";
//   - clear node.Args, add *test_driver.ValueExpr(value 1) to node.Args;
//   - set node.Distinct to false;
//   - set node.Order to nil.
//   In particular, we need to set *ast.SelectStmt.GroupBy to nil
//   to avoid the semantic error caused by removing the aggregate functions.
//   example:
//   ----------input----------
//   SELECT * FROM (
//     SELECT SUM(ID+1) AS S, GROUP_CONCAT(NAME ORDER BY NAME DESC), CITY
//     FROM COMPANY
//     GROUP BY CITY
//     HAVING COUNT(DISTINCT AGE) >= 1
//   ) AS T
//   WHERE T.S > 0;
//   ----------output----------
//   SELECT * FROM (
//     SELECT (1) AS S, (1), CITY
//     FROM COMPANY
//     HAVING (1) >= 1
//   ) AS T
//   WHERE T.S > 0;
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
