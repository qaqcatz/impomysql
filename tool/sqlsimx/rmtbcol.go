package sqlsimx

import (
	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	_ "github.com/pingcap/tidb/parser/test_driver"
	"github.com/pkg/errors"
	"github.com/qaqcatz/impomysql/connector"
	"strings"
)

// rmtbcol: rm unused tables, columns
func rmtbcol(dml string, ddls []*connector.EachSql) (string, error) {
	newDDls := ""

	tableNameMap, columnNameMap, err := getTbCol(dml)
	if err != nil {
		return "", err
	}

	table_svColIds_map := make(map[string][]int)

	for _, ddl := range ddls {
		p := parser.New()
		stmtNodes, _, err := p.Parse(ddl.Sql, "", "")
		if err != nil {
			return "", errors.Wrap(err, "[rmtbcol]parse error")
		}
		if stmtNodes == nil || len(stmtNodes) == 0 {
			return "", errors.New("[rmtbcol]stmtNodes == nil || len(stmtNodes) == 0 ")
		}
		rootNode := &stmtNodes[0]

		switch (*rootNode).(type) {
		case *ast.CreateTableStmt:
			create := (*rootNode).(*ast.CreateTableStmt)
			tableName := strings.ToLower(create.Table.Name.O)
			if _, ok := tableNameMap[tableName]; !ok {
				continue
			}
			newDDls += "drop table if exists " + tableName + ";\n"

			// remove refer
			create.ReferTable = nil
			// remove select
			create.Select = nil
			// remove constraint
			create.Constraints = nil
			// remove options
			create.Options = nil
			// remove partition
			create.Partition = nil

			// remove cols
			var svColIds []int = nil
			newCols := []*ast.ColumnDef{}
			for i, colDef := range create.Cols {
				colName := strings.ToLower(colDef.Name.Name.O)
				if _, ok := columnNameMap[colName]; ok {
					newCols = append(newCols, colDef)
					svColIds = append(svColIds, i)
				}
			}
			create.Cols = newCols
			table_svColIds_map[tableName] = svColIds

			simplifiedSql, err := restore(*rootNode)
			if err != nil {
				return "", errors.Wrap(err, "[rmtbcol]restore error")
			}

			newDDls += sqlSemi(string(simplifiedSql)) + "\n"
		case *ast.InsertStmt:
			insert := (*rootNode).(*ast.InsertStmt)
			tableRef := insert.Table.TableRefs
			if tableRef.Right != nil {
				continue
			}
			ok := false
			var tbs *ast.TableSource
			if tbs, ok = (tableRef.Left).(*ast.TableSource); !ok {
				continue
			}
			var tbn *ast.TableName
			if tbn, ok = (tbs.Source).(*ast.TableName); !ok {
				continue
			}
			tableName := strings.ToLower(tbn.Name.O)
			if _, ok := tableNameMap[tableName]; !ok {
				continue
			}

			var svColIds []int = nil
			if svColIds, ok = table_svColIds_map[tableName]; !ok {
				return "", errors.New("[rmtbcol]please check table_rmColIds_map!")
			}

			// remove columns
			insert.Columns = nil
			// remove set list
			insert.Setlist = nil
			// remove select
			insert.Select = nil
			// remove table hints
			insert.TableHints = nil
			// remove duplicate
			insert.OnDuplicate = nil
			// remove partition
			insert.PartitionNames = nil

			listsLen := len(insert.Lists)
			svLen := len(svColIds)
			for i := 0; i < listsLen; i += 1 {
				newList := make([]ast.ExprNode, svLen)
				for j := 0; j < svLen; j += 1 {
					newList[j] = insert.Lists[i][svColIds[j]]
				}
				insert.Lists[i] = newList
			}

			simplifiedSql, err := restore(*rootNode)
			if err != nil {
				return "", errors.Wrap(err, "[rmtbcol]restore error")
			}

			newDDls += sqlSemi(string(simplifiedSql)) + "\n"
		default:
			newDDls += sqlSemi(strings.TrimSpace(ddl.Sql)) + "\n"
		}
	}
	return newDDls, nil
}

func sqlSemi(sql string) string {
	if strings.HasSuffix(sql, ";") {
		return sql
	}
	return sql + ";"
}

// getTbCol: get all tables, columns in sql
func getTbCol(sql string) (map[string]bool, map[string]bool, error) {
	p := parser.New()
	stmtNodes, _, err := p.Parse(sql, "", "")
	if err != nil {
		return nil, nil, errors.Wrap(err, "[rmFields]parse error")
	}
	if stmtNodes == nil || len(stmtNodes) == 0 {
		return nil, nil, errors.New("[rmFields]stmtNodes == nil || len(stmtNodes) == 0 ")
	}
	rootNode := &stmtNodes[0]

	v := &getTbColVisitor{
		TableNameMap: make(map[string]bool),
		ColumnNameMap: make(map[string]bool),
	}
	(*rootNode).Accept(v)

	return v.TableNameMap, v.ColumnNameMap, nil
}

type getTbColVisitor struct {
	TableNameMap map[string]bool
	ColumnNameMap map[string]bool
}

func (v *getTbColVisitor) Enter(in ast.Node) (ast.Node, bool) {
	switch in.(type) {
	case *ast.TableName:
		tableName := in.(*ast.TableName)
		v.TableNameMap[strings.ToLower(tableName.Name.O)] = true
	case *ast.ColumnName:
		columnName := in.(*ast.ColumnName)
		v.ColumnNameMap[strings.ToLower(columnName.Name.O)] = true
	}
	return in, false
}

func (v *getTbColVisitor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}