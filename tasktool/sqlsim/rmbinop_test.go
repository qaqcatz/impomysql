package sqlsim

import (
	"github.com/qaqcatz/impomysql/testsqls"
	"testing"
)

func TestRmBinOp(t *testing.T) {
	err := testsqls.InitTableCOMPANY("")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	conn, err := testsqls.GetConnector("")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	sql := testsqls.SQLBinaryOp2
	result := conn.ExecSQL(sql)
	if result.Err != nil {
		t.Fatalf("%+v", result.Err)
	} else {
		t.Log(sql+"\n", result.ToString())
	}

	newSql, err := rmBinOpTrueAllUnit(sql, result, conn)
	result = conn.ExecSQL(newSql)
	if result.Err != nil {
		t.Fatalf("%+v", result.Err)
	} else {
		t.Log(newSql+"\n", result.ToString())
	}

	newSql, err = rmBinOpFalseAllUnit(newSql, result, conn)
	result = conn.ExecSQL(newSql)
	if result.Err != nil {
		t.Fatalf("%+v", result.Err)
	} else {
		t.Log(newSql+"\n", result.ToString())
	}
}
