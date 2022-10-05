package stage1

import (
	"github.com/qaqcatz/IMPOMySQL/IMPOMySQL/testsqls"
	"testing"
)

func TestStage1(t *testing.T) {
	if err := testsqls.InitDBTEST(); err != nil {
		t.Fatal(err.Error())
	}
	if err := testsqls.InitTableCOMPANY(); err != nil {
		t.Fatal(err.Error())
	}
}

func TestStage12(t *testing.T) {
	if err := testsqls.SQLExec(testsqls.SQLAGG); err != nil {
		t.Fatal(err.Error())
	}
	if sql, err := Stage1(testsqls.SQLAGG); err != nil {
		t.Fatal(err.Error())
	} else {
		if err := testsqls.SQLExec(sql); err != nil {
			t.Fatal(err.Error())
		}
	}
}

func TestStage13(t *testing.T) {
	if err := testsqls.SQLExec(testsqls.SQLWindow); err != nil {
		t.Fatal(err.Error())
	}
	if sql, err := Stage1(testsqls.SQLWindow); err != nil {
		t.Fatal(err.Error())
	} else {
		if err := testsqls.SQLExec(sql); err != nil {
			t.Fatal(err.Error())
		}
	}
}