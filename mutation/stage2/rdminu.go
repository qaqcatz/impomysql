package stage2

import (
	"errors"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/test_driver"
	_ "github.com/pingcap/tidb/parser/test_driver"
	"math/rand"
	"reflect"
)

// addRdMInU: RdMInU, *ast.PatternInExpr: in(x,x,x) -> in(x,x,x,...)
// type conversion, only support in(x,x,x) -> in(x,x,x,null)
func (v *MutateVisitor) addRdMInU(in *ast.PatternInExpr, flag int) {
	if in.Sel == nil && in.List != nil {
		v.addCandidate(RdMInU, 1, in, flag)
	}
}

// doRdMInU: RdMInU, *ast.PatternInExpr: in(x,x,x) -> in(x,x,x,...)
// type conversion, only support in(x,x,x) -> in(x,x,x,null)
func doRdMInU(rootNode ast.Node, in ast.Node, seed int64) ([]byte, error) {
	rander := rand.New(rand.NewSource(seed))
	switch in.(type) {
	case *ast.PatternInExpr:
		pin := in.(*ast.PatternInExpr)
		// check
		if pin.Sel != nil || pin.List == nil {
			return nil, errors.New("doRdMInU: pin.Sel != nil || pin.List == nil")
		}
		// mutate
		oldList := pin.List
		newList := make([]ast.ExprNode, 0)
		for _, expr := range oldList {
			newList = append(newList, expr)
		}
		// add 1 ~ 3 null expr
		nullExprNum := rander.Intn(3)+1
		for i := 0; i < nullExprNum; i++ {
			newList = append(newList, &test_driver.ValueExpr{
				Datum: test_driver.NewDatum(nil),
			})
		}
		pin.List = newList
		sql, err := restore(rootNode)
		if err != nil {
			return nil, errors.New("doRdMInU: " +  err.Error())
		}
		// recover
		pin.List = oldList
		return sql, nil
	case nil:
		return nil, errors.New("doRdMInU: type error: nil")
	default:
		return nil, errors.New("doRdMInU: type error: " + reflect.TypeOf(in).String())
	}
}
