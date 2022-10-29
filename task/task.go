package task

import (
	"github.com/qaqcatz/impomysql/connector"
	"github.com/qaqcatz/impomysql/mutation/oracle"
	"github.com/qaqcatz/impomysql/mutation/stage1"
	"github.com/qaqcatz/impomysql/mutation/stage2"
	"github.com/qaqcatz/impomysql/randgen"
	"strconv"
)

// Config: task config
type Config struct {
	Conn *connector.Connector
	RandGenConfig *randgen.Config
}

func (config *Config) ToShortString() string {
	s := config.Conn.ToString() + "\n" + config.RandGenConfig.ToString()
	return s
}

// Result: task result
type Result struct {
	RandGenRes *randgen.Results       // result of randgen.RandGen
	Stage1Res  []*stage1.InitResult   // result of stage1.InitAndExec
	Stage1ErrNum int
	Stage1ExecErrNum int
	Stage2Res  []*stage2.MutateResult // result of stage2.MutateAllAndExec, nil if stage1 error
	Stage2ErrNum int
	Stage2ExecNum int
	Stage2ExecErrNum int
	ImpoBugs   []*ImpoBug             // logical bugs, see oracle.Check
	DriverErrNum int // see *(connector.Connector).ExecSQLX
	Err        error
}

func (result *Result) ToShortString() string {
	s := ""
	if result.Err != nil {
		s = result.Err.Error()
		return s
	}
	s += "[|RandGenDDL|] " + strconv.Itoa(len(result.RandGenRes.DDLs)) + "\n"
	s += "[|RandGenDML|] " + strconv.Itoa(len(result.RandGenRes.RandSQLs)) + "\n"
	s += "[Stage1ErrNum] " + strconv.Itoa(result.Stage1ErrNum) + "\n"
	s += "[Stage1ExecErrNum] " + strconv.Itoa(result.Stage1ExecErrNum) + "\n"
	s += "[Stage2ErrNum] " + strconv.Itoa(result.Stage2ErrNum) + "\n"
	s += "[Stage2ExecNum] " + strconv.Itoa(result.Stage2ExecNum) + "\n"
	s += "[Stage2ExecErrNum] " + strconv.Itoa(result.Stage2ExecErrNum) + "\n"
	s += "[DriverErrNum] " + strconv.Itoa(result.DriverErrNum) + "\n"
	s += "[|ImpoBugs|] " + strconv.Itoa(len(result.ImpoBugs))
	return s
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

func (impoBug *ImpoBug) ToString() string {
	s := ""
	s += "[mutation name]" + impoBug.MutationName + "\n"
	s += "[is upper]" + strconv.FormatBool(impoBug.IsUpper) + "\n"
	s += "[origin sql]" + impoBug.OriginSql + "\n"
	s += "[origin result]" + impoBug.OriginResult.ToString() + "\n"
	s += "[new sql]" + impoBug.NewSql + "\n"
	s += "[new result]" + impoBug.NewResult.ToString()
	return s
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
	seed := randGenConfig.Seed

	result := &Result {
		Stage1Res: make([]*stage1.InitResult, 0),
		Stage1ErrNum: 0,
		Stage1ExecErrNum: 0,
		Stage2Res: make([]*stage2.MutateResult, 0),
		Stage2ErrNum: 0,
		Stage2ExecNum: 0,
		Stage2ExecErrNum: 0,
		ImpoBugs: make([]*ImpoBug, 0),
		DriverErrNum: 0,
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
	for i, initSql := range result.Stage1Res {
		if initSql.Err != nil || initSql.ExecResult.Err != nil {
			result.Stage2Res = append(result.Stage2Res, nil)
			continue
		}
		mutateResult := stage2.MutateAllAndExec(initSql.InitSql, seed+int64(i), conn)
		result.Stage2Res = append(result.Stage2Res, mutateResult)
	}
	// 5
	for i, initResult := range result.Stage1Res {
		if initResult.Err != nil {
			result.Stage1ErrNum += 1
			continue
		}
		if initResult.ExecResult.Err != nil {
			result.Stage1ExecErrNum += 1
			continue
		}
		originSql := initResult.InitSql
		originResult := initResult.ExecResult
		if result.Stage2Res[i].Err != nil {
			result.Stage2ErrNum += 1
			continue
		}
		for j, newResult := range result.Stage2Res[i].ExecResults {
			result.Stage2ExecNum += 1
			if newResult.Err != nil {
				result.Stage2ExecErrNum += 1
				continue
			}
			mutName := result.Stage2Res[i].MutNames[j]
			mutSql := result.Stage2Res[i].MutSqls[j]
			isUpper := result.Stage2Res[i].IsUppers[j]
			if !oracle.Check(originResult, newResult, isUpper) {
				if !oracle.DoubleCheck(conn, originSql, mutSql, originResult.Err != nil, newResult.Err != nil) {
					result.DriverErrNum += 1
					continue
				}
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