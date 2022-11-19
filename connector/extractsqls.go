package connector

type EachSql struct {
	Id  int    `json:"id"`
	Sql string `json:"sql"`
}

// ExtractSQL: Extract sql statements by ';':
//   - ignore the ';' in ``, '', "";
//   - ignore the escaped characters in ``, '', "";
// Note that: Comments cannot have ';'
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