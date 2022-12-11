package sqlsimx

import (
	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/test_driver"
	"github.com/pkg/errors"
	"github.com/qaqcatz/impomysql/connector"
)

// rmFields: remove unused fields
func rmFields(sql string, result *connector.Result, conn *connector.Connector) (string, error) {
	p := parser.New()
	stmtNodes, _, err := p.Parse(sql, "", "")
	if err != nil {
		return "", errors.Wrap(err, "[rmFields]parse error")
	}
	if stmtNodes == nil || len(stmtNodes) == 0 {
		return "", errors.New("[rmFields]stmtNodes == nil || len(stmtNodes) == 0 ")
	}
	rootNode := &stmtNodes[0]

	v := &rmFieldsVisitor{
		rootNode: *rootNode,
		result: result,
		conn: conn,
	}
	(*rootNode).Accept(v)

	simplifiedSql, err := restore(*rootNode)
	if err != nil {
		return "", errors.Wrap(err, "[rmFields]restore error")
	}

	tempSql := string(simplifiedSql)
	tempResult := conn.ExecSQL(tempSql)
	cmp, err := tempResult.CMP(result)
	if err == nil && cmp == 0 {
		return tempSql, nil
	}

	return sql, nil
}

type rmFieldsVisitor struct {
	rootNode ast.Node
	result *connector.Result
	conn *connector.Connector
}

func (v *rmFieldsVisitor) Enter(in ast.Node) (ast.Node, bool) {
	switch in.(type) {
	case *ast.FieldList:
		fieldList := in.(*ast.FieldList)
		if fieldList.Fields != nil && len(fieldList.Fields) > 1 {

			newFields := make([]*ast.SelectField, 0)

			for i := 0; i < len(fieldList.Fields); i += 1 {
				// dp copy
				oldFields := append([]*ast.SelectField{}, fieldList.Fields...)

				fieldList.Fields = append(fieldList.Fields[:i], fieldList.Fields[i+1:]...)

				simplifiedSql, err := restore(v.rootNode)

				fieldList.Fields = oldFields

				if err == nil {
					tempSql := string(simplifiedSql)
					tempResult := v.conn.ExecSQL(tempSql)
					cmp, err := tempResult.CMP(v.result)
					if err == nil && cmp == 0 {
						continue
					}
				}

				newFields = append(newFields, fieldList.Fields[i])
			}

			if len(newFields) == 0 {
				fieldList.Fields = []*ast.SelectField{
					&ast.SelectField{
						WildCard: nil,
						Expr: &test_driver.ValueExpr{
							Datum: test_driver.NewDatum(1),
						},
					},
				}
			} else {
				fieldList.Fields = newFields
			}
		}
	}
	return in, false
}

func (v *rmFieldsVisitor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}