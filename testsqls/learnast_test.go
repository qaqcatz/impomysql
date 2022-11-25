package testsqls

import (
	"testing"
)

func testLearnASTCommon(t *testing.T, sql string) {
	if err := InitTableCOMPANY(""); err != nil {
		t.Fatalf("%+v", err)
	}

	if err := SQLExec(sql, ""); err != nil {
		t.Fatalf("%+v", err)
	}
	if sql, err := learnAST(sql); err != nil {
		t.Fatalf("%+v", err)
	} else {
		t.Log(sql)
	}
}

func TestLearnASTAGG(t *testing.T) {
	testLearnASTCommon(t, SQLAGG);
}

func TestLearnASTWindow(t *testing.T) {
	testLearnASTCommon(t, SQLWindow);
}

func TestLearnASTSelectValue(t *testing.T) {
	testLearnASTCommon(t, SQLSelectValue);
}
func TestLearnASTSelectValue2(t *testing.T) {
	testLearnASTCommon(t, SQLSelectValue2);
}

func TestLearnASTSelectValue3(t *testing.T) {
	testLearnASTCommon(t, SQLSelectValue3);
}

func TestLearnASTSubQuery(t *testing.T) {
	testLearnASTCommon(t, SQLSubQuery);
}

func TestLearnASTSubQuery2(t *testing.T) {
	testLearnASTCommon(t, SQLSubQuery2);
}

func TestLearnASTSubQuery3(t *testing.T) {
	testLearnASTCommon(t, SQLSubQuery3);
}

func TestLearnASTSubQuery4(t *testing.T) {
	testLearnASTCommon(t, SQLSubQuery4);
}

func TestLearnASTSubQuery5(t *testing.T) {
	testLearnASTCommon(t, SQLSubQuery5);
}

func TestLearnASTJOIN(t *testing.T) {
	testLearnASTCommon(t, SQLJOIN);
}

func TestLearnASTJOIN2(t *testing.T) {
	testLearnASTCommon(t, SQLJOIN2);
}

func TestLearnASTJOIN3(t *testing.T) {
	testLearnASTCommon(t, SQLJOIN3);
}

func TestLearnASTJOIN4(t *testing.T) {
	testLearnASTCommon(t, SQLJOIN4);
}

func TestLearnASTJOIN5(t *testing.T) {
	testLearnASTCommon(t, SQLJOIN5);
}

func TestLearnASTJOIN6(t *testing.T) {
	testLearnASTCommon(t, SQLJOIN6);
}

func TestLearnASTLIMIT(t *testing.T) {
	testLearnASTCommon(t, SQLLIMIT);
}

func TestLearnASTLIMIT2(t *testing.T) {
	testLearnASTCommon(t, SQLLIMIT2);
}

func TestLearnASTUNION(t *testing.T) {
	testLearnASTCommon(t, SQLUNION);
}

func TestLearnASTUNION2(t *testing.T) {
	testLearnASTCommon(t, SQLUNION2);
}

func TestLearnASTWITH(t *testing.T) {
	testLearnASTCommon(t, SQLWITH);
}

func TestLearnASTWITH2(t *testing.T) {
	testLearnASTCommon(t, SQLWITH2);
}

func TestLearnASTIN(t *testing.T) {
	testLearnASTCommon(t, SQLIN);
}

func TestLearnASTWHERE(t *testing.T) {
	testLearnASTCommon(t, SQLWHERE);
}

func TestLearnASTLIKE(t *testing.T) {
	testLearnASTCommon(t, SQLLIKE);
}

func TestLearnASTRegExp(t *testing.T) {
	testLearnASTCommon(t, SQLRegExp);
}