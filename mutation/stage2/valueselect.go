package stage2

import (
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/test_driver"
	_ "github.com/pingcap/tidb/parser/test_driver"
	"math/rand"
)

func GenValueSelect(columnNum int, seed int64) *ast.SelectStmt {
	rander := rand.New(rand.NewSource(seed))

	selectFields := make([]*ast.SelectField, 0)
	for i := 0; i < columnNum; i++ {
		var data interface{}
		rd := rander.Intn(3)
		switch rd {
		case 0:
			data = rander.Int63()
		case 1:
			data = rander.Uint64()
		case 2:
			data = rander.Float64()
		}

		var expr ast.ExprNode = nil
		if rd == 0 {
			expr =  &test_driver.ValueExpr{
				Datum: test_driver.NewDatum(data),
			}
		}

		selectField := &ast.SelectField {
			Expr: expr,
		}

		selectFields = append(selectFields, selectField)
	}

	var fieldlist *ast.FieldList = &ast.FieldList {
		Fields: selectFields,
	}

	sel := &ast.SelectStmt{
		Fields: fieldlist,
	}
	return sel
}