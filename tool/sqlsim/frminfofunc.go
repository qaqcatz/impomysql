package sqlsim

import (
	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	_ "github.com/pingcap/tidb/parser/test_driver"
	"github.com/pkg/errors"
	"github.com/qaqcatz/impomysql/connector"
	"github.com/qaqcatz/impomysql/mutation/oracle"
	"github.com/qaqcatz/impomysql/task"
	"strings"
)

// frmInfoFunc: remove information function:
//
// format_bytes -> charset
func frmInfoFunc(bug *task.BugReport, conn *connector.Connector) error {
	sql1, err := frmInfoFuncUnit(bug.OriginalSql)
	if err != nil {
		return err
	}
	sql2, err := frmInfoFuncUnit(bug.MutatedSql)
	if err != nil {
		return err
	}
	res1 := conn.ExecSQL(sql1)
	res2 := conn.ExecSQL(sql2)
	check, err := oracle.Check(res1, res2, bug.IsUpper)
	if err == nil && !check {
		bug.OriginalSql = sql1
		bug.OriginalResult = res1
		bug.MutatedSql = sql2
		bug.MutatedResult = res2
	}
	return nil
}

type frmInfoFuncVisitor struct {
}

func (v *frmInfoFuncVisitor) Enter(in ast.Node) (ast.Node, bool) {
	switch in.(type) {
	case *ast.FuncCallExpr:
		funcCall := in.(*ast.FuncCallExpr)
		if strings.ToLower(funcCall.FnName.String()) == "format_bytes" {
			funcCall.FnName.O = "charset"
			funcCall.FnName.L = "charset"
		}
	}
	return in, false
}

func (v *frmInfoFuncVisitor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

func frmInfoFuncUnit(sql string) (string, error) {
	p := parser.New()
	stmtNodes, _, err := p.Parse(sql, "", "")
	if err != nil {
		return "", errors.Wrap(err, "[frmInfoFuncUnit]parse error")
	}
	if stmtNodes == nil || len(stmtNodes) == 0 {
		return "", errors.New("[frmInfoFuncUnit]stmtNodes == nil || len(stmtNodes) == 0 ")
	}
	rootNode := &stmtNodes[0]

	v := &frmInfoFuncVisitor{}
	(*rootNode).Accept(v)

	simplifiedSql, err := restore(*rootNode)
	if err != nil {
		return "", errors.Wrap(err, "[frmInfoFuncUnit]restore error")
	}
	return string(simplifiedSql), nil
}