package stage1

import (
	"github.com/qaqcatz/impomysql/testsqls"
	"testing"
)

func testInitCommon(t *testing.T, sql string) {
	if err := testsqls.InitTableCOMPANY(""); err != nil {
		t.Fatalf("%+v", err)
	}

	if err := testsqls.SQLExec(sql, ""); err != nil {
		t.Fatalf("%+v", err)
	}
	if initResult := Init(sql); initResult.Err != nil {
		t.Fatalf("%+v", initResult.Err)
	} else {
		if err := testsqls.SQLExec(initResult.InitSql, ""); err != nil {
			t.Fatalf("%+v", err)
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