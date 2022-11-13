package task

import (
	"github.com/pkg/errors"
	"github.com/qaqcatz/impomysql/connector"
	"strconv"
)

type RdSql struct {
	Id  int    `json:"id"`
	Sql string `json:"sql"`
}

// ExtractSQL: Extract sql statements by ';':
//   - ignore the ';' in ``, '', "";
//   - ignore the escaped characters in ``, '', "";
// Note that: Comments cannot have ';'
func ExtractSQL(s string) []*RdSql {
	res := make([]*RdSql, 0)
	start := 0
	flag := -1
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '\'':
			if flag == -1 {
				flag = '\''
			} else {
				if flag == '\'' {
					flag = -1
				}
			}
		case '"':
			if flag == -1 {
				flag = '"'
			} else {
				if flag == '"' {
					flag = -1
				}
			}
		case '`':
			if flag == -1 {
				flag = '`'
			} else {
				if flag == '`' {
					flag = -1
				}
			}
		case '\\':
			if flag != -1 {
				i++
			}
		case ';':
			if flag == -1 {
				res = append(res, &RdSql{
					Id:  len(res),
					Sql: s[start : i+1],
				})
				start = i + 1
			}
		default:
			continue
		}
	}
	return res
}

// InitDDLSqls: init database and execute ddl sqls
func InitDDLSqls(ddlSqls []*RdSql, conn *connector.Connector) error {
	err := conn.InitDB()
	if err != nil {
		return errors.Wrap(err, "[InitDDLSqls]init database error")
	}
	for i, ddlSql := range ddlSqls {
		result := conn.ExecSQL(ddlSql.Sql)
		if result.Err != nil {
			return errors.Wrap(result.Err, "[InitDDLSqls]exec ddl sql " + strconv.Itoa(i) + " error")
		}
	}
	return nil
}