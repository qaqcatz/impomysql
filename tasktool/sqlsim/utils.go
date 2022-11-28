package sqlsim

import (
	"bytes"
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

func restore(rootNode ast.Node) ([]byte, error) {
	buf := new(bytes.Buffer)
	ctx := format.NewRestoreCtx(format.DefaultRestoreFlags, buf)
	err := rootNode.Restore(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "restore error")
	}
	return buf.Bytes(), nil
}
