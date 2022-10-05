package mutation

import (
	"github.com/qaqcatz/IMPOMySQL/IMPOMySQL/testsqls"
	"testing"
)

func TestLearnAST(t *testing.T) {
	if err := testsqls.InitDBTEST(); err != nil {
		t.Fatal(err.Error())
	}
	if err := testsqls.InitTableCOMPANY(); err != nil {
		t.Fatal(err.Error())
	}
}

func TestLearnAST2(t *testing.T) {
	if err := testsqls.SQLExec(testsqls.SQLAGG); err != nil {
		t.Fatal(err.Error())
	}
	if sql, err := learnAST(testsqls.SQLAGG); err != nil {
		t.Fatal(err.Error())
	} else {
		t.Log(sql)
	}
}

func TestLearnAST3(t *testing.T) {
	if err := testsqls.SQLExec(testsqls.SQLWindow); err != nil {
		t.Fatal(err.Error())
	}
	if sql, err := learnAST(testsqls.SQLWindow); err != nil {
		t.Fatal(err.Error())
	} else {
		t.Log(sql)
	}
}

func TestLearnAST4(t *testing.T) {
	if err := testsqls.SQLExec(testsqls.SQLSelectValue); err != nil {
		t.Fatal(err.Error())
	}
	if sql, err := learnAST(testsqls.SQLSelectValue); err != nil {
		t.Fatal(err.Error())
	} else {
		t.Log(sql)
	}
}
func TestLearnAST5(t *testing.T) {
	if err := testsqls.SQLExec(testsqls.SQLSelectValue2); err != nil {
		t.Fatal(err.Error())
	}
	if sql, err := learnAST(testsqls.SQLSelectValue2); err != nil {
		t.Fatal(err.Error())
	} else {
		t.Log(sql)
	}
}

func TestLearnAST6(t *testing.T) {
	if err := testsqls.SQLExec(testsqls.SQLSelectValue3); err != nil {
		t.Fatal(err.Error())
	}
	if sql, err := learnAST(testsqls.SQLSelectValue3); err != nil {
		t.Fatal(err.Error())
	} else {
		t.Log(sql)
	}
}

func TestLearnAST7(t *testing.T) {
	if err := testsqls.SQLExec(testsqls.SQLSubQuery); err != nil {
		t.Fatal(err.Error())
	}
	if sql, err := learnAST(testsqls.SQLSubQuery); err != nil {
		t.Fatal(err.Error())
	} else {
		t.Log(sql)
	}
}

func TestLearnAST8(t *testing.T) {
	if err := testsqls.SQLExec(testsqls.SQLSubQuery2); err != nil {
		t.Fatal(err.Error())
	}
	if sql, err := learnAST(testsqls.SQLSubQuery2); err != nil {
		t.Fatal(err.Error())
	} else {
		t.Log(sql)
	}
}

func TestLearnAST9(t *testing.T) {
	if err := testsqls.SQLExec(testsqls.SQLSubQuery3); err != nil {
		t.Fatal(err.Error())
	}
	if sql, err := learnAST(testsqls.SQLSubQuery3); err != nil {
		t.Fatal(err.Error())
	} else {
		t.Log(sql)
	}
}

func TestLearnAST10(t *testing.T) {
	if err := testsqls.SQLExec(testsqls.SQLSubQuery4); err != nil {
		t.Fatal(err.Error())
	}
	if sql, err := learnAST(testsqls.SQLSubQuery4); err != nil {
		t.Fatal(err.Error())
	} else {
		t.Log(sql)
	}
}