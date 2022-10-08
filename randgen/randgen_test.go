package randgen

import (
	"encoding/json"
	"github.com/qaqcatz/impomysql/sqlsexecutor"
	"github.com/qaqcatz/impomysql/testsqls"
	"io/ioutil"
	"testing"
)

func TestRandGen(t *testing.T) {
	sqls, err := RandGen(ZZDefault, YYDefault, 10, 123456)
	if err != nil {
		t.Fatal(err.Error())
	}
	data, err := json.Marshal(sqls)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = ioutil.WriteFile("./test.json", data, 0777)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestRandGen2(t *testing.T) {
	sqls, err := RandGen(ZZDefault, YYDefault, 10, 123456)
	if err != nil {
		t.Fatal(err.Error())
	}

	err = testsqls.InitDBTEST()
	if err != nil {
		t.Fatal(err.Error())
	}
	conn, err := testsqls.GetConnector()
	if err != nil {
		t.Fatal(err.Error())
	}

	sqlsExecutor1, err := sqlsexecutor.NewSQLSExecutorS(sqls.DDLs)
	if err != nil {
		t.Fatal(err.Error())
	}
	sqlsExecutor1.Exec(conn)
	err = ioutil.WriteFile("./results_"+ZZDefault+".txt", []byte(sqlsExecutor1.ToString()), 0777)
	if err != nil {
		t.Fatal(err.Error())
	}

	sqlsExecutor2, err := sqlsexecutor.NewSQLSExecutorS(sqls.RandSQLs)
	if err != nil {
		t.Fatal(err.Error())
	}
	sqlsExecutor2.Exec(conn)
	err = ioutil.WriteFile("./results_"+YYDefault+".txt", []byte(sqlsExecutor2.ToString()), 0777)
	if err != nil {
		t.Fatal(err.Error())
	}
}
