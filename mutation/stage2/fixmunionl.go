package stage2

import (
	"errors"
	_ "github.com/pingcap/tidb/parser/test_driver"
	"github.com/pingcap/tidb/parser/ast"
	"reflect"
)

// addFixMUnionL: FixMUnionL, *ast.SetOprSelectList: remove Selects[1:]
func (v *MutateVisitor) addFixMUnionL(in *ast.SetOprSelectList, flag int) {
	if in.Selects != nil && len(in.Selects) > 1 {
		v.addCandidate(FixMUnionL, 0, in, flag)
	}
}

// doFixMUnionL: FixMUnionL, *ast.SetOprSelectList: remove Selects[1:]
func doFixMUnionL(rootNode ast.Node, in ast.Node) ([]byte, error) {
	switch in.(type) {
	case *ast.SetOprSelectList:
		lst := in.(*ast.SetOprSelectList)
		// check
		if lst.Selects == nil || len(lst.Selects) <= 1 {
			return nil, errors.New("doFixMUnionL: lst.Selects == nil || len(lst.Selects) <= 1")
		}
		// mutate
		oldSels := lst.Selects
		newSels := make([]ast.Node, 0)
		newSels = append(newSels, oldSels[0])
		lst.Selects = newSels
		sql, err := restore(rootNode)
		if err != nil {
			return nil, errors.New("doFixMUnionL: " +  err.Error())
		}
		// recover
		lst.Selects = oldSels
		return sql, nil
	case nil:
		return nil, errors.New("doFixMUnionL: type error: nil")
	default:
		return nil, errors.New("doFixMUnionL: type error: " + reflect.TypeOf(in).String())
	}
}