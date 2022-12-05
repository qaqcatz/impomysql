package sqlsim

import (
	"github.com/qaqcatz/impomysql/testsqls"
	"testing"
)

func TestFrmInfoFuncUnit(t *testing.T) {
	err := testsqls.InitTableCOMPANY("")
	if err != nil {
		t.Fatalf("%+v\n", err)
	}
	conn, err := testsqls.GetConnector("")
	if err != nil {
		t.Fatalf("%+v\n", err)
	}
	simplifiedSql, err := frmInfoFuncUnit(testsqls.SQLInfoFunc)
	if err != nil {
		t.Fatalf("%+v\n", err)
	}
	t.Log(testsqls.SQLInfoFunc)
	res1 := conn.ExecSQL(testsqls.SQLInfoFunc)
	if res1.Err != nil {
		t.Fatalf("%+v\n", res1.Err)
	}
	t.Log(res1.ToString())
	t.Log(simplifiedSql)
	res2 := conn.ExecSQL(simplifiedSql)
	if res2.Err != nil {
		t.Fatalf("%+v\n", res2.Err)
	}
	t.Log(res2.ToString())
}
