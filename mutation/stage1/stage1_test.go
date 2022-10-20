package stage1

import (
	"errors"
	"github.com/qaqcatz/impomysql/sqlsexecutor"
	"github.com/qaqcatz/impomysql/testsqls"
	"io/ioutil"
	"path"
	"runtime"
	"strconv"
	"testing"
)

// getPackagePath: get the package actual path, then you can read files under the path.
func getPackagePath() (string, error) {
	if _, file, _, ok := runtime.Caller(0); !ok {
		return "", errors.New("PackagePath: runtime.Caller(0) error ")
	} else {
		return path.Join(file, "../"), nil
	}
}

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
	if sqlm, err := Init(sql); err != nil {
		t.Fatal(err.Error())
	} else {
		if err := testsqls.SQLExec(sqlm); err != nil {
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

	packagePath, err := getPackagePath()
	if err != nil {
		t.Fatal(err.Error())
	}

	err = ioutil.WriteFile(path.Join(packagePath, "results_"+sqlFileName+".txt"), []byte(sqlsExecutor.ToString()), 0777)
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

		sqlm, err := Init(sqlsExecutor.ASTs[i].Text())
		if err != nil {
			failedNum += 1
			failedStr += "========================================\n"
			failedStr += "[sql " + strconv.Itoa(i) + "] " + sqlsExecutor.ASTs[i].Text() + "\n"
			failedStr += "@@@@@@@@@@Init failed!@@@@@@@@@@\n" + err.Error() + "\n"
			continue
		}

		resultm := conn.ExecSQL(sqlm)
		if resultm.Err != nil {
			failedNum += 1
			failedStr += "========================================\n"
			failedStr += "[sql " + strconv.Itoa(i) + "] " + sqlsExecutor.ASTs[i].Text() + "\n"
			failedStr += "[Init] " + sqlm + "\n"
			failedStr += resultm.ToString() + "\n"
		} else {
			passedNum += 1
			passedStr += "========================================\n"
			passedStr += "[sql " + strconv.Itoa(i) + "] " + sqlsExecutor.ASTs[i].Text() + "\n"
			passedStr += "[Init] " + sqlm + "\n"
			passedStr += resultm.ToString() + "\n"
		}
	}
	passedStr = strconv.Itoa(passedNum) + "\n" + passedStr
	failedStr = strconv.Itoa(failedNum) + "\n" + failedStr
	err = ioutil.WriteFile(path.Join(packagePath, "results_"+sqlFileName+"_pass.txt"), []byte(passedStr), 0777)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = ioutil.WriteFile(path.Join(packagePath, "results_"+sqlFileName+"_fail.txt"), []byte(failedStr), 0777)
	if err != nil {
		t.Fatal(err.Error())
	}

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