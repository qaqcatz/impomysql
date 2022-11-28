// Package oracle: check results to see if there is a logical bug according to implication oracle
package oracle

import (
	"github.com/qaqcatz/impomysql/connector"
)

// Check: check results to see if there is a logical bug according to implication oracle.
// return false if there is a logical bug, otherwise return true.
//
// Note that implication oracle cannot support error oracle.
// You cannot have any errors in your results, otherwise we will return an error
func Check(originResult *connector.Result, mutatedResult *connector.Result, isUpper bool) (bool, error) {
	cmp, err := originResult.CMP(mutatedResult)
	if err != nil {
		return false, err
	}
	if cmp == 0 {
		return true, nil
	}
	if (isUpper && cmp == -1) || (!isUpper && cmp == 1) {
		return true, nil
	}
	return false, nil
}