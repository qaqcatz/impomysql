package stage1

import (
	"github.com/pingcap/tidb/parser/ast"
	_ "github.com/pingcap/tidb/parser/test_driver"
)

type InitVisitor struct {
}

func (v *InitVisitor) Enter(in ast.Node) (ast.Node, bool) {
	rmAgg(in)
	rmWindow(in)
	rmLRJoin(in)
	rmLimit(in)
	rmUncertain(in)
	return in, false
}

func (v *InitVisitor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}