package stage2

import (
	"fmt"
	"github.com/pingcap/tidb/parser"
	_ "github.com/pingcap/tidb/parser/test_driver"
	"github.com/qaqcatz/impomysql/connector"
	"github.com/qaqcatz/impomysql/testsqls"
	"reflect"
	"testing"
)

func check(originResult *connector.Result, newResult *connector.Result, uflag int) bool {
	empty1 := len(originResult.ColumnNames) == 0
	empty2 := len(newResult.ColumnNames) == 0
	if empty1 || empty2 {
		// empty1&&!empty2, !empty1&&empty2, empty1&&empty2
		if (empty1 && empty2) {
			return true
		}
		// origin < new
		if (empty1) {
			// empty1&&!empty2
			return uflag == 1;
		} else {
			// !empty1&&empty2
			return uflag == 0;
		}
	}
	if len(originResult.ColumnNames) != len(newResult.ColumnNames) {
		return false
	}
	// Due to the difference between the restored sql and the original sql,
	// we can not compare compare column names and types. (consider value select)
	//for i, _ := range originResult.ColumnNames {
	//	if originResult.ColumnNames[i] != newResult.ColumnNames[i] {
	//		return false
	//	}
	//	if originResult.ColumnTypes[i] != newResult.ColumnTypes[i] {
	//		return false
	//	}
	//}
	res1 := make([]string, 0)
	for _, r := range originResult.Rows {
		t := ""
		for _, e := range r {
			t += e
			t += ","
		}
		res1 = append(res1, t)
	}
	res2 := make([]string, 0)
	for _, r := range newResult.Rows {
		t := ""
		for _, e := range r {
			t += e
			t += ","
		}
		res2 = append(res2, t)
	}
	if uflag == 0 {
		// negative
		t := res1
		res1 = res2
		res2 = t
	}
	// res1 < res2
	mp := make(map[string]int)
	for i := 0; i < len(res2); i++ {
		if num, ok := mp[res2[i]]; ok {
			mp[res2[i]] = num + 1
		} else {
			mp[res2[i]] = 1
		}
	}
	for i := 0; i < len(res1); i++ {
		if num, ok := mp[res1[i]]; ok {
			if num <= 1 {
				delete(mp, res1[i])
			} else {
				mp[res1[i]] = num - 1
			}
		} else {
			return false
		}
	}
	return true
}

func testImpoMutateCommon(t *testing.T, sql string, seed int64) {
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

	t.Log("[origin]", sql)

	conn, err := testsqls.GetConnector()
	if err != nil {
		t.Fatal(err.Error())
	}

	originResult := conn.ExecSQL(sql)
	if originResult.Err != nil {
		t.Fatal(originResult.Err.Error())
	}

	t.Log("[origin result]", originResult.ToString())

	i := 0
	for k, v := range v.Candidates {
		t.Log(i, "====================")
		i += 1
		t.Log("[MutationName]", k)
		j := 0
		for _, can := range v {
			t.Log(i, ".", j, "==========")
			j += 1
			t.Log("[type]", reflect.TypeOf(can.Node))
			t.Log("[candidate]", can.Node)
			t.Log("[flag]", can.Flag)

			newSql, err := ImpoMutate(*rootNode, can, seed)
			if err != nil {
				t.Fatal(err.Error())
			}
			t.Log("[newSql]", string(newSql))

			result := conn.ExecSQL(string(newSql))
			if result.Err != nil {
				t.Fatal(result.Err.Error())
			}

			t.Log("[new result]", result.ToString())

			if !check(originResult, result, (can.U ^ can.Flag) ^ 1) {
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
	//testImpoMutateCommon(t, testsqls.SQLBetween3, 110003)
}