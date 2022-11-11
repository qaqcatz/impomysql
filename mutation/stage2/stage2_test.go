package stage2

import (
	"fmt"
	_ "github.com/pingcap/tidb/parser/test_driver"
	"github.com/qaqcatz/impomysql/mutation/oracle"
	"github.com/qaqcatz/impomysql/testsqls"
	"reflect"
	"testing"
)

func testImpoMutateCommon(t *testing.T, sql string, seed int64) {
	fmt.Println("==================================================")
	// prepare table COMPANY
	if err := testsqls.InitTableCOMPANY(""); err != nil {
		t.Fatalf("%+v", err)
	}
	// prepare connector
	conn, err := testsqls.GetConnector("")
	if err != nil {
		t.Fatalf("%+v", err)
	}

	// execute original sql
	t.Log("[origin]", sql)
	originResult := conn.ExecSQL(sql)
	if originResult.Err != nil {
		t.Fatalf("%+v", originResult.Err)
	}
	t.Log("[origin result]", originResult.ToString())

	// calculate candidates
	v, err := CalCandidates(sql)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	root := v.Root

	i := 0
	for mutationName, mutationPoints := range v.Candidates {
		t.Log(i, "====================")
		i += 1

		t.Log("[MutationName]", mutationName)
		// mutate all points
		j := 0
		for _, point := range mutationPoints {
			t.Log(i, ".", j, "==========")
			j += 1

			t.Log("[type]", reflect.TypeOf(point.Node))
			t.Log("[candidate]", point.Node)
			t.Log("[flag]", point.Flag)

			// do!
			newSql, err := ImpoMutate(root, point, seed)
			if err != nil {
				t.Fatalf("%+v", err)
			}

			// execute new sql
			t.Log("[newSql]", newSql)
			result := conn.ExecSQL(newSql)
			if result.Err != nil {
				t.Fatalf("%+v", result.Err)
			}
			t.Log("[new result]", result.ToString())

			// check impo result
			if !oracle.Check(originResult, result, ((point.U ^ point.Flag) ^ 1) == 1) {
				t.Fatal("!IMPO")
			}
		}
	}
}

func TestImpoMutateSelectValue(t *testing.T) {
	testImpoMutateCommon(t, testsqls.SQLSelectValue, 10001)
	testImpoMutateCommon(t, testsqls.SQLSelectValue2, 10002)
	testImpoMutateCommon(t, testsqls.SQLSelectValue3, 10003)
}
func TestImpoMutateSubQuery(t *testing.T) {
	testImpoMutateCommon(t, testsqls.SQLSubQuery, 20001)
	testImpoMutateCommon(t, testsqls.SQLSubQuery2, 20002)
	testImpoMutateCommon(t, testsqls.SQLSubQuery3, 20003)
	testImpoMutateCommon(t, testsqls.SQLSubQuery4, 20004)
	testImpoMutateCommon(t, testsqls.SQLSubQuery5, 20005)
}

func TestImpoMutateJOIN(t *testing.T) {
	testImpoMutateCommon(t, testsqls.SQLJOIN, 30001)
	testImpoMutateCommon(t, testsqls.SQLJOIN2, 30002)
	testImpoMutateCommon(t, testsqls.SQLJOIN3, 30003)
	testImpoMutateCommon(t, testsqls.SQLJOIN4, 30004)
	testImpoMutateCommon(t, testsqls.SQLJOIN5, 30005)
	testImpoMutateCommon(t, testsqls.SQLJOIN6, 30006)
}

func TestImpoMutateUNION(t *testing.T) {
	testImpoMutateCommon(t, testsqls.SQLUNION, 40001)
	testImpoMutateCommon(t, testsqls.SQLUNION2, 40002)
}

func TestImpoMutateWITH(t *testing.T) {
	testImpoMutateCommon(t, testsqls.SQLWITH, 50001)
	testImpoMutateCommon(t, testsqls.SQLWITH2, 50002)
}

func TestImpoMutateIN(t *testing.T) {
	testImpoMutateCommon(t, testsqls.SQLIN, 60001)
	testImpoMutateCommon(t, testsqls.SQLIN2, 60002)
}

func TestImpoMutateWhere(t *testing.T) {
	testImpoMutateCommon(t, testsqls.SQLWHERE, 70001)
	testImpoMutateCommon(t, testsqls.SQLWHERE, 70003)
}

func TestImpoMutateHaving(t *testing.T) {
	testImpoMutateCommon(t, testsqls.SQLHAVING, 80001)
	testImpoMutateCommon(t, testsqls.SQLHAVING, 80002)
}

func TestImpoMutateLIKE(t *testing.T) {
	testImpoMutateCommon(t, testsqls.SQLLIKE, 90001)
	testImpoMutateCommon(t, testsqls.SQLLIKE, 90003)
}

func TestImpoMutateRegExp(t *testing.T) {
	testImpoMutateCommon(t, testsqls.SQLRegExp, 100001)
	testImpoMutateCommon(t, testsqls.SQLRegExp, 100003)
}

func TestImpoMutateBetween(t *testing.T) {
	testImpoMutateCommon(t, testsqls.SQLBetween, 110001)
	testImpoMutateCommon(t, testsqls.SQLBetween2, 110002)
}