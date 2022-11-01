package randgen

import (
	"encoding/json"
	"errors"
	"github.com/qaqcatz/impomysql/connector"
	"github.com/qaqcatz/nanoshlib"
	"path"
	"runtime"
	"strconv"
)

const gorandgensh = "gorandgensh"

// Config: see https://github.com/qaqcatz/gorandgensh
type Config struct {
	GoRandGenPath string
	ZZFilePath string
	YYFilePath string
	QueriesNum int
	Seed int64
}

func (config *Config) ToString() string {
	s := "[randgen config]\n"
	s += "==================================================\n"
	s += "[GoRandGenPath] " + config.GoRandGenPath + "\n"
	s += "[ZZFilePath] " + config.ZZFilePath + "\n"
	s += "[YYFilePath] " + config.YYFilePath + "\n"
	s += "[QueriesNum] " + strconv.Itoa(config.QueriesNum) + "\n"
	s += "[Seed] " + strconv.FormatInt(config.Seed, 10) + "\n"
	s += "=================================================="
	return s
}

// Results: see https://github.com/qaqcatz/gorandgensh
type Results struct {
	ZZFilePath string `json:"zzFilePath"`
	YYFilePath string `json:"yyFilePath"`
	Seed int64 `json:"seed"`
	DDLs []string `json:"ddls"`
	RandSQLs []string `json:"randsqls"`
	Err error `json:"-"`
}

// RandGen: see Config and Results
func RandGen(config *Config) *Results {
	sqls := &Results{
		Err: nil,
	}

	goRandGenPath := config.GoRandGenPath
	if goRandGenPath == "" {
		// default path
		temp, err := getPackagePath()
		if err != nil {
			sqls.Err = errors.New("RandGen: getPackagePath() error: " + err.Error())
			return sqls
		}
		goRandGenPath = path.Join(temp, gorandgensh)
	}
	zzFilePath := config.ZZFilePath
	yyFilePath := config.YYFilePath
	queriesNum := config.QueriesNum
	seed := config.Seed
	// see https://github.com/qaqcatz/nanoshlib
	outStream, errStream, err := nanoshlib.Exec(goRandGenPath+" "+
		zzFilePath+" "+yyFilePath+" "+strconv.Itoa(queriesNum)+" "+strconv.FormatInt(seed, 10), 0)
	if err != nil {
		sqls.Err = errors.New("RandGen: gen test error: "+err.Error()+": "+string(errStream))
		return sqls
	}

	err = json.Unmarshal(outStream, sqls)
	if err != nil {
		sqls.Err = errors.New("RandGen: json.Unmarshal() error: "+err.Error())
		return sqls
	}
	return sqls
}

// RandGenAndExecDDL: RandGen + exec ddl
func RandGenAndExecDDL(config *Config, conn *connector.Connector) *Results {
	result := RandGen(config)
	for _, ddl := range result.DDLs {
		execRes := conn.ExecSQL(ddl)
		if execRes.Err != nil {
			result.Err = errors.New("RandGenAndExecDDL: " + execRes.Err.Error())
			return result
		}
	}
	return result
}

// getPackagePath: get the package actual path, then you can read files under the path.
func getPackagePath() (string, error) {
	if _, file, _, ok := runtime.Caller(0); !ok {
		return "", errors.New("PackagePath: runtime.Caller(0) error ")
	} else {
		return path.Join(file, "../"), nil
	}
}