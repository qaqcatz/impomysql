package tool

import (
	"github.com/pkg/errors"
	"github.com/qaqcatz/impomysql/connector"
	"strconv"
	"strings"
)

// DropDBLike: drop all databases matching like.
// We will first:
//   SHOW DATABASES LIKE 'like'.
// Then for each dbname, we will:
//   DROP DATABASE dbname.
// For example:
//   like: TEST%
//   ->
//   drop database TEST, TEST0, TEST1, TEST2, ...
// dsn format:
//   username$password$host$port$dbname
//   Obviously you cannot use '$' in any of username, password, host, port, dbname
//   In DropDBLike, we will ignore your dbname
func DropDBLike(dsn string, like string) error {
	dsnUnits := strings.Split(dsn, "$")
	if len(dsnUnits) != 5 {
		return errors.New("[DropDBLike]len(dsnUnits) != 5")
	}
	username := dsnUnits[0]
	password := dsnUnits[1]
	host := dsnUnits[2]
	portStr := dsnUnits[3]
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return errors.Wrap(err, "[DropDBLike]parse port error")
	}
	conn, err := connector.NewConnector(host, port, username, password, "")
	if err != nil {
		return err
	}
	showResult := conn.ExecSQL(`SHOW DATABASES LIKE '`+like+`'`)
	if showResult.Err != nil {
		return showResult.Err
	}
	for _, row := range showResult.Rows {
		if len(row) != 1 {
			return errors.New("[DropDBLike]len(row) != 1")
		}
		dbname := row[0]
		dropResult := conn.ExecSQL(`DROP DATABASE `+dbname)
		if dropResult.Err != nil {
			return dropResult.Err
		}
	}
	return nil
}
