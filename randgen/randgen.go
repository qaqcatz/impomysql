package randgen

import (
	"encoding/json"
	"errors"
	"path"
	"runtime"
	"github.com/qaqcatz/nanoshlib"
	"strconv"
	"strings"
)

// getPackagePath: get the package actual path, then you can read files under the path.
func getPackagePath() (string, error) {
	if _, file, _, ok := runtime.Caller(0); !ok {
		return "", errors.New("PackagePath: runtime.Caller(0) error ")
	} else {
		return path.Join(file, "../"), nil
	}
}

// SQLS: see https://github.com/qaqcatz/gorandgensh
type SQLS struct {
	ZZFilePath string `json:"zzFilePath"`
	YYFilePath string `json:"yyFilePath"`
	Seed int64 `json:"seed"`
	DDLs []string `json:"ddls"`
	RandSQLs []string `json:"randsqls"`
}

// default zz, yy files under the package
const (
	gorandgensh = "gorandgensh"
	magic = "!@#$%^&*()"
	ZZDefault = magic + "default.zz.lua"
	YYDefault = magic + "default.yy"
)

// RandGen: see https://github.com/qaqcatz/gorandgensh
func RandGen(zzFilePath string, yyFilePath string, queriesNum int, seed int64) (*SQLS, error) {
	packagePath, err := getPackagePath()
	if err != nil {
		return nil, errors.New("RandGen: getPackagePath() error ")
	}
	if strings.HasPrefix(zzFilePath, magic) {
		zzFilePath = path.Join(packagePath, zzFilePath[len(magic):])
	}
	if strings.HasPrefix(yyFilePath, magic) {
		yyFilePath = path.Join(packagePath, yyFilePath[len(magic):])
	}
	// see https://github.com/qaqcatz/nanoshlib
	outStream, errStream, err := nanoshlib.Exec(packagePath+"/"+gorandgensh+" "+
		zzFilePath+" "+yyFilePath+" "+strconv.Itoa(queriesNum)+" "+strconv.FormatInt(seed, 10), 0)
	if err != nil {
		return nil, errors.New("RandGen: gen test error: "+err.Error()+": "+string(errStream))
	}
	sqls := &SQLS{}
	err = json.Unmarshal(outStream, sqls)
	if err != nil {
		return nil, errors.New("RandGen: json.Unmarshal() error: "+err.Error())
	}
	return sqls, nil
}