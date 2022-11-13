package connector

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
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
	db              *sql.DB
}

// NewConnector: create Connector. CREATE DATABASE IF NOT EXISTS dbname + USE dbname when dbname != ""
func NewConnector(host string, port int, username string, password string, dbname string) (*Connector, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		username, password, host, port, "")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "[NewConnector]open dsn error")
	}
	conn := &Connector{
		DSN:             dsn,
		Host:            host,
		Port:            port,
		Username:        username,
		Password:        password,
		DbName:          dbname,
		db:              db,
	}
	if dbname != "" {
		// CREATE DATABASE IF NOT EXISTS conn.DbName
		result := conn.ExecSQL("CREATE DATABASE IF NOT EXISTS " + conn.DbName)
		if result.Err != nil {
			return nil, result.Err
		}
		// USE conn.DbName
		result = conn.ExecSQL("USE " + conn.DbName)
		if result.Err != nil {
			return nil, result.Err
		}
	}
	return conn, nil
}

// Connector.ExecSQL: execute sql, return *Result.
func (conn *Connector) ExecSQL(sql string) *Result {
	startTime := time.Now()
	rows, err := conn.db.Query(sql)
	if err != nil {
		return &Result{
			Err: errors.Wrap(err, "[Connector.ExecSQL]execute sql error"),
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
				Err: errors.Wrap(err, "[Connector.ExecSQL]get column type error"),
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
					Err: errors.New("[Connector.ExecSQL]|columnTypes|("+strconv.Itoa(len(columnTypes))+") != "+
						"|columnNames|("+strconv.Itoa(len(result.ColumnNames))+")"),
				}
			}
			for i, columnType := range columnTypes {
				if columnType.Name() != result.ColumnNames[i] {
					return &Result{
						Err: errors.New("[Connector.ExecSQL]columnType.Name()("+columnType.Name()+") != "+
							"result.ColumnNames[i]("+result.ColumnNames[i]+")"),
					}
				}
				if columnType.DatabaseTypeName() != result.ColumnTypes[i] {
					return &Result{
						Err: errors.New("[Connector.ExecSQL]columnType.DatabaseTypeName()("+columnType.DatabaseTypeName()+") != "+
							"result.ColumnTypes[i]("+result.ColumnTypes[i]+")"),
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
				Err: errors.Wrap(err, "[Connector.ExecSQL]scan rows error"),
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

	if rows.Err() != nil {
		return &Result{
			Err: errors.Wrap(rows.Err(), "[Connector.ExecSQL]rows error"),
		}
	}

	result.Time = time.Since(startTime)
	return result
}

// Connector.InitDB:
//   DROP DATABASE IF EXISTS Connector.DbName
//   CREATE DATABASE Connector.DbName
//   USE Connector.DbName
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

func (conn *Connector) Close() {
	_ = conn.db.Close()
}