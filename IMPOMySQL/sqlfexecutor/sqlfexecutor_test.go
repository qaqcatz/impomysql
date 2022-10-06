package sqlfexecutor

import (
	"github.com/qaqcatz/IMPOMySQL/IMPOMySQL/testsqls"
	"io/ioutil"
	"testing"
)

func TestExtractSQL(t *testing.T) {
	err := testsqls.EnsureDBTEST()
	if err != nil {
		t.Fatal(err.Error())
	}
	conn, err := testsqls.GetConnector()
	if err != nil {
		t.Fatal(err.Error())
	}
	data, err, _ := testsqls.ReadSQLFile(testsqls.SQLFileQuote)
	if err != nil {
		t.Log(err.Error())
	}
	sqls := ExtractSQL(string(data))
	for i, sql := range sqls {
		t.Log(i, ":", sql)
		result := conn.ExecSQL(sql)
		if result.Err != nil {
			t.Fatal(result.Err.Error())
		}
		t.Log(result.ToString())
	}
}

func TestExtractSQL2(t *testing.T) {
	data, err, _ := testsqls.ReadSQLFile(testsqls.SQLFileTest)
	if err != nil {
		t.Log(err.Error())
	}
	sqls := ExtractSQL(string(data))
	if len(sqls) != 136 {
		t.Fatal("len(sqls) != 136")
	} else {
		t.Log(len(sqls))
	}
}

func testNewSQLFExecutorCommon(t *testing.T, sqlFile string) {
	data, err, _ := testsqls.ReadSQLFile(sqlFile)
	if err != nil {
		t.Fatal(err.Error())
	}
	sqlFExecutor, err := NewSQLFExecutorB(data)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(sqlFExecutor.ToString())
}

func TestNewSQLFExecutor(t *testing.T) {
	testNewSQLFExecutorCommon(t, testsqls.SQLFileTest)
}

func TestNewSQLFExecutor2(t *testing.T) {
	testNewSQLFExecutorCommon(t, testsqls.SQLFileWindow)
}

func testSQLFExecutor_ExecCommon(t *testing.T, sqlFile string) {
	err := testsqls.InitDBTEST()
	if err != nil {
		t.Fatal(err.Error())
	}
	conn, err := testsqls.GetConnector()
	if err != nil {
		t.Fatal(err.Error())
	}
	data, err, _ := testsqls.ReadSQLFile(sqlFile)
	if err != nil {
		t.Log(err.Error())
	}
	sqlFExecutor, err := NewSQLFExecutorB(data)
	if err != nil {
		t.Fatal(err.Error())
	}
	sqlFExecutor.Exec(conn)

	err = ioutil.WriteFile("./results_"+sqlFile+".txt", []byte(sqlFExecutor.ToString()), 0777)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestSQLFExecutor_Exec(t *testing.T) {
	testSQLFExecutor_ExecCommon(t, testsqls.SQLFileTest)
}

func TestSQLFExecutor_Exec2(t *testing.T) {
	testSQLFExecutor_ExecCommon(t, testsqls.SQLFileWindow)
}