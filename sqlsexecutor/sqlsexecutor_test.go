package sqlsexecutor

import (
	"github.com/qaqcatz/impomysql/testsqls"
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

func testNewSQLSExecutorCommon(t *testing.T, sqlFile string) {
	data, err, _ := testsqls.ReadSQLFile(sqlFile)
	if err != nil {
		t.Fatal(err.Error())
	}
	sqlsExecutor, err := NewSQLSExecutorB(data)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(sqlsExecutor.ToString())
}

func TestNewSQLSExecutor(t *testing.T) {
	testNewSQLSExecutorCommon(t, testsqls.SQLFileTest)
}

func TestNewSQLSExecutor2(t *testing.T) {
	testNewSQLSExecutorCommon(t, testsqls.SQLFileWindow)
}

func testSQLSExecutor_ExecCommon(t *testing.T, sqlFile string) {
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
	sqlsExecutor, err := NewSQLSExecutorB(data)
	if err != nil {
		t.Fatal(err.Error())
	}
	sqlsExecutor.Exec(conn)

	err = ioutil.WriteFile("./results_"+sqlFile+".txt", []byte(sqlsExecutor.ToString()), 0777)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestSQLSExecutor_Exec(t *testing.T) {
	testSQLSExecutor_ExecCommon(t, testsqls.SQLFileTest)
}

func TestSQLSExecutor_Exec2(t *testing.T) {
	testSQLSExecutor_ExecCommon(t, testsqls.SQLFileWindow)
}