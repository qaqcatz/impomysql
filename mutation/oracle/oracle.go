// Package oracle: check results to see if there is a logical bug according to implication oracle
package oracle

import "github.com/qaqcatz/impomysql/connector"

// Check: check results to see if there is a logical bug according to implication oracle.
// return false if there is a logical bug, otherwise return true.
func Check(originResult *connector.Result, mutatedResult *connector.Result, isUpper bool) bool {
	// ignore error
	isErr1 := (originResult.Err != nil)
	isErr2 := (mutatedResult.Err != nil)
	if isErr1 || isErr2 {
		return true
	}

	empty1 := originResult.IsEmpty()
	empty2 := mutatedResult.IsEmpty()
	if empty1 || empty2 {
		// empty1&&!empty2, !empty1&&empty2, empty1&&empty2
		if (empty1 && empty2) {
			return true
		}
		// origin < new
		if (empty1) {
			// empty1&&!empty2
			return isUpper;
		} else {
			// !empty1&&empty2
			return !isUpper;
		}
	}

	if len(originResult.ColumnNames) != len(mutatedResult.ColumnNames) {
		return false
	}
	// Due to the difference between the restored sql and the original sql,
	// we can not compare compare column names and types. (consider value select)
	//for i, _ := range originResult.ColumnNames {
	//	if originResult.ColumnNames[i] != mutatedResult.ColumnNames[i] {
	//		return false
	//	}
	//	if originResult.ColumnTypes[i] != mutatedResult.ColumnTypes[i] {
	//		return false
	//	}
	//}

	// Rows -> []string
	res1 := originResult.FlatRows()
	res2 := mutatedResult.FlatRows()

	if !isUpper {
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