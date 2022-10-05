package sqlfexecutor

import (
	"github.com/qaqcatz/IMPOMySQL/IMPOMySQL/testsqls"
	"io/ioutil"
	"strconv"
	"testing"
)

func TestNewSQLFExecutor(t *testing.T) {
	sqlFExecutor, err := NewSQLFExecutor("./test.sql")
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log("|ASTs|:", len(sqlFExecutor.ASTs))
	t.Log("Read Time:", sqlFExecutor.ReadTime)
	t.Log("Parse Time:", sqlFExecutor.ParseTime)
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
	sqlFExecutor, err := NewSQLFExecutor("./test.sql")
	if err != nil {
		t.Fatal(err.Error())
	}
	sqlFExecutor.Exec(conn)

	str := ""
	for i, result := range sqlFExecutor.Results {
		str += "==================================================\n"
		str += "[sql "+strconv.Itoa(i)+"]: " + sqlFExecutor.ASTs[i].Text() + "\n"
		str += "[result "+strconv.Itoa(i)+"]: " + result.ToString() + "\n"
		str += "==================================================\n\n"
	}
	err = ioutil.WriteFile("./results.txt", []byte(str), 0777)
	if err != nil {
		t.Fatal(err.Error())
	}

	t.Log("|ASTs|:", len(sqlFExecutor.ASTs))
	t.Log("Read Time:", sqlFExecutor.ReadTime)
	t.Log("Parse Time:", sqlFExecutor.ParseTime)
	t.Log("Exec Time:", sqlFExecutor.ExecuteTime)
	t.Log("Passed SQL Num:", sqlFExecutor.PassedSQLNum)
	t.Log("Failed SQL Num:", sqlFExecutor.FailedSQLNum)
	t.Log("Results: ./results.txt")
}