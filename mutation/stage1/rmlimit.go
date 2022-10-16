package stage1

import (
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/test_driver"
	_ "github.com/pingcap/tidb/parser/test_driver"
)

// remove Limit
func RmLimit(in ast.Node) {
	if limit, ok := in.(*ast.Limit); ok {
		limit.Count = &test_driver.ValueExpr{
			Datum: test_driver.NewDatum(2147483647),
		}
		limit.Offset = nil
	}
}
