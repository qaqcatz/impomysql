package stage1

import (
	"github.com/pingcap/tidb/parser/ast"
	_ "github.com/pingcap/tidb/parser/test_driver"
)

// rmLRJoin: change LEFT/RIGHT [OUTER] JOIN to JOIN
func rmLRJoin(in ast.Node) bool {
	if join, ok := in.(*ast.Join); ok {
		if join.Tp == ast.LeftJoin || join.Tp == ast.RightJoin {
			join.Tp = ast.CrossJoin
			join.NaturalJoin = false
			join.StraightJoin = false
			return true
		}
	}
	return false
}
