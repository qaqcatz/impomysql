package stage1

import (
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/test_driver"
)

// rmWindow: remove window function.
// *ast.WindowFuncExpr can only appear in *ast.SelectField.
// In particular *ast.WindowSpec can also appear in *ast.SelectStmt. Therefore, we:
//   - Iterate each *ast.FieldList, replace each *ast.WindowFuncExpr with *test_driver.ValueExpr(value 1);
//   - set *ast.SelectStmt.WindowSpecs to nil(not empty, nil!).
//   example:
//   ----------input----------
//   SELECT ID AS id, CITY, AGE,
//   SUM(AGE) OVER w AS sum_age,
//   AVG(AGE) OVER (PARTITION BY CITY ORDER BY ID ROWS BETWEEN 1 PRECEDING AND 1 FOLLOWING) AS avg_age,
//   ROW_NUMBER() OVER (PARTITION BY CITY ORDER BY ID) AS rn
//   FROM COMPANY
//   WINDOW w AS (PARTITION BY CITY ORDER BY ID ROWS UNBOUNDED PRECEDING)
//   ----------output----------
//   SELECT ID AS id, CITY, AGE,
//   1 AS sum_age,
//   1 AS avg_age,
//   1 AS rn
//   FROM COMPANY
func rmWindow(in ast.Node) bool {
	change := false
	if selectStmt, ok := in.(*ast.SelectStmt); ok {
		if selectStmt.WindowSpecs != nil {
			change = true
			selectStmt.WindowSpecs = nil
		}
	}
	if fieldList, ok := in.(*ast.FieldList); ok {
		for _, field := range fieldList.Fields {
			if field.Expr == nil {
				continue
			}
			if _, ok := field.Expr.(*ast.WindowFuncExpr); ok {
				change = true
				field.Expr = &test_driver.ValueExpr{
					Datum: test_driver.NewDatum(1),
				}
			} // end of field.Expr.(*ast.WindowFuncExpr)
		} // end of range fieldList.Fields
	} // end of in.(*ast.FieldList)
	return change
}