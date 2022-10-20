package stage2

import (
	"errors"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/opcode"
	"github.com/pingcap/tidb/parser/test_driver"
	_ "github.com/pingcap/tidb/parser/test_driver"
	"math/rand"
	"reflect"
)

// addRdMWhereU: RdMWhereU, *ast.SelectStmt: WHERE xxx -> WHERE TRUE | WHERE (xxx) OR 1
func (v *MutateVisitor) addRdMWhereU(in *ast.SelectStmt, flag int) {
	if in.Where != nil {
		v.addCandidate(RdMWhereU, 1, in, flag)
	}
}

// doRdMWhereU: RdMWhereU, *ast.SelectStmt: WHERE xxx -> WHERE TRUE | WHERE (xxx) OR 1
func doRdMWhereU(rootNode ast.Node, in ast.Node, seed int64) ([]byte, error) {
	rander := rand.New(rand.NewSource(seed))
	switch in.(type) {
	case *ast.SelectStmt:
		sel := in.(*ast.SelectStmt)
		// check
		if sel.Where == nil {
			return nil, errors.New("doRdMWhereU: sel.Where == nil")
		}
		// mutate
		old := sel.Where

		rd := rander.Intn(2)
		if rd == 0 {
			// WHERE xxx -> WHERE TRUE
			sel.Where = &test_driver.ValueExpr{
				Datum: test_driver.NewDatum(1),
			}
		} else {
			// WHERE xxx -> WHERE (xxx) OR 1
			sel.Where = &ast.BinaryOperationExpr{
				Op: opcode.LogicOr,
				L: old,
				R: &test_driver.ValueExpr{
					Datum: test_driver.NewDatum(1),
				},
			}
		}

		sql, err := restore(rootNode)
		if err != nil {
			return nil, errors.New("doRdMWhereU: " +  err.Error())
		}
		// recover
		sel.Where = old
		return sql, nil
	case nil:
		return nil, errors.New("doRdMWhereU: type error: nil")
	default:
		return nil, errors.New("doRdMWhereU: type error: " + reflect.TypeOf(in).String())
	}
}