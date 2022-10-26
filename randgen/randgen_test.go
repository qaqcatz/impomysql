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
	config := &Config{
		ZZFilePath: testsqls.GetTestZZPath(),
		YYFilePath: testsqls.GetTestYYPath(),
		QueriesNum: 100,
		Seed:       time.Now().UnixNano(),
	}
	sqls := RandGen(config)
	if sqls.Err != nil {
		t.Fatal(sqls.Err.Error())
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

func testRandGenCommon(t *testing.T, config *Config, name string, log bool) {
	sqls := RandGen(config)
	if sqls.Err != nil {
		t.Fatal(sqls.Err.Error())
	}

	err := testsqls.InitDBTEST()
	if err != nil {
		t.Fatal(err.Error())
	}
	conn, err := testsqls.GetConnector()
	if err != nil {
		t.Fatal(err.Error())
	}

	sqlsExecutor1 := sqlsexecutor.NewSQLSExecutorS(sqls.DDLs)
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
		t.Log(config.ZZFilePath, ":\n", sqlsExecutor1.ToShortString())
	}

	sqlsExecutor2 := sqlsexecutor.NewSQLSExecutorS(sqls.RandSQLs)
	sqlsExecutor2.Exec(conn)

	if log {
		err = ioutil.WriteFile(path.Join(packagePath, name+"_yy.txt"), []byte(sqlsExecutor2.ToString()), 0777)
		if err != nil {
			t.Fatal(err.Error())
		}
	} else {
		t.Log(config.YYFilePath, ":\n", sqlsExecutor2.ToShortString())
	}
}

func TestRandGenRd100Log(t *testing.T) {
	config := &Config{
		ZZFilePath: testsqls.GetTestZZPath(),
		YYFilePath: testsqls.GetTestYYPath(),
		QueriesNum: 100,
		Seed:       time.Now().UnixNano(),
	}
	testRandGenCommon(t, config, "fix100", true)
}

func TestRandGenRd100(t *testing.T) {
	config := &Config{
		ZZFilePath: testsqls.GetTestZZPath(),
		YYFilePath: testsqls.GetTestYYPath(),
		QueriesNum: 100,
		Seed:       time.Now().UnixNano(),
	}
	testRandGenCommon(t, config, "", false)
}

//3MB Memory
//func TestRandGenRd10000(t *testing.T) {
//	config := &Config{
//		ZZFilePath: testsqls.GetTestZZPath(),
//		YYFilePath: testsqls.GetTestYYPath(),
//		QueriesNum: 10000,
//		Seed: time.Now().UnixNano(),
//	}
//	testRandGenCommon(t, config, "",false)
//}