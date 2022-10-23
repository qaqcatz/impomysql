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

// addRdMOnL: RdMOnL, *ast.Join: ON xxx -> ON FALSE | ON (xxx) AND 0
func (v *MutateVisitor) addRdMOnL(in *ast.Join, flag int) {
	if in.On != nil {
		v.addCandidate(RdMOnL, 0, in, flag)
	}
}

// doRdMOnL: RdMOnL, *ast.Join: ON xxx -> ON FALSE | ON (xxx) AND 0
func doRdMOnL(rootNode ast.Node, in ast.Node, seed int64) ([]byte, error) {
	rander := rand.New(rand.NewSource(seed))
	switch in.(type) {
	case *ast.Join:
		join := in.(*ast.Join)
		// check
		if join.On == nil {
			return nil, errors.New("doRdMOnL: join.On == nil")
		}
		// mutate
		old := join.On.Expr

		rd := rander.Intn(2)
		if rd == 0 {
			// ON xxx -> ON FALSE
			join.On.Expr = &test_driver.ValueExpr{
				Datum: test_driver.NewDatum(0),
			}
		} else {
			// ON xxx -> ON (xxx) AND 0
			join.On.Expr = &ast.BinaryOperationExpr{
				Op: opcode.LogicAnd,
				L: old,
				R: &test_driver.ValueExpr{
					Datum: test_driver.NewDatum(0),
				},
			}
		}

		sql, err := restore(rootNode)
		if err != nil {
			return nil, errors.New("doRdMOnL: " +  err.Error())
		}
		// recover
		join.On.Expr = old
		return sql, nil
	case nil:
		return nil, errors.New("doRdMOnL: type error: nil")
	default:
		return nil, errors.New("doRdMOnL: type error: " + reflect.TypeOf(in).String())
	}
}