package task

import (
	"github.com/qaqcatz/impomysql/connector"
	"github.com/qaqcatz/impomysql/mutation/stage2"
)

var ignoreErrorCode = map[int]string {
	1690 : "value out of range",
}

var ignoreMutationName = map[string]int {
	stage2.FixMHaving0L: 0,
	stage2.FixMHaving1U: 1,
	stage2.FixMOn0L: 0,
	stage2.FixMOn1U: 1,
	stage2.FixMWhere0L: 0,
	stage2.FixMWhere1U: 1,
}

func IgnoreError(mutationName string, result *connector.Result) bool {
	if _, ok := ignoreMutationName[mutationName]; ok {
		return true
	}
	errCode, err := result.GetErrorCode()
	if err == nil {
		if _, ok := ignoreErrorCode[errCode]; ok {
			return true
		}
	}
	return false
}