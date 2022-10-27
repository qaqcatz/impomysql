package task

import (
	"github.com/qaqcatz/impomysql/connector"
	"github.com/qaqcatz/impomysql/mutation/oracle"
	"github.com/qaqcatz/impomysql/mutation/stage1"
	"github.com/qaqcatz/impomysql/mutation/stage2"
	"github.com/qaqcatz/impomysql/randgen"
)

// Config: task config
type Config struct {
	Conn *connector.Connector
	RandGenConfig *randgen.Config
	MutSeed int64 // for stage2.MutateAll
}

// Result: task result
type Result struct {
	RandGenRes *randgen.Results       // result of randgen.RandGen
	Stage1Res  []*stage1.InitResult   // result of stage1.InitAndExec
	Stage2Res  []*stage2.MutateResult // result of stage2.MutateAllAndExec, nil if stage1 error
	ImpoBugs   []*ImpoBug             // logical bugs, see oracle.Check
	Err        error
}

// ImpoBug: logical bug
type ImpoBug struct {
	//DDLs []string
	OriginSql string
	OriginResult *connector.Result
	NewSql string
	NewResult *connector.Result
	MutationName string
	IsUpper bool // true: theoretically, OriginResult < NewResult
}

// Run:  Basic task for finding logical bug:
//
// 1. random generate sql statements -- randgen.RandGen
//
// 2. initialize random sqls, filter parse error. -- stage1.Init
//
// 3. execute random sqls, filter execute error. -- connector.Connector
//
// 4. for each random sqls, try all of its mutation points, get mutated sqls. -- stage2.MutateAll
//
// 5. execute each mutated sqls, compare their results with the original result,
// detect logical bugs. -- connector.Connector, oracle.Check
//
// see Config and Result for the input/output message
//
// Note that: you should init database yourself!
func Run(config *Config) *Result {
	conn := config.Conn
	randGenConfig := config.RandGenConfig
	mutSeed := config.MutSeed

	result := &Result {
		Stage1Res: make([]*stage1.InitResult, 0),
		Stage2Res: make([]*stage2.MutateResult, 0),
		ImpoBugs: make([]*ImpoBug, 0),
		Err: nil,
	}

	// 1
	result.RandGenRes = randgen.RandGenAndExecDDL(randGenConfig, conn)
	if result.RandGenRes.Err != nil {
		return result
	}
	// 2, 3
	for _, sql := range result.RandGenRes.RandSQLs {
		initResult := stage1.InitAndExec(sql, conn)
		result.Stage1Res = append(result.Stage1Res, initResult)
	}
	// 4
	for _, initSql := range result.Stage1Res {
		if initSql.Err != nil || initSql.ExecResult.Err != nil {
			result.Stage2Res = append(result.Stage2Res, nil)
			continue
		}
		mutateResult := stage2.MutateAllAndExec(initSql.InitSql, mutSeed, conn)
		result.Stage2Res = append(result.Stage2Res, mutateResult)
	}
	// 5
	for i, initResult := range result.Stage1Res {
		if initResult.Err != nil || initResult.ExecResult.Err != nil {
			continue
		}
		originSql := initResult.InitSql
		originResult := initResult.ExecResult
		if result.Stage2Res[i].Err != nil {
			continue
		}
		for j, newResult := range result.Stage2Res[i].ExecResults {
			if newResult.Err != nil {
				continue
			}
			mutName := result.Stage2Res[i].MutNames[j]
			mutSql := result.Stage2Res[i].MutSqls[j]
			isUpper := result.Stage2Res[i].IsUppers[j]
			if !oracle.Check(originResult, newResult, isUpper) {
				result.ImpoBugs = append(result.ImpoBugs, &ImpoBug {
					//DDLs: result.RandGenRes.DDLs,
					OriginSql: originSql,
					OriginResult: originResult,
					NewSql: mutSql,
					NewResult: newResult,
					MutationName: mutName,
					IsUpper: isUpper,
				})
			}
		}
	}
	return result
}