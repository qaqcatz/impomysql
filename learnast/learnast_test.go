package learnast

import (
	"github.com/qaqcatz/impomysql/testsqls"
	"testing"
)

// init
func TestLearnAST(t *testing.T) {
	if err := testsqls.EnsureDBTEST(); err != nil {
		t.Fatal(err.Error())
	}
	if err := testsqls.InitTableCOMPANY(); err != nil {
		t.Fatal(err.Error())
	}
}

func testLearnASTCommon(t *testing.T, sql string) {
	if err := testsqls.SQLExec(sql); err != nil {
		t.Fatal(err.Error())
	}
	if sql, err := learnAST(sql); err != nil {
		t.Fatal(err.Error())
	} else {
		t.Log(sql)
	}
}

func TestLearnAST2(t *testing.T) {
	testLearnASTCommon(t, testsqls.SQLAGG);
}

func TestLearnAST3(t *testing.T) {
	testLearnASTCommon(t, testsqls.SQLWindow);
}

func TestLearnAST4(t *testing.T) {
	testLearnASTCommon(t, testsqls.SQLSelectValue);
}
func TestLearnAST5(t *testing.T) {
	testLearnASTCommon(t, testsqls.SQLSelectValue2);
}

func TestLearnAST6(t *testing.T) {
	testLearnASTCommon(t, testsqls.SQLSelectValue3);
}

func TestLearnAST7(t *testing.T) {
	testLearnASTCommon(t, testsqls.SQLSubQuery);
}

func TestLearnAST8(t *testing.T) {
	testLearnASTCommon(t, testsqls.SQLSubQuery2);
}

func TestLearnAST9(t *testing.T) {
	testLearnASTCommon(t, testsqls.SQLSubQuery3);
}

func TestLearnAST10(t *testing.T) {
	testLearnASTCommon(t, testsqls.SQLSubQuery4);
}

func TestLearnAST11(t *testing.T) {
	testLearnASTCommon(t, testsqls.SQLJOIN);
}

func TestLearnAST12(t *testing.T) {
	testLearnASTCommon(t, testsqls.SQLJOIN2);
}

func TestLearnAST13(t *testing.T) {
	testLearnASTCommon(t, testsqls.SQLJOIN3);
}

func TestLearnAST14(t *testing.T) {
	testLearnASTCommon(t, testsqls.SQLJOIN4);
}

func TestLearnAST15(t *testing.T) {
	testLearnASTCommon(t, testsqls.SQLJOIN5);
}

func TestLearnAST16(t *testing.T) {
	testLearnASTCommon(t, testsqls.SQLJOIN6);
}

func TestLearnAST17(t *testing.T) {
	testLearnASTCommon(t, testsqls.SQLLIMIT);
}

func TestLearnAST18(t *testing.T) {
	testLearnASTCommon(t, testsqls.SQLLIMIT2);
}

func TestLearnAST19(t *testing.T) {
	testLearnASTCommon(t, testsqls.SQLUNION);
}