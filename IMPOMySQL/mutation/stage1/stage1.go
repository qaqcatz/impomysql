package stage1

import (
	"bytes"
	"errors"
	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/format"
	_ "github.com/pingcap/tidb/parser/test_driver"
)

// InitVisitor: Remove aggregate function(and group by), window function.
type InitVisitor struct {
}

func (v *InitVisitor) Enter(in ast.Node) (ast.Node, bool) {
	RmAgg(in)
	RmWindow(in)
	return in, false
}

func (v *InitVisitor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

// Stage1: Remove aggregate function(and group by), window function.
//
// The transformed sql may fail to execute. It is recommended to execute
// the transformed sql to do some verification.
//
// Only Support SELECT statement.
func Stage1(sql string) (string, error) {
	p := parser.New()
	stmtNodes, _, err := p.Parse(sql, "", "")
	if err != nil {
		return "", errors.New("Stage1: p.Parse() error: " + err.Error())
	}
	if stmtNodes == nil || len(stmtNodes) == 0 {
		return "", errors.New("Stage1: stmtNodes == nil || len(stmtNodes) == 0 ")
	}
	rootNode := &stmtNodes[0]
	if selectStmt, ok := (*rootNode).(*ast.SelectStmt); ok {
		v := &InitVisitor{}
		selectStmt.Accept(v)

		buf := new(bytes.Buffer)
		ctx := format.NewRestoreCtx(format.DefaultRestoreFlags, buf)
		err := selectStmt.Restore(ctx)
		if err != nil {
			return "", errors.New("Stage1: selectStmt.Restore() error: " + err.Error())
		}
		return buf.String(), nil
	} else {
		return "", errors.New("Stage1: *rootNode is not *ast.SelectStmt ")
	}
}
