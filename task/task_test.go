package task

import (
	"github.com/qaqcatz/impomysql/randgen"
	"github.com/qaqcatz/impomysql/testsqls"
	"testing"
)

func TestRun(t *testing.T) {
	err := testsqls.InitDBTEST()
	if err != nil {
		t.Fatal(err.Error())
	}

	conn, err := testsqls.GetConnector()
	if err != nil {
		t.Fatal(err.Error())
	}
	randGenConfig := &randgen.Config {
		ZZFilePath: testsqls.GetTestZZPath(),
		YYFilePath: testsqls.GetTestYYPath(),
		QueriesNum: 100,
		Seed:       123456,
	}
	config := &Config{
		Conn: conn,
		RandGenConfig: randGenConfig,
		MutSeed: 123456,
	}
	result := Run(config)
	if result.RandGenRes.Err != nil {
		t.Fatal(result.RandGenRes.Err)
	}
	for i, randsql := range result.RandGenRes.RandSQLs {
		t.Log(i, "========================================")
		t.Log("[random sql]", randsql)
		t.Log("[init sql]", result.Stage1Res[i].InitSql)
		if result.Stage1Res[i].Err != nil {
			t.Log("stage1 error:", result.Stage1Res[i].Err)
			continue
		}
		if result.Stage1Res[i].ExecResult.Err != nil {
			t.Log("stage1 error:", result.Stage1Res[i].ExecResult.Err)
			continue
		}
		if result.Stage2Res[i].Err != nil {
			t.Log("stage2 error:", result.Stage2Res[i].Err)
			continue
		}
		t.Log("|stage2|", len(result.Stage2Res[i].ExecResults))
		for j, _ := range result.Stage2Res[i].ExecResults {
			if result.Stage2Res[i].MutErrs[j] != nil {
				t.Log("stage2 error", j, result.Stage2Res[i].MutErrs[j])
			}
		}
	}
	t.Log("|bug|", len(result.ImpoBugs))
	for i, bug := range result.ImpoBugs {
		t.Log(i, "!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		t.Log("[origin sql]", bug.OriginSql)
		t.Log("[origin result]", bug.OriginResult.ToString())
		t.Log("[new sql]", bug.NewSql)
		t.Log("[new result]", bug.NewResult.ToString())
		t.Log("[mutation name]", bug.MutationName)
		t.Log("[is upper]", bug.IsUpper)
	}
}
