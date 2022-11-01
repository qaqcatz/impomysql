package stage2

import (
	"errors"
	"github.com/pingcap/tidb/parser/ast"
	_ "github.com/pingcap/tidb/parser/test_driver"
	"math/rand"
	"reflect"
)

// addRdMInL: RdMInL, *ast.PatternInExpr: in(x,x,x,...) -> in(x,x,x)
// may false positive, skim
func (v *MutateVisitor) addRdMInL(in *ast.PatternInExpr, flag int) {
	if in.Sel == nil && in.List != nil && len(in.List) >= 2 {
		v.addCandidate(RdMInL, 0, in, flag)
	}
}

// doRdMInL: RdMInL, *ast.PatternInExpr: in(x,x,x,...) -> in(x,x,x)
// may false positive, skim
func doRdMInL(rootNode ast.Node, in ast.Node, seed int64) ([]byte, error) {
	rander := rand.New(rand.NewSource(seed))
	switch in.(type) {
	case *ast.PatternInExpr:
		pin := in.(*ast.PatternInExpr)
		// check
		if pin.Sel != nil || pin.List == nil || len(pin.List) < 2 {
			return nil, errors.New("doRdMInL: pin.Sel != nil || pin.List == nil || len(pin.List) < 2")
		}
		// mutate
		oldList := pin.List
		// random remove [1,len(pin.List)-1] expr
		newList := make([]ast.ExprNode, 0)
		p := 0
		for {
			p = rander.Intn(len(oldList)-p)+p
			newList = append(newList, oldList[p])
			p += 1
			if p >= len(oldList) {
				break
			}
		}
		pin.List = newList
		sql, err := restore(rootNode)
		if err != nil {
			return nil, errors.New("doRdMInL: " +  err.Error())
		}
		// recover
		pin.List = oldList
		return sql, nil
	case nil:
		return nil, errors.New("doRdMInL: type error: nil")
	default:
		return nil, errors.New("doRdMInL: type error: " + reflect.TypeOf(in).String())
	}
}
