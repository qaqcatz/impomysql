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

// frmStrFunc: remove string function:
//
// to_base64 -> oct
func frmStrFunc(bug *task.BugReport, conn *connector.Connector) error {
	sql1, err := frmStrFuncUnit(bug.OriginalSql)
	if err != nil {
		return err
	}
	sql2, err := frmStrFuncUnit(bug.MutatedSql)
	if err != nil {
		return err
	}
	res1 := conn.ExecSQL(sql1)
	res2 := conn.ExecSQL(sql2)
	check, err := oracle.Check(res1, res2, bug.IsUpper)
	if !check {
		bug.OriginalSql = sql1
		bug.OriginalResult = res1
		bug.MutatedSql = sql2
		bug.MutatedResult = res2
	}
	return nil
}

type frmStrFuncVisitor struct {
}

func (v *frmStrFuncVisitor) Enter(in ast.Node) (ast.Node, bool) {
	switch in.(type) {
	case *ast.FuncCallExpr:
		funcCall := in.(*ast.FuncCallExpr)
		if strings.ToLower(funcCall.FnName.String()) == "to_base64" {
			funcCall.FnName.O = "oct"
			funcCall.FnName.L = "oct"
		}
	}
	return in, false
}

func (v *frmStrFuncVisitor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

func frmStrFuncUnit(sql string) (string, error) {
	p := parser.New()
	stmtNodes, _, err := p.Parse(sql, "", "")
	if err != nil {
		return "", errors.Wrap(err, "[frmStrFuncUnit]parse error")
	}
	if stmtNodes == nil || len(stmtNodes) == 0 {
		return "", errors.New("[frmStrFuncUnit]stmtNodes == nil || len(stmtNodes) == 0 ")
	}
	rootNode := &stmtNodes[0]

	v := &frmStrFuncVisitor{}
	(*rootNode).Accept(v)

	simplifiedSql, err := restore(*rootNode)
	if err != nil {
		return "", errors.Wrap(err, "[frmStrFuncUnit]restore error")
	}
	return string(simplifiedSql), nil
}