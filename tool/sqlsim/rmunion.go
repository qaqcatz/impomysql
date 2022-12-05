package sqlsim

import (
	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	_ "github.com/pingcap/tidb/parser/test_driver"
	"github.com/pkg/errors"
	"github.com/qaqcatz/impomysql/connector"
	"github.com/qaqcatz/impomysql/mutation/oracle"
	"github.com/qaqcatz/impomysql/task"
)

// rmUnion: A UNION B -> A or B
func rmUnion(bug *task.BugReport, conn *connector.Connector) error {
	sql2 := []*string {
		&(bug.OriginalSql),
		&(bug.MutatedSql),
	}
	res2 := []**connector.Result {
		&(bug.OriginalResult),
		&(bug.MutatedResult),
	}

	// (i, j)
	// i: 0-front, 1-back;
	// j: 0-original, 1-mutated;
	union := [][]string {
		[]string {
			"",
			"",
		},[]string {
			"",
			"",
		},
	}
	for j := 0; j < 2; j++ {
		p := parser.New()
		stmtNodes, _, err := p.Parse(*sql2[j], "", "")
		if err != nil {
			return errors.Wrap(err, "[rmUnion]parse error")
		}
		if stmtNodes == nil || len(stmtNodes) == 0 {
			return errors.New("[rmUnion]stmtNodes == nil || len(stmtNodes) == 0 ")
		}
		rootNode := &stmtNodes[0]

		switch (*rootNode).(type) {
		case *ast.SetOprStmt:
			sos := (*rootNode).(*ast.SetOprStmt)
			sels := sos.SelectList.Selects
			if len(sels) != 2 {
				return nil
			}
			for i := 0; i < 2; i++ {
				simplifiedSql, err := restore(sels[i])
				if err != nil {
					return errors.Wrap(err, "[rmUnion]restore error")
				}
				union[i][j] = string(simplifiedSql)
			}
		default:
			return nil
		}
	}
	for i := 0; i < 2; i++ {
		sqlX := union[i][0]
		resX := conn.ExecSQL(sqlX)
		if resX.Err != nil {
			continue
		}
		sqlY := union[i][1]
		resY := conn.ExecSQL(sqlY)
		if resY.Err != nil {
			continue
		}
		check, err := oracle.Check(resX, resY, bug.IsUpper)
		if err != nil {
			continue
		}
		if !check {
			*sql2[0] = sqlX
			*res2[0] = resX
			*sql2[1] = sqlY
			*res2[1] = resY
			break
		}
	}
	return nil
}
