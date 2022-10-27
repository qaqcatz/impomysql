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

	if err := testsqls.SQLExecS(sql); err != nil {
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

func TestLearnASTBug(t *testing.T) {
	//testLearnASTCommon(t, "SELECT 1")
	testLearnASTCommon(t, "WITH `MYWITH` AS (SELECT (DATE_ADD(`f6`, INTERVAL 1 MICROSECOND)) AS `f1`,(COERCIBILITY(NULL)%`f4`) AS `f2`,(UCASE(`f4`) DIV `f6`>>3) AS `f3` FROM (SELECT `col_bigint_undef_signed` AS `f4`,`col_bigint_key_unsigned` AS `f5`,`col_double_key_unsigned` AS `f6` FROM `table_3_utf8_2` USE INDEX (`col_float_key_unsigned`)) AS `t1`) SELECT * FROM `MYWITH`")

	//testLearnASTCommon(t, testsqls.SQLEX2);
	//testLearnASTCommon(t, "select exists ("+testsqls.SQLEX2+")");

	//testLearnASTCommon(t, "select * from COMPANY where 9223372036854775807 + 1 > 1;");

	//testLearnASTCommon(t, "select 9223372036854775807 + 1 > 1");
	//testLearnASTCommon(t, "select exists (select 9223372036854775807 + 1 > 1)");

	//testLearnASTCommon(t, "select exists (select * from COMPANY WHERE ID = 0)");
}