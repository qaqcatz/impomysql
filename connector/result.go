package connector

import (
	"github.com/pkg/errors"
	"github.com/go-sql-driver/mysql"
	"reflect"
	"strconv"
	"time"
)

// Result:
//
// query result, for example:
//   +-----+------+------+
//   | 1+2 | ID   | NAME | -> ColumnNames: 1+2,    ID,  NAME
//   +-----+------+------+ -> ColumnTypes: BIGINT, INT, TEXT
//   |   3 |    1 | H    | -> Rows[0]:     3,      1,   H
//   |   3 |    2 | Z    | -> Rows[1]:     3,      2,   Z
//   |   3 |    3 | Y    | -> Rows[2]:     3,      3,   Y
//   +-----+------+------+
// or error, for example:
//  Err: ERROR 1054 (42S22): Unknown column 'T' in 'field list'
//
// note that:
//
// len(ColumnNames) = len(ColumnTypes) = len(Rows[i]);
//
// if the statement is not SELECT, then the ColumnNames, ColumnTypes and Rows are empty
type Result struct {
	ColumnNames []string
	ColumnTypes []string
	Rows [][]string
	Err error
	Time time.Duration // total time
}

func (result *Result) ToString() string {
	str := ""
	str += "ColumnName(ColumnType)s: "
	for i, columnName := range result.ColumnNames {
		str += " " + columnName + "(" + result.ColumnTypes[i] + ")"
	}
	str += "\n"
	for i, row := range result.Rows {
		str += "row " + strconv.Itoa(i) + ":"
		for _, data := range row {
			str += " " + data
		}
		str += "\n"
	}
	if result.Err != nil {
		str += "Error: " + result.Err.Error() + "\n"
	}

	str += result.Time.String()
	return str
}

// Result.FlatRows: [["1","2"],["3","4"]] -> ["1,2", "3,4"]
func (result *Result) FlatRows() []string {
	flt := make([]string, 0)
	for _, r := range result.Rows {
		t := ""
		for i, e := range r {
			if i != 0 {
				t += ","
			}
			t += e
		}
		flt = append(flt, t)
	}
	return flt
}

// Result.IsEmpty: if the result is empty
func (result *Result) IsEmpty() bool {
	return len(result.ColumnNames) == 0
}

func (result *Result) GetErrorCode() (int, error) {
	if result.Err == nil {
		return -1, errors.New("[Result.GetErrorCode]result.Err == nil")
	}
	rootCause := errors.Cause(result.Err)
	if driverErr, ok := rootCause.(*mysql.MySQLError); ok { // Now the error number is accessible directly
		return int(driverErr.Number), nil
	} else {
		return -1, errors.New("[Result.GetErrorCode]not *mysql.MySQLError " + reflect.TypeOf(rootCause).String())
	}
}

// Result.CMP:
//   -1: another contains this
//   0: eq
//   1: this contains another
//   2: others
//   error: this.Err or another.Err
//   do not consider the column name
func (this *Result) CMP(another *Result) (int, error) {
	if this.Err != nil {
		return -2, errors.New("[Result.CMP]this error")
	}
	if another.Err != nil {
		return -2, errors.New("[Result.CMP]another error")
	}

	empty1 := this.IsEmpty()
	empty2 := another.IsEmpty()
	if empty1 || empty2 {
		// empty1&&!empty2, !empty1&&empty2, empty1&&empty2
		if (empty1 && empty2) {
			return 0, nil
		}
		if empty1 {
			// empty1&&!empty2
			return -1, nil;
		} else {
			// !empty1&&empty2
			return 1, nil;
		}
	}

	if len(this.ColumnNames) != len(another.ColumnNames) {
		return 2, nil
	}

	res1 := this.FlatRows()
	res2 := another.FlatRows()

	mp := make(map[string]int)
	for i := 0; i < len(res2); i++ {
		if num, ok := mp[res2[i]]; ok {
			mp[res2[i]] = num + 1
		} else {
			mp[res2[i]] = 1
		}
	}
	allInAnother := true
	for i := 0; i < len(res1); i++ {
		if num, ok := mp[res1[i]]; ok {
			if num <= 1 {
				delete(mp, res1[i])
			} else {
				mp[res1[i]] = num - 1
			}
		} else {
			allInAnother = false
		}
	}

	if allInAnother {
		if len(mp) == 0 {
			return 0, nil
		} else {
			return -1, nil
		}
	} else {
		if len(mp) == 0 {
			return 1, nil
		} else {
			return 2, nil
		}
	}
}