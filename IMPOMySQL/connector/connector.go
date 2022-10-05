package connector

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strconv"
	"time"
)

// Connector: connect to MySQL, execute raw sql statements, return raw execution result or error.
type Connector struct {
	db *gorm.DB
}

// NewConnector: create Connector.
func NewConnector(host string, port int, username string, password string, dbname string) (*Connector, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		username, password, host, port, dbname)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, errors.New("NewConnector: create Connector error: " + err.Error())
	}
	return &Connector{
		db: db,
	}, nil
}

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

// ExecSQL: execute sql, return *Result.
func (conn *Connector) ExecSQL(sql string) *Result {
	startTime := time.Now()
	rows, err := conn.db.Raw(sql).Rows()
	if err != nil {
		return &Result{
			Err: errors.New("Connector.ExecSQL: execute '" + sql + "' error: " + err.Error()),
		}
	}
	defer rows.Close()

	result := &Result{
		ColumnNames: make([]string, 0),
		ColumnTypes: make([]string, 0),
		Rows: make([][]string, 0),
		Err: nil,
	}
	for rows.Next() {
		columnTypes, err := rows.ColumnTypes()
		if err != nil {
			return &Result{
				Err: errors.New("Connector.ExecSQL: get columns' type error: " + err.Error()),
			}
		}
		if len(result.ColumnNames) == 0 {
			for _, columnType := range columnTypes {
				result.ColumnNames = append(result.ColumnNames, columnType.Name())
				result.ColumnTypes = append(result.ColumnTypes, columnType.DatabaseTypeName())
			}
		} else {
			if len(columnTypes) != len(result.ColumnNames) {
				return &Result{
					Err: errors.New("Connector.ExecSQL: column mismatch: " +
						"len(columnTypes) != len(result.ColumnNames)"),
				}
			}
			for i, columnType := range columnTypes {
				if columnType.Name() != result.ColumnNames[i] {
					return &Result{
						Err: errors.New("Connector.ExecSQL: column mismatch: " +
							"columnType.Name() != result.ColumnNames[i]"),
					}
				}
				if columnType.DatabaseTypeName() != result.ColumnTypes[i] {
					return &Result{
						Err: errors.New("Connector.ExecSQL: column mismatch: " +
							"columnType.DatabaseTypeName() != result.ColumnTypes[i]"),
					}
				}
			}
		}

		data := make([]string, len(columnTypes))
		dataI := make([]interface{}, len(columnTypes))
		for i, _ := range data {
			dataI[i] = &data[i]
		}
		err = rows.Scan(dataI...)
		if err != nil {
			return &Result{
				Err: errors.New("Connector.ExecSQL: scan row error: " + err.Error()),
			}
		}
		result.Rows = append(result.Rows, data)
	}

	result.Time = time.Since(startTime)
	return result
}