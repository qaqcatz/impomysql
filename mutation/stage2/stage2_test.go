package stage2

import (
	"fmt"
	"github.com/pingcap/tidb/parser"
	_ "github.com/pingcap/tidb/parser/test_driver"
	"github.com/qaqcatz/impomysql/learnast"
	"github.com/qaqcatz/impomysql/testsqls"
	"reflect"
	"sort"
	"strconv"
	"testing"
)

func testCalCandidatesCommon(t *testing.T, sql string, oracle []string) {
	fmt.Println("==================================================")
	if err := testsqls.EnsureDBTEST(); err != nil {
		t.Fatal(err.Error())
	}
	if err := testsqls.InitTableCOMPANY(); err != nil {
		t.Fatal(err.Error())
	}

	p := parser.New()
	stmtNodes, _, err := p.Parse(sql, "", "")
	if err != nil {
		t.Fatal(err.Error())
	}
	if stmtNodes == nil || len(stmtNodes) == 0 {
		t.Fatal("stmtNodes == nil || len(stmtNodes) == 0 ")
	}
	rootNode := &stmtNodes[0]

	v := CalCandidates(*rootNode)

	fmt.Println("[origin]", sql)
	i := 0
	ck := make([]string, 0)
	for can, flag := range v.Candidates {
		fmt.Println(i, "==========")
		fmt.Println("[type]", reflect.TypeOf(can))
		fmt.Print("[candidate]")
		learnast.PrintNode(can)
		fmt.Println()
		fmt.Println("[flag]", flag)
		i += 1
		ck = append(ck, reflect.TypeOf(can).String()+"#"+strconv.Itoa(flag))
	}

	if !check(ck, oracle) {
		t.Fatal("!check:\n[ck]", ck, "\n[oracle]", oracle)
	}
}

// type#flag
func check(ck []string, oracle []string) bool {
	if len(ck) != len(oracle) {
		return false
	}
	sort.Strings(ck)
	sort.Strings(oracle)
	for i, _ := range ck {
		if ck[i] != oracle[i] {
			return false;
		}
	}
	return true
}

func TestCalCandidatesAGG(t *testing.T) {
	testCalCandidatesCommon(t, testsqls.SQLAGG,
		[]string{"*ast.BinaryOperationExpr#1",
		"*ast.BinaryOperationExpr#1",
		"*ast.SelectStmt#1",
		"*ast.SelectStmt#1"});
}

func TestCalCandidatesWindow(t *testing.T) {
	testCalCandidatesCommon(t, testsqls.SQLWindow,
		[]string{"*ast.SelectStmt#1"});
}

func TestCalCandidatesSelectValue(t *testing.T) {
	testCalCandidatesCommon(t, testsqls.SQLSelectValue,
		[]string{"*ast.SelectStmt#1"});
	testCalCandidatesCommon(t, testsqls.SQLSelectValue2,
		[]string{"*ast.SelectStmt#1"});
	testCalCandidatesCommon(t, testsqls.SQLSelectValue3,
		[]string{"*ast.SelectStmt#1"});
}
func TestCalCandidatesSQLSubQuery(t *testing.T) {
	testCalCandidatesCommon(t, testsqls.SQLSubQuery,
		[]string{"*ast.SelectStmt#1",
			"*ast.BinaryOperationExpr#1"});
	testCalCandidatesCommon(t, testsqls.SQLSubQuery2,
		[]string{"*ast.SelectStmt#1",
			"*ast.BinaryOperationExpr#1",
			"*ast.SelectStmt#1",
			"*ast.CompareSubqueryExpr#1"});
	testCalCandidatesCommon(t, testsqls.SQLSubQuery3,
		[]string{"*ast.PatternInExpr#0",
			"*ast.SelectStmt#1",
			"*ast.SelectStmt#0"});
	testCalCandidatesCommon(t, testsqls.SQLSubQuery4,
		[]string{"*ast.SelectStmt#1",
			"*ast.CompareSubqueryExpr#1",
			"*ast.SelectStmt#0",
			"*ast.BinaryOperationExpr#0"});
	testCalCandidatesCommon(t, testsqls.SQLSubQuery5,
		[]string{"*ast.SelectStmt#0",
			"*ast.BinaryOperationExpr#0",
			"*ast.SelectStmt#1"});
}

func TestCalCandidatesJOIN(t *testing.T) {
	testCalCandidatesCommon(t, testsqls.SQLJOIN,
		[]string{"*ast.BinaryOperationExpr#1",
			"*ast.BinaryOperationExpr#1",
			"*ast.SelectStmt#1",
			"*ast.SelectStmt#1"});
	testCalCandidatesCommon(t, testsqls.SQLJOIN2,
		[]string{"*ast.SelectStmt#1"});
	testCalCandidatesCommon(t, testsqls.SQLJOIN3,
		[]string{"*ast.SelectStmt#1"});
	testCalCandidatesCommon(t, testsqls.SQLJOIN4,
		[]string{"*ast.SelectStmt#1",
			"*ast.SelectStmt#1",
			"*ast.BinaryOperationExpr#1",
			"*ast.BinaryOperationExpr#1"});
	testCalCandidatesCommon(t, testsqls.SQLJOIN5,
		[]string{"*ast.SelectStmt#1",
			"*ast.SelectStmt#1",
			"*ast.BinaryOperationExpr#1"});
	testCalCandidatesCommon(t, testsqls.SQLJOIN6,
		[]string{"*ast.SelectStmt#1"});
}

func TestCalCandidatesLIMIT(t *testing.T) {
	testCalCandidatesCommon(t, testsqls.SQLLIMIT,
		[]string{"*ast.SelectStmt#1"});
	testCalCandidatesCommon(t, testsqls.SQLLIMIT2,
		[]string{"*ast.SelectStmt#1"});
}

func TestCalCandidatesUNION(t *testing.T) {
	testCalCandidatesCommon(t, testsqls.SQLUNION,
		[]string{"*ast.SelectStmt#1",
			"*ast.SelectStmt#1",
			"*ast.SelectStmt#1",
			"*ast.SelectStmt#1"});
	testCalCandidatesCommon(t, testsqls.SQLUNION2,
		[]string{"*ast.SelectStmt#1",
			"*ast.SelectStmt#1",
			"*ast.SelectStmt#1",
			"*ast.SelectStmt#1",
			"*ast.SelectStmt#1",
			"*ast.SelectStmt#1",
			"*ast.SelectStmt#1"});
}

func TestCalCandidatesWITH(t *testing.T) {
	testCalCandidatesCommon(t, testsqls.SQLWITH,
		[]string{"*ast.SelectStmt#1",
			"*ast.SelectStmt#1"});
	testCalCandidatesCommon(t, testsqls.SQLWITH2,
		[]string{"*ast.BinaryOperationExpr#1",
			"*ast.SelectStmt#1",
			"*ast.SelectStmt#1",
			"*ast.SelectStmt#1"});
}

func TestCalCandidatesIN(t *testing.T) {
	testCalCandidatesCommon(t, testsqls.SQLIN,
		[]string{"*ast.SelectStmt#1"});
}