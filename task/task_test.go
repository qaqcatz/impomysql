package task

import (
	"github.com/qaqcatz/impomysql/randgen"
	"github.com/qaqcatz/impomysql/testsqls"
	"testing"
)

func testRunCommon(t *testing.T, queriesNum int, seed int64) {
	err := testsqls.InitDBTEST()
	if err != nil {
		t.Fatal(err.Error())
	}

	conn, err := testsqls.GetConnector()
	if err != nil {
		t.Fatal(err.Error())
	}
	randGenConfig := &randgen.Config {
		ZZFilePath: testsqls.GetTestZZPath(),
		YYFilePath: testsqls.GetTestYYPath(),
		QueriesNum: queriesNum,
		Seed:       seed,
	}
	config := &Config{
		Conn: conn,
		RandGenConfig: randGenConfig,
	}
	result := Run(config)
	if result.RandGenRes.Err != nil {
		t.Fatal(result.RandGenRes.Err)
	}
	t.Log(result.ToShortString())

	t.Log("DDL")
	t.Log("==================================================")
	for i, ddl := range result.RandGenRes.DDLs {
		t.Log("ddl", i)
		t.Log(ddl)
	}
	for i, bug := range result.ImpoBugs {
		t.Log("[bug]", i, "==================================================")
		t.Log(bug.ToString())
	}
}

//func TestRun(t *testing.T) {
//	testRunCommon(t, 100, 123456)
//}
//
//func TestRun10000(t *testing.T) {
//	testRunCommon(t, 10000, 654321)
//}

func testRunOtherCommon(t *testing.T, queriesNum int, seed int64, DBMS string) {
	err := testsqls.InitOtherDBTEST(DBMS)
	if err != nil {
		t.Fatal(err.Error())
	}

	conn, err := testsqls.GetOtherDBConnector(DBMS)
	if err != nil {
		t.Fatal(err.Error())
	}
	randGenConfig := &randgen.Config {
		ZZFilePath: testsqls.GetTestZZPath(),
		YYFilePath: testsqls.GetTestYYPath(),
		QueriesNum: queriesNum,
		Seed:       seed,
	}
	config := &Config{
		Conn: conn,
		RandGenConfig: randGenConfig,
	}
	result := Run(config)
	if result.RandGenRes.Err != nil {
		t.Fatal(result.RandGenRes.Err)
	}
	t.Log(result.ToShortString())

	t.Log("DDL")
	t.Log("==================================================")
	for i, ddl := range result.RandGenRes.DDLs {
		t.Log("ddl", i)
		t.Log(ddl)
	}
	for i, bug := range result.ImpoBugs {
		t.Log("[bug]", i, "==================================================")
		t.Log(bug.ToString())
	}
}

//func TestRunOtherMariaDB(t *testing.T) {
//	testRunOtherCommon(t, 100, 123456, testsqls.MariaDB)
//}
//
//func TestRunOtherTiDB(t *testing.T) {
//	testRunOtherCommon(t, 100, 123456, testsqls.TiDB)
//}
//
//func TestRunOtherOceanBase(t *testing.T) {
//	testRunOtherCommon(t, 100, 123456, testsqls.OceanBase)
//}

func TestRun10000OtherMariaDB(t *testing.T) {
	testRunOtherCommon(t, 10000, 654321, testsqls.MariaDB)
}

func TestRun10000OtherTiDB(t *testing.T) {
	testRunOtherCommon(t, 10000, 654321, testsqls.TiDB)
}

func TestRun10000OtherOceanBase(t *testing.T) {
	testRunOtherCommon(t, 10000, 654321, testsqls.OceanBase)
}