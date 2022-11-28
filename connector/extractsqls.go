package connector

import (
	"github.com/pkg/errors"
	"io/ioutil"
)

type EachSql struct {
	Id  int    `json:"id"`
	Sql string `json:"sql"`
}

// ExtractSQL: s is a sqls string, each sql statement is separated by ';' in s.
// We will extract each sql statement into []*EachSql.
//
// Note that:
//   - we will ignore the ';' in ``, '', "";
//   - we will ignore the escaped characters in ``, '', "";
//   - your comments cannot have ';'
func ExtractSQL(s string) []*EachSql {
	res := make([]*EachSql, 0)
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
				res = append(res, &EachSql{
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

// ExtractSqlFromPath: extract sql from sql file, see ExtractSQL
func ExtractSqlFromPath(sqlPath string) ([]*EachSql, error) {
	sqlData, err := ioutil.ReadFile(sqlPath)
	if err != nil {
		return nil, errors.Wrap(err, "[ExtractSqlFromPath]read file error")
	}
	return ExtractSQL(string(sqlData)), nil
}