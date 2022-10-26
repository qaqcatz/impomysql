package stage1

import (
	"github.com/qaqcatz/impomysql/sqlsexecutor"
	"github.com/qaqcatz/impomysql/testsqls"
	"strconv"
	"testing"
)

func testInitCommon(t *testing.T, sql string) {
	if err := testsqls.EnsureDBTEST(); err != nil {
		t.Fatal(err.Error())
	}
	if err := testsqls.InitTableCOMPANY(); err != nil {
		t.Fatal(err.Error())
	}

	if err := testsqls.SQLExec(sql); err != nil {
		t.Fatal(err.Error())
	}
	if initResult := Init(sql); initResult.Err != nil {
		t.Fatal(initResult.Err.Error())
	} else {
		if err := testsqls.SQLExec(initResult.InitSql); err != nil {
			t.Fatal(err.Error())
		}
	}
}

func TestInitAGG(t *testing.T) {
	testInitCommon(t, testsqls.SQLAGG)
}

func TestInitWindow(t *testing.T) {
	testInitCommon(t, testsqls.SQLWindow)
}

func TestInitJOIN(t *testing.T) {
	testInitCommon(t, testsqls.SQLJOIN2)
	testInitCommon(t, testsqls.SQLJOIN3)
	testInitCommon(t, testsqls.SQLJOIN6)
}

func TestInitLIMIT(t *testing.T) {
	testInitCommon(t, testsqls.SQLLIMIT)
	testInitCommon(t, testsqls.SQLLIMIT2)
}

func testInitCommon2(t *testing.T, sqlFileName string, oracle int) {
	err := testsqls.InitDBTEST()
	if err != nil {
		t.Fatal(err.Error())
	}

	conn, err := testsqls.GetConnector()
	if err != nil {
		t.Fatal(err.Error())
	}
	data, err, _ := testsqls.ReadSQLFile(sqlFileName)
	if err != nil {
		t.Log(err.Error())
	}
	sqlsExecutor, err := sqlsexecutor.NewSQLSExecutorB(data)
	if err != nil {
		t.Fatal(err.Error())
	}
	sqlsExecutor.Exec(conn)

	// t.Log(sqlsExecutor.ToString())

	passedNum := 0
	failedNum := 0
	passedStr := ""
	failedStr := ""
	for i, result := range sqlsExecutor.Results {
		if result.Err != nil {
			continue
		}

		initResult := Init(sqlsExecutor.ASTs[i].Text())
		if initResult.Err != nil {
			failedNum += 1
			failedStr += "========================================\n"
			failedStr += "[sql " + strconv.Itoa(i) + "] " + sqlsExecutor.ASTs[i].Text() + "\n"
			failedStr += "@@@@@@@@@@Init failed!@@@@@@@@@@\n" + initResult.Err.Error() + "\n"
			continue
		}

		resultm := conn.ExecSQL(initResult.InitSql)
		if resultm.Err != nil {
			failedNum += 1
			failedStr += "========================================\n"
			failedStr += "[sql " + strconv.Itoa(i) + "] " + sqlsExecutor.ASTs[i].Text() + "\n"
			failedStr += "[Init] " + initResult.InitSql + "\n"
			failedStr += resultm.ToString() + "\n"
		} else {
			passedNum += 1
			passedStr += "========================================\n"
			passedStr += "[sql " + strconv.Itoa(i) + "] " + sqlsExecutor.ASTs[i].Text() + "\n"
			passedStr += "[Init] " + initResult.InitSql + "\n"
			passedStr += resultm.ToString() + "\n"
		}
	}
	passedStr = strconv.Itoa(passedNum) + "\n" + passedStr
	failedStr = strconv.Itoa(failedNum) + "\n" + failedStr

	//t.Log(passedStr)
	//t.Log(failedStr)

	if passedNum != oracle {
		t.Fatal("passedNum != oracle: [passedNum]", passedNum, "[oracle]", oracle)
	}
}

func TestInitFileAgg(t *testing.T) {
	testInitCommon2(t, testsqls.SQLFileAgg, 100)
}

func TestInitFileWindow(t *testing.T) {
	testInitCommon2(t, testsqls.SQLFileWindow, 10)
}