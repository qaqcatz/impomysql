package sqlfexecutor

import (
	"github.com/qaqcatz/IMPOMySQL/IMPOMySQL/testsqls"
	"io/ioutil"
	"testing"
)

func TestNewSQLFExecutor(t *testing.T) {
	data, err, _ := testsqls.ReadSQLFile(testsqls.SQLFileTest)
	if err != nil {
		t.Fatal(err.Error())
	}
	sqlFExecutor, err := NewSQLFExecutorB(data)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(sqlFExecutor.ToString())
}

func TestSQLFExecutor_Exec(t *testing.T) {
	err := testsqls.InitDBTEST()
	if err != nil {
		t.Fatal(err.Error())
	}
	conn, err := testsqls.GetConnector()
	if err != nil {
		t.Fatal(err.Error())
	}
	data, err, _ := testsqls.ReadSQLFile(testsqls.SQLFileTest)
	if err != nil {
		t.Log(err.Error())
	}
	sqlFExecutor, err := NewSQLFExecutorB(data)
	if err != nil {
		t.Fatal(err.Error())
	}
	sqlFExecutor.Exec(conn)

	err = ioutil.WriteFile("./results.txt", []byte(sqlFExecutor.ToString()), 0777)
	if err != nil {
		t.Fatal(err.Error())
	}
}