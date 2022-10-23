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

// addRdMOnU: RdMOnU, *ast.Join: ON xxx -> ON TRUE | ON (xxx) OR 1
func (v *MutateVisitor) addRdMOnU(in *ast.Join, flag int) {
	if in.On != nil {
		v.addCandidate(RdMOnU, 1, in, flag)
	}
}

// doRdMOnU: RdMOnU, *ast.Join: ON xxx -> ON TRUE | ON (xxx) OR 1
func doRdMOnU(rootNode ast.Node, in ast.Node, seed int64) ([]byte, error) {
	rander := rand.New(rand.NewSource(seed))
	switch in.(type) {
	case *ast.Join:
		join := in.(*ast.Join)
		// check
		if join.On == nil {
			return nil, errors.New("doRdMOnU: join.On == nil")
		}
		// mutate
		old := join.On.Expr

		rd := rander.Intn(2)
		if rd == 0 {
			// HAVING xxx -> HAVING TRUE
			join.On.Expr = &test_driver.ValueExpr{
				Datum: test_driver.NewDatum(1),
			}
		} else {
			// HAVING xxx -> HAVING (xxx) OR 1
			join.On.Expr = &ast.BinaryOperationExpr{
				Op: opcode.LogicOr,
				L: old,
				R: &test_driver.ValueExpr{
					Datum: test_driver.NewDatum(1),
				},
			}
		}

		sql, err := restore(rootNode)
		if err != nil {
			return nil, errors.New("doRdMOnU: " +  err.Error())
		}
		// recover
		join.On.Expr = old
		return sql, nil
	case nil:
		return nil, errors.New("doRdMOnU: type error: nil")
	default:
		return nil, errors.New("doRdMOnU: type error: " + reflect.TypeOf(in).String())
	}
}