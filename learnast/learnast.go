package learnast

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/format"
	"github.com/pingcap/tidb/parser/test_driver"
	_ "github.com/pingcap/tidb/parser/test_driver"
	"reflect"
)

// only for learning, ignore it
type learnASTVisitor struct {
	depth int
}

func (v *learnASTVisitor) Enter(in ast.Node) (ast.Node, bool) {
	prefix := ""
	for i := 0; i < v.depth; i++ {
		prefix += "| "
	}
	fmt.Printf("\x1b[%dm%s\x1b[0m", 33, prefix)
	v.depth += 1

	fmt.Printf("\x1b[%dm%T\x1b[0m", 31, in)
	switch in.(type) {
	case *test_driver.ValueExpr:
		v.EnterValueExpr(in.(*test_driver.ValueExpr))
	case *ast.BinaryOperationExpr:
		v.EnterBinaryOperationExpr(in.(*ast.BinaryOperationExpr))
	case *ast.CompareSubqueryExpr:
		v.EnterCompareSubqueryExpr(in.(*ast.CompareSubqueryExpr))
	case *ast.PatternInExpr:
		v.EnterPatternInExpr(in.(*ast.PatternInExpr))
	case *ast.TableSource:
		v.EnterTableSource(in.(*ast.TableSource))
	case *ast.TableName:
		v.EnterTableName(in.(*ast.TableName))
	case *ast.SelectField:
		v.EnterSelectField(in.(*ast.SelectField))
	case *ast.ColumnNameExpr:
		v.EnterColumnNameExpr(in.(*ast.ColumnNameExpr))
	case *ast.ColumnName:
		v.EnterColumnName(in.(*ast.ColumnName))
	case *ast.AggregateFuncExpr:
		v.EnterAggregateFuncExpr(in.(*ast.AggregateFuncExpr))
	case *ast.WindowFuncExpr:
		v.EnterWindowFuncExpr(in.(*ast.WindowFuncExpr))
	case *ast.WindowSpec:
		v.EnterWindowSpec(in.(*ast.WindowSpec))
	default:
	}
	fmt.Println()
	return in, false
}

var valueKindMap = map[byte]string {
	0: "KindNull",
	1: "KindInt64",
	2: "KindUint64",
	3: "KindFloat32",
	4: "KindFloat64",
	5: "KindString",
	6: "KindBytes",
	7: "KindBinaryLiteral",
	8: "KindMysqlDecimal",
	9: "KindMysqlDuration",
	10: "KindMysqlEnum",
	11: "KindMysqlBit",
	12: "KindMysqlSet",
	13: "KindMysqlTime",
	14: "KindInterface",
	15: "KindMinNotNull",
	16: "KindMaxValue",
	17: "KindRaw",
	18: "KindMysqlJSON",
}

func (v *learnASTVisitor) EnterValueExpr(in *test_driver.ValueExpr) {
	fmt.Print("[Kind] ", valueKindMap[in.Kind()], " [Value] ", in.GetValue())
}

func (v *learnASTVisitor) EnterBinaryOperationExpr(in *ast.BinaryOperationExpr) {
	fmt.Print(" [Op] ", in.Op)
}

func (v *learnASTVisitor) EnterCompareSubqueryExpr(in *ast.CompareSubqueryExpr) {
	fmt.Print(" [Op] ", in.Op, " [All] ", in.All)
}

func (v *learnASTVisitor) EnterPatternInExpr(in *ast.PatternInExpr) {
	fmt.Print(" [|List|] ", len(in.List), " [Not] ", in.Not, " [Sel type] ", reflect.TypeOf(in.Sel))
}

func (v *learnASTVisitor) EnterTableSource(in *ast.TableSource) {
	fmt.Print(" [AsName] ", in.AsName)
}

func (v *learnASTVisitor) EnterTableName(in *ast.TableName) {
	fmt.Print(" [Schema] ", in.Schema, " [Name] ", in.Name,
		" [DBInfo] ", in.DBInfo, " [TableInfo] ", in.TableInfo, " [IndexHints] ", in.IndexHints,
		" [PartitionNames] ", in.PartitionNames, " [TableSample] ", in.TableSample, " [AsOf] ", in.AsOf)
}

func (v *learnASTVisitor) EnterSelectField(in *ast.SelectField) {
	fmt.Print(" [AsName] ", in.AsName, " [WildCard] ", in.WildCard)
}

func (v *learnASTVisitor) EnterColumnNameExpr(in *ast.ColumnNameExpr) {
	fmt.Print(" [Refer] ", in.Refer)
}

func (v *learnASTVisitor) EnterColumnName(in *ast.ColumnName) {
	fmt.Print(" [Schema] ", in.Schema, " [Table] ", in.Table, " [Name] ", in.Name)
}

func (v *learnASTVisitor) EnterAggregateFuncExpr(in *ast.AggregateFuncExpr) {
	fmt.Print(" [F] ", in.F, " [|Args|] ", len(in.Args), " [Distinct] ", in.Distinct, " [Order?] ", in.Order != nil)
}

func (v *learnASTVisitor) EnterWindowFuncExpr(in *ast.WindowFuncExpr) {
	fmt.Print(" [F] ", in.F, " [|Args|] ", len(in.Args), " [Distinct] ", in.Distinct,
		" [IgnoreNull] ", in.IgnoreNull, " [FromLast] ", in.FromLast)
}

func (v *learnASTVisitor) EnterWindowSpec(in *ast.WindowSpec) {
	fmt.Print(" [Name] ", in.Name, " [Ref] ", in.Ref, " [PartitionBy?] ", in.PartitionBy != nil,
		" [OrderBy?] ", in.OrderBy != nil, " [Frame?] ", in.Frame != nil, " [OnlyAlias] ", in.OnlyAlias)
}

func (v *learnASTVisitor) Leave(in ast.Node) (ast.Node, bool) {
	v.depth -= 1
	return in, true
}

func learnAST(sql string) (string, error) {
	p := parser.New()
	stmtNodes, _, err := p.Parse(sql, "", "")
	if err != nil {
		return "", err
	}
	if stmtNodes == nil || len(stmtNodes) == 0 {
		return "", errors.New("stmtNodes == nil || len(stmtNodes) == 0")
	}
	rootNode := &stmtNodes[0]
	v := &learnASTVisitor{depth: 0}
	(*rootNode).Accept(v)
	buf := new(bytes.Buffer)
	ctx := format.NewRestoreCtx(format.DefaultRestoreFlags, buf)
	err = (*rootNode).Restore(ctx)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}