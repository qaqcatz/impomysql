package randgen

import (
	"encoding/json"
	"github.com/qaqcatz/impomysql/sqlsexecutor"
	"github.com/qaqcatz/impomysql/testsqls"
	"io/ioutil"
	"testing"
	"time"
)

func TestRandGen(t *testing.T) {
	//sqls, err := RandGen(ZZTest, YYTest, 100, 123456)
	sqls, err := RandGen(ZZTest, YYTest, 100, time.Now().UnixNano())
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

func testRandGenCommon(t *testing.T, zzFilePath string, yyFilePath string, queriesNum int, seed int64, log bool) {
	//sqls, err := RandGen(ZZTest, YYTest, 100, 123456)
	sqls, err := RandGen(zzFilePath, yyFilePath, queriesNum, seed)
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

	if log {
		err = ioutil.WriteFile("./results_"+ZZDefault+".txt", []byte(sqlsExecutor1.ToString()), 0777)
		if err != nil {
			t.Fatal(err.Error())
		}
	} else {
		t.Log(ZZDefault, ":\n", sqlsExecutor1.ToShortString())
	}

	sqlsExecutor2, err := sqlsexecutor.NewSQLSExecutorS(sqls.RandSQLs)
	if err != nil {
		t.Fatal(err.Error())
	}
	sqlsExecutor2.Exec(conn)

	if log {
		err = ioutil.WriteFile("./results_"+YYDefault+".txt", []byte(sqlsExecutor2.ToString()), 0777)
		if err != nil {
			t.Fatal(err.Error())
		}
	} else {
		t.Log(YYDefault, ":\n", sqlsExecutor2.ToShortString())
	}
}

func TestRandGen2(t *testing.T) {
	testRandGenCommon(t, ZZTest, YYTest, 100, time.Now().UnixNano(), true)
}

func TestRandGen3(t *testing.T) {
	testRandGenCommon(t, ZZTest, YYTest, 100, time.Now().UnixNano(), false)
}

// 3MB Memory
func TestRandGen4(t *testing.T) {
	testRandGenCommon(t, ZZTest, YYTest, 10000, time.Now().UnixNano(), false)
}