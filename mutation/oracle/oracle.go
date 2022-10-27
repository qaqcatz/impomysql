// Package oracle: check impo
package oracle

import "github.com/qaqcatz/impomysql/connector"

// Check: check impo
func Check(originResult *connector.Result, newResult *connector.Result, isUpper bool) bool {
	isErr1 := (originResult.Err != nil)
	isErr2 := (newResult.Err != nil)
	if isErr1 || isErr2 {
		// isErr1&&!isErr2, !isErr1&&isErr2, isErr1&&isErr2
		if (isErr1 && isErr2) {
			return true
		}
		return false
	}

	empty1 := originResult.IsEmpty()
	empty2 := newResult.IsEmpty()
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

	// Rows -> string
	res1 := originResult.FlatRows()
	res2 := newResult.FlatRows()

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