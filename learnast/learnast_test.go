package learnast

import (
	"github.com/qaqcatz/impomysql/testsqls"
	"testing"
)

func testLearnASTCommon(t *testing.T, sql string) {
	if err := testsqls.EnsureDBTEST(); err != nil {
		t.Fatal(err.Error())
	}
	if err := testsqls.InitTableCOMPANY(); err != nil {
		t.Fatal(err.Error())
	}

	if err := testsqls.SQLExec(sql); err != nil {
		t.Fatal(err.Error())
	}
	if sql, err := learnAST(sql); err != nil {
		t.Fatal(err.Error())
	} else {
		t.Log(sql)
	}
}

func TestLearnASTAGG(t *testing.T) {
	testLearnASTCommon(t, testsqls.SQLAGG);
}

func TestLearnASTWindow(t *testing.T) {
	testLearnASTCommon(t, testsqls.SQLWindow);
}

func TestLearnASTSelectValue(t *testing.T) {
	testLearnASTCommon(t, testsqls.SQLSelectValue);
}
func TestLearnASTSelectValue2(t *testing.T) {
	testLearnASTCommon(t, testsqls.SQLSelectValue2);
}

func TestLearnASTSelectValue3(t *testing.T) {
	testLearnASTCommon(t, testsqls.SQLSelectValue3);
}

func TestLearnASTSubQuery(t *testing.T) {
	testLearnASTCommon(t, testsqls.SQLSubQuery);
}

func TestLearnASTSubQuery2(t *testing.T) {
	testLearnASTCommon(t, testsqls.SQLSubQuery2);
}

func TestLearnASTSubQuery3(t *testing.T) {
	testLearnASTCommon(t, testsqls.SQLSubQuery3);
}

func TestLearnASTSubQuery4(t *testing.T) {
	testLearnASTCommon(t, testsqls.SQLSubQuery4);
}

func TestLearnASTSubQuery5(t *testing.T) {
	testLearnASTCommon(t, testsqls.SQLSubQuery5);
}

func TestLearnASTJOIN(t *testing.T) {
	testLearnASTCommon(t, testsqls.SQLJOIN);
}

func TestLearnASTJOIN2(t *testing.T) {
	testLearnASTCommon(t, testsqls.SQLJOIN2);
}

func TestLearnASTJOIN3(t *testing.T) {
	testLearnASTCommon(t, testsqls.SQLJOIN3);
}

func TestLearnASTJOIN4(t *testing.T) {
	testLearnASTCommon(t, testsqls.SQLJOIN4);
}

func TestLearnASTJOIN5(t *testing.T) {
	testLearnASTCommon(t, testsqls.SQLJOIN5);
}

func TestLearnASTJOIN6(t *testing.T) {
	testLearnASTCommon(t, testsqls.SQLJOIN6);
}

func TestLearnASTLIMIT(t *testing.T) {
	testLearnASTCommon(t, testsqls.SQLLIMIT);
}

func TestLearnASTLIMIT2(t *testing.T) {
	testLearnASTCommon(t, testsqls.SQLLIMIT2);
}

func TestLearnASTUNION(t *testing.T) {
	testLearnASTCommon(t, testsqls.SQLUNION);
}

func TestLearnASTUNION2(t *testing.T) {
	testLearnASTCommon(t, testsqls.SQLUNION2);
}

func TestLearnASTWITH(t *testing.T) {
	testLearnASTCommon(t, testsqls.SQLWITH);
}

func TestLearnASTWITH2(t *testing.T) {
	testLearnASTCommon(t, testsqls.SQLWITH2);
}

func TestLearnASTIN(t *testing.T) {
	testLearnASTCommon(t, testsqls.SQLIN);
}

func TestLearnASTWHERE(t *testing.T) {
	testLearnASTCommon(t, testsqls.SQLWHERE);
}

func TestLearnASTLIKE(t *testing.T) {
	testLearnASTCommon(t, testsqls.SQLLIKE);
}

func TestLearnASTRegExp(t *testing.T) {
	testLearnASTCommon(t, testsqls.SQLRegExp);
}