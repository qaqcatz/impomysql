package sqlsimx

import (
	"bytes"
	"github.com/pingcap/tidb/parser/test_driver"
	_ "github.com/pingcap/tidb/parser/test_driver"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pkg/errors"
	"github.com/pingcap/tidb/parser/format"
	"os"
)

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, errors.Wrap(err, "[PathExists]file stat error")
}

type frmCharsetVisitor struct {
}

func (v *frmCharsetVisitor) Enter(in ast.Node) (ast.Node, bool) {
	switch in.(type) {
	case *test_driver.ValueExpr:
		valueExpr := in.(*test_driver.ValueExpr)
		valueExpr.Type.Charset = ""
	}
	return in, false
}

func (v *frmCharsetVisitor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

func restore(rootNode ast.Node) ([]byte, error) {
	v := &frmCharsetVisitor{}
	rootNode.Accept(v)

	buf := new(bytes.Buffer)
	ctx := format.NewRestoreCtx(format.DefaultRestoreFlags, buf)
	err := rootNode.Restore(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "restore error")
	}
	return buf.Bytes(), nil
}
