package stage1

import (
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/test_driver"
)

// remove window function
func RmWindow(in ast.Node) {
	if selectStmt, ok := in.(*ast.SelectStmt); ok {
		selectStmt.WindowSpecs = nil
	}
	if fieldList, ok := in.(*ast.FieldList); ok {
		for _, field := range fieldList.Fields {
			if field.Expr == nil {
				continue
			}
			if _, ok := field.Expr.(*ast.WindowFuncExpr); ok {
				field.Expr = &test_driver.ValueExpr{
					Datum: test_driver.NewDatum(1),
				}
			} // end of field.Expr.(*ast.WindowFuncExpr)
		} // end of range fieldList.Fields
	} // end of in.(*ast.FieldList)
}