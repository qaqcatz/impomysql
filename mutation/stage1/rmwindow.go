package stage1

import (
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/test_driver"
)

// rmWindow: remove window functions.
//
// For example:
//
// SELECT SUM(C1) OVER w as sum_c1 FROM T WINDOW w AS (...) -> SELECT 1 FROM T
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