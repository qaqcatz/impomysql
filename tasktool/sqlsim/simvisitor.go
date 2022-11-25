package sqlsim

import (
	"github.com/pingcap/tidb/parser/ast"
	_ "github.com/pingcap/tidb/parser/test_driver"
)

type SimVisitor struct {
}

func (v *SimVisitor) Enter(in ast.Node) (ast.Node, bool) {

	return in, false
}

func (v *SimVisitor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}