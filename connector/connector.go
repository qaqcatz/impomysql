package connector

import (
	"bytes"
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os/exec"
	"strconv"
	"time"
)

// Connector: connect to MySQL, execute raw sql statements, return raw execution result or error.
type Connector struct {
	DSN             string
	Host            string
	Port            int
	Username        string
	Password        string
	DbName          string
	MysqlClientPath string
	db              *gorm.DB
}

// NewConnector: create Connector. CREATE DATABASE IF NOT EXISTS dbname
//
// Default mysqlClientPath: /usr/bin/mysql
func NewConnector(host string, port int, username string, password string, dbname string, mysqlClientPath string) (*Connector, error) {
	if mysqlClientPath == "" {
		mysqlClientPath = "/usr/bin/mysql"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		username, password, host, port, "")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		return nil, errors.New("NewConnector: create Connector error: " + err.Error())
	}
	conn := &Connector{
		DSN:             dsn,
		Host:            host,
		Port:            port,
		Username:        username,
		Password:        password,
		DbName:          dbname,
		MysqlClientPath: mysqlClientPath,
		db:              db,
	}
	// CREATE DATABASE IF NOT EXISTS conn.DbName
	result := conn.ExecSQL("CREATE DATABASE IF NOT EXISTS " + conn.DbName)
	if result.Err != nil {
		return nil, errors.New("NewConnector: create database if not exists error: " + result.Err.Error())
	}
	return conn, nil
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

// FlatRows: [["1","2"],["3","4"]] -> ["1,2", "3,4"]
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

// IsEmpty: if the result is empty
func (result *Result) IsEmpty() bool {
	return len(result.ColumnNames) == 0
}

// ExecSQL: execute sql, return *Result.
func (conn *Connector) ExecSQL(sql string) *Result {
	startTime := time.Now()
	rows, err := conn.db.Raw(sql).Rows()
	if err != nil {
		return &Result{
			Err: errors.New("Connector.ExecSQL: execute error: " + err.Error()),
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

		// gorm cannot convert NULL to string, we should use []byte
		data := make([][]byte, len(columnTypes))
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

		dataS := make([]string, len(columnTypes))
		for i, _ := range data {
			if data[i] == nil {
				dataS[i] = "NULL"
			} else {
				dataS[i] = string(data[i])
			}
		}
		result.Rows = append(result.Rows, dataS)
	}

	result.Time = time.Since(startTime)
	return result
}

// ExecSQLS: execute sql, return *Result.
//
// There is a bug in golang mysql driver:
//
// If you execute the following sql in mysql-client, you will see:
//   mysql> select 9223372036854775807 + 1 > 1;
//   ERROR 1690 (22003): BIGINT value is out of range in '(9223372036854775807 + 1)'
// However, when execute this sql in gorm, no error, just an empty result.
//
// We will double check the result:
//
// when sql1 returns an empty result and no error occurred, change it to sql2=SELECT EXISTS (sql1),
// if sql2 returns 0, then sql1 is really empty, otherwise an error occurred.
//
// Therefore, it is recommended to use ExecSQLS to query.
// Note that do not use this function for other ddl/dml, only use it to query without side effects!
func (conn *Connector) ExecSQLS(sql string) *Result {
	res1 := conn.ExecSQL(sql)
	if res1.Err == nil && res1.IsEmpty() {
		sql2 := "SELECT EXISTS (" + sql + " )"
		res2 := conn.ExecSQL(sql2)
		if res2.Err == nil && len(res2.Rows) == 1 && len(res2.Rows[0]) == 1 && res2.Rows[0][0] == "0" {
			return res1
		} else {
			res1.Err = errors.New("ExecSQLS: unknown error, maybe BIGINT value is out of range. ")
			return res1
		}
	} else {
		return res1
	}
}

// ExecSQLX: see ExecSQLS first.
// Unfortunately, sometimes golang mysql driver will return a non-empty result, while mysql-client will return an error.
// Therefore we will eventually use mysql-client to execute the sql and check for errors.
// Note that this function is very slow, only use it to verify bugs.
//
// Actually, we execute sql | Connector.MysqlClientPath -h Connector.Host -P Connector.Port
// -u Connector.Username --password=Connector.Password Connector.DbName by pipeline, and
// return output stream, error stream, error
func (conn *Connector) ExecSQLX(sql string) (string, string, error) {
	sqlBuf := bytes.NewBufferString(sql)

	mysqlClient := exec.Command("/bin/bash", "-c",
		conn.MysqlClientPath+
		" -h " + conn.Host +
		" -P " + strconv.Itoa(conn.Port) +
		" -u " + conn.Username +
		" -p" + conn.Password +
		" " + conn.DbName)

	readPipe, err := mysqlClient.StdinPipe()
	if err != nil {
		return "", "", errors.New("ExecSQLX: mysqlClient.StdinPipe() error: " + err.Error())
	}

	_, err = sqlBuf.WriteTo(readPipe)
	if err != nil {
		return "", "", errors.New("ExecSQLX: sqlBuf.WriteTo(readPipe) error: " + err.Error())
	}

	var outBuf bytes.Buffer
	var errBuf bytes.Buffer
	mysqlClient.Stdout = &outBuf
	mysqlClient.Stderr = &errBuf

	err = mysqlClient.Start()
	if err != nil {
		return "", "", errors.New("ExecSQLX: mysqlClient.Start() error: " + err.Error())
	}

	err = readPipe.Close()
	if err != nil {
		return "", "", errors.New("ExecSQLX: readPipe.Close() error: " + err.Error())
	}

	err = mysqlClient.Wait()
	if err != nil {
		return outBuf.String(), errBuf.String(), errors.New("ExecSQLX: mysqlClient.Wait() error: " + err.Error())
	}
	return outBuf.String(), errBuf.String(), nil
}

// InitDBTEST:
//   DROP DATABASE IF EXISTS Connector.DbName
//   CREATE DATABASE Connector.DbName
func (conn *Connector) InitDB() error {
	result := conn.ExecSQL("DROP DATABASE IF EXISTS " + conn.DbName)
	if result.Err != nil {
		return result.Err
	}
	result = conn.ExecSQL("CREATE DATABASE " + conn.DbName)
	if result.Err != nil {
		return result.Err
	}
	result = conn.ExecSQL("USE " + conn.DbName)
	if result.Err != nil {
		return result.Err
	}
	return nil
}

// RmDB:
//   DROP DATABASE IF EXISTS Connector.DbName
func (conn *Connector) RmDB() error {
	result := conn.ExecSQL("DROP DATABASE IF EXISTS " + conn.DbName)
	if result.Err != nil {
		return result.Err
	}
	return nil
}