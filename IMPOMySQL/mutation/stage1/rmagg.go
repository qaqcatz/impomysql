package stage1

import (
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/test_driver"
)

// remove aggregate function and group by
func RmAgg(in ast.Node) {
	if selectStmt, ok := in.(*ast.SelectStmt); ok {
		selectStmt.GroupBy = nil
	}
	if aggFunExpr, ok := in.(*ast.AggregateFuncExpr); ok {
		aggFunExpr.F = ""
		aggFunExpr.Distinct = false
		aggFunExpr.Order = nil
		aggFunExpr.Args = make([]ast.ExprNode, 0)
		aggFunExpr.Args = append(aggFunExpr.Args, &test_driver.ValueExpr{
			Datum: test_driver.NewDatum(1),
		})
	}
}
