package stage1

import (
	"github.com/qaqcatz/impomysql/sqlsexecutor"
	"github.com/qaqcatz/impomysql/testsqls"
	"io/ioutil"
	"strconv"
	"testing"
)

// init
func TestStage1(t *testing.T) {
	if err := testsqls.EnsureDBTEST(); err != nil {
		t.Fatal(err.Error())
	}
	if err := testsqls.InitTableCOMPANY(); err != nil {
		t.Fatal(err.Error())
	}
}

func testStage1Common(t *testing.T, sql string) {
	if err := testsqls.SQLExec(sql); err != nil {
		t.Fatal(err.Error())
	}
	if sqlm, err := Stage1(sql); err != nil {
		t.Fatal(err.Error())
	} else {
		if err := testsqls.SQLExec(sqlm); err != nil {
			t.Fatal(err.Error())
		}
	}
}

func TestStage12(t *testing.T) {
	testStage1Common(t, testsqls.SQLAGG)
}

func TestStage13(t *testing.T) {
	testStage1Common(t, testsqls.SQLWindow)
}

func testStage1Common2(t *testing.T, sqlFileName string) {
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

	err = ioutil.WriteFile("./results_"+sqlFileName+".txt", []byte(sqlsExecutor.ToString()), 0777)
	if err != nil {
		t.Fatal(err.Error())
	}

	passedNum := 0
	failedNum := 0
	passedStr := ""
	failedStr := ""
	for i, result := range sqlsExecutor.Results {
		if result.Err != nil {
			continue
		}

		sqlm, err := Stage1(sqlsExecutor.ASTs[i].Text())
		if err != nil {
			failedNum += 1
			failedStr += "========================================\n"
			failedStr += "[sql " + strconv.Itoa(i) + "] " + sqlsExecutor.ASTs[i].Text() + "\n"
			failedStr += "@@@@@@@@@@Stage1 failed!@@@@@@@@@@\n" + err.Error() + "\n"
			continue
		}

		resultm := conn.ExecSQL(sqlm)
		if resultm.Err != nil {
			failedNum += 1
			failedStr += "========================================\n"
			failedStr += "[sql " + strconv.Itoa(i) + "] " + sqlsExecutor.ASTs[i].Text() + "\n"
			failedStr += "[stage1] " + sqlm + "\n"
			failedStr += resultm.ToString() + "\n"
		} else {
			passedNum += 1
			passedStr += "========================================\n"
			passedStr += "[sql " + strconv.Itoa(i) + "] " + sqlsExecutor.ASTs[i].Text() + "\n"
			passedStr += "[stage1] " + sqlm + "\n"
			passedStr += resultm.ToString() + "\n"
		}
	}
	passedStr = strconv.Itoa(passedNum) + "\n" + passedStr
	failedStr = strconv.Itoa(failedNum) + "\n" + failedStr
	err = ioutil.WriteFile("./results_"+sqlFileName+"_pass.txt", []byte(passedStr), 0777)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = ioutil.WriteFile("./results_"+sqlFileName+"_fail.txt", []byte(failedStr), 0777)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestStage14(t *testing.T) {
	testStage1Common2(t, testsqls.SQLFileAgg)
}

func TestStage15(t *testing.T) {
	testStage1Common2(t, testsqls.SQLFileWindow)
}