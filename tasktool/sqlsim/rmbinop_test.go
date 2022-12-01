package sqlsim

import (
	"github.com/qaqcatz/impomysql/testsqls"
	"testing"
)

func TestRmBinOp(t *testing.T) {
	rmBinOp(testsqls.SQLBinaryOp)
}