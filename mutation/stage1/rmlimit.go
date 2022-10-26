package stage1

import (
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/test_driver"
	_ "github.com/pingcap/tidb/parser/test_driver"
)

// rmLimit: limit x -> limit 2147483647
func rmLimit(in ast.Node) bool {
	if limit, ok := in.(*ast.Limit); ok {
		limit.Count = &test_driver.ValueExpr{
			Datum: test_driver.NewDatum(2147483647),
		}
		limit.Offset = nil
		return true
	}
	return false
}
