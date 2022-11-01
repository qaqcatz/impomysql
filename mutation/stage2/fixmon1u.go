package stage2

import (
	"errors"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/test_driver"
	_ "github.com/pingcap/tidb/parser/test_driver"
	"reflect"
)

// addFixMOn1U: FixMOn1U, *ast.Join: ON xxx -> ON 1
func (v *MutateVisitor) addFixMOn1U(in *ast.Join, flag int) {
	if in.On != nil {
		v.addCandidate(FixMOn1U, 1, in, flag)
	}
}

// doFixMOn1U: FixMOn1U, *ast.Join: ON xxx -> ON 1
func doFixMOn1U(rootNode ast.Node, in ast.Node) ([]byte, error) {
	switch in.(type) {
	case *ast.Join:
		join := in.(*ast.Join)
		// check
		if join.On == nil {
			return nil, errors.New("FixMOn1U: join.On == nil")
		}
		// mutate
		old := join.On.Expr

		// ON xxx -> ON 1
		join.On.Expr = &test_driver.ValueExpr{
			Datum: test_driver.NewDatum(1),
		}

		sql, err := restore(rootNode)
		if err != nil {
			return nil, errors.New("FixMOn1U: " +  err.Error())
		}
		// recover
		join.On.Expr = old
		return sql, nil
	case nil:
		return nil, errors.New("FixMOn1U: type error: nil")
	default:
		return nil, errors.New("FixMOn1U: type error: " + reflect.TypeOf(in).String())
	}
}