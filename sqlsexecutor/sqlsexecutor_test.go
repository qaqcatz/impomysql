package sqlsexecutor

import (
	"github.com/qaqcatz/impomysql/testsqls"
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

func testNewSQLSExecutorCommon(t *testing.T, sqlFile string, oracle int) {
	data, err, _ := testsqls.ReadSQLFile(sqlFile)
	if err != nil {
		t.Fatal(err.Error())
	}
	sqlsExecutor, err := NewSQLSExecutorB(data)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(sqlsExecutor.ToString())

	if len(sqlsExecutor.ParseErrs) != oracle {
		t.Fatal("len(sqlsExecutor.ParseErrs) != oracle: ", len(sqlsExecutor.ParseErrs), oracle)
	}
}

func TestNewSQLSExecutorTest(t *testing.T) {
	testNewSQLSExecutorCommon(t, testsqls.SQLFileTest, 0)
}

func TestNewSQLSExecutorWindow(t *testing.T) {
	testNewSQLSExecutorCommon(t, testsqls.SQLFileWindow, 49)
}

func testSQLSExecutor_ExecCommon(t *testing.T, sqlFile string, oracle int) {
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

	//t.Log(sqlsExecutor.ToString())

	if sqlsExecutor.PassedSQLNum != oracle {
		t.Fatal("sqlsExecutor.PassedSQLNum != oracle: ", sqlsExecutor.PassedSQLNum, oracle)
	}
}

func TestSQLSExecutor_ExecTest(t *testing.T) {
	testSQLSExecutor_ExecCommon(t, testsqls.SQLFileTest, 136)
}

func TestSQLSExecutor_ExecWindow(t *testing.T) {
	testSQLSExecutor_ExecCommon(t, testsqls.SQLFileWindow, 46)
}