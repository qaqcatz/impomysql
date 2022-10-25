package randgen

import (
	"encoding/json"
	"github.com/qaqcatz/impomysql/sqlsexecutor"
	"github.com/qaqcatz/impomysql/testsqls"
	"io/ioutil"
	"path"
	"testing"
	"time"
)

func TestRandGenJson(t *testing.T) {
	//sqls, err := RandGen(ZZTest, YYImpo, 100, 123456)
	sqls, err := RandGen(ZZTest, YYImpo, 100, time.Now().UnixNano())
	if err != nil {
		t.Fatal(err.Error())
	}
	data, err := json.Marshal(sqls)
	if err != nil {
		t.Fatal(err.Error())
	}
	packagePath, err := getPackagePath()
	if err != nil {
		t.Fatal(err.Error())
	}
	err = ioutil.WriteFile(path.Join(packagePath, "test.json"), data, 0777)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func testRandGenCommon(t *testing.T, zzFilePath string, yyFilePath string, queriesNum int, seed int64, name string, log bool) {
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

	packagePath, err := getPackagePath()
	if err != nil {
		t.Fatal(err.Error())
	}

	if log {
		err = ioutil.WriteFile(path.Join(packagePath, name+"_zz.txt"), []byte(sqlsExecutor1.ToString()), 0777)
		if err != nil {
			t.Fatal(err.Error())
		}
	} else {
		t.Log(zzFilePath, ":\n", sqlsExecutor1.ToShortString())
	}

	sqlsExecutor2, err := sqlsexecutor.NewSQLSExecutorS(sqls.RandSQLs)
	if err != nil {
		t.Fatal(err.Error())
	}
	sqlsExecutor2.Exec(conn)

	if log {
		err = ioutil.WriteFile(path.Join(packagePath, name+"_yy.txt"), []byte(sqlsExecutor2.ToString()), 0777)
		if err != nil {
			t.Fatal(err.Error())
		}
	} else {
		t.Log(yyFilePath, ":\n", sqlsExecutor2.ToShortString())
	}
}

func TestRandGenRd100Log(t *testing.T) {
	testRandGenCommon(t, ZZTest, YYImpo, 100, time.Now().UnixNano(), "fix100", true)
}

func TestRandGenRd100(t *testing.T) {
	testRandGenCommon(t, ZZTest, YYImpo, 100, time.Now().UnixNano(), "", false)
}

//3MB Memory
//func TestRandGenRd10000(t *testing.T) {
//	testRandGenCommon(t, ZZTest, YYImpo, 10000, time.Now().UnixNano(), "",false)
//}