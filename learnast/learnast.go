package learnast

import (
	"bytes"
	"github.com/pkg/errors"
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

	PrintNode(in)

	fmt.Println()
	return in, false
}

func (v *learnASTVisitor) Leave(in ast.Node) (ast.Node, bool) {
	v.depth -= 1
	return in, true
}

func PrintNode(in ast.Node) {
	switch in.(type) {
	case *ast.WithClause:
		printWithClause(in.(*ast.WithClause))
	case *ast.SetOprStmt:
		printSetOprStmt(in.(*ast.SetOprStmt))
	case *ast.SetOprSelectList:
		printSetOprSelectList(in.(*ast.SetOprSelectList))
	case *ast.SubqueryExpr:
		printSubqueryExpr(in.(*ast.SubqueryExpr))
	case *ast.SelectStmt:
		printSelectStmt(in.(*ast.SelectStmt))
	case *ast.Limit:
		printLimit(in.(*ast.Limit))
	case *test_driver.ValueExpr:
		printValueExpr(in.(*test_driver.ValueExpr))
	case *ast.BinaryOperationExpr:
		printBinaryOperationExpr(in.(*ast.BinaryOperationExpr))
	case *ast.CompareSubqueryExpr:
		printCompareSubqueryExpr(in.(*ast.CompareSubqueryExpr))
	case *ast.ExistsSubqueryExpr:
		printExistsSubqueryExpr(in.(*ast.ExistsSubqueryExpr))
	case *ast.PatternInExpr:
		printPatternInExpr(in.(*ast.PatternInExpr))
	case *ast.TableRefsClause:
		printTableRefsClause(in.(*ast.TableRefsClause))
	case *ast.Join:
		printJoin(in.(*ast.Join))
	case *ast.TableSource:
		printTableSource(in.(*ast.TableSource))
	case *ast.TableName:
		printTableName(in.(*ast.TableName))
	case *ast.FieldList:
		printFieldList(in.(*ast.FieldList))
	case *ast.SelectField:
		printSelectField(in.(*ast.SelectField))
	case *ast.ColumnNameExpr:
		printColumnNameExpr(in.(*ast.ColumnNameExpr))
	case *ast.ColumnName:
		printColumnName(in.(*ast.ColumnName))
	case *ast.AggregateFuncExpr:
		printAggregateFuncExpr(in.(*ast.AggregateFuncExpr))
	case *ast.WindowFuncExpr:
		printWindowFuncExpr(in.(*ast.WindowFuncExpr))
	case *ast.WindowSpec:
		printWindowSpec(in.(*ast.WindowSpec))
	case *ast.PatternLikeExpr:
		printPatternLikeExpr(in.(*ast.PatternLikeExpr))
	case *ast.PatternRegexpExpr:
		printPatternRegexpExpr(in.(*ast.PatternRegexpExpr))
	default:
	}
}

func printWithClause(in *ast.WithClause) {
	fmt.Print(" [|CTEs|] ", len(in.CTEs), " [IsRecursive] ", in.IsRecursive)
}

func printSetOprStmt(in *ast.SetOprStmt) {
	fmt.Print(" [OrderBy?] ", in.OrderBy != nil)
	fmt.Print(" [Limit?] ", in.Limit != nil)
	fmt.Print(" [With?] ", in.With != nil)
}

func printSetOprSelectList(in *ast.SetOprSelectList) {
	fmt.Print(" [AfterSetOperator] ")
	if in.AfterSetOperator != nil {
		fmt.Print(in.AfterSetOperator.String())
	}
	fmt.Print(" [With?] ", in.With != nil)
	fmt.Print(" [Selects] ", len(in.Selects))
	for i, sel := range in.Selects {
		fmt.Print(" [", i, "] ", reflect.TypeOf(sel))
	}
}

func printSelectStmt(in *ast.SelectStmt) {
	fmt.Print(" [Distinct] ", in.Distinct, " [AfterSetOperator] ", in.AfterSetOperator)
	fmt.Print(" [IsInBraces] ", in.IsInBraces)
}

func printSubqueryExpr(in *ast.SubqueryExpr) {
	fmt.Print(" [Query type] ", reflect.TypeOf(in.Query))
	fmt.Print(" [Evaluated] ", in.Evaluated)
	fmt.Print(" [Correlated] ", in.Correlated)
	fmt.Print(" [MultiRows] ", in.MultiRows)
	fmt.Print(" [Exists] ", in.Exists)
}

func printLimit(in *ast.Limit) {
	fmt.Print(" [Count type] ", reflect.TypeOf(in.Count), " [Offset type] ", reflect.TypeOf(in.Offset))
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

func printValueExpr(in *test_driver.ValueExpr) {
	fmt.Print("[Kind] ", valueKindMap[in.Kind()], " [Value] ", in.GetValue())
}

func printBinaryOperationExpr(in *ast.BinaryOperationExpr) {
	fmt.Print(" [Op] ", in.Op)
}

func printCompareSubqueryExpr(in *ast.CompareSubqueryExpr) {
	fmt.Print(" [Op] ", in.Op, " [All] ", in.All)
}

func printExistsSubqueryExpr(in *ast.ExistsSubqueryExpr) {
	fmt.Print(" [Not] ", in.Not)
}

func printPatternInExpr(in *ast.PatternInExpr) {
	fmt.Print(" [|List|] ", len(in.List), " [Not] ", in.Not, " [Sel type] ", reflect.TypeOf(in.Sel))
}

func printTableRefsClause(in *ast.TableRefsClause) {
	fmt.Print(" [TableRefs?] ", in.TableRefs != nil)
}

var joinTpMap = map[int]string {
	0: "none",
	1: "CrossJoin",
	2: "LeftJoin",
	3: "RightJoin",
}

func printJoin(in *ast.Join) {
	fmt.Print(" [Tp] ", joinTpMap[int(in.Tp)], " [NaturalJoin] ", in.NaturalJoin, " [StraightJoin] ", in.StraightJoin, " [ExplicitParens] ", in.ExplicitParens)
}

func printTableSource(in *ast.TableSource) {
	fmt.Print(" [AsName] ", in.AsName)
}

func printTableName(in *ast.TableName) {
	fmt.Print(" [Schema] ", in.Schema, " [Name] ", in.Name,
		" [DBInfo] ", in.DBInfo, " [TableInfo] ", in.TableInfo, " [IndexHints] ", in.IndexHints,
		" [PartitionNames] ", in.PartitionNames, " [TableSample] ", in.TableSample, " [AsOf] ", in.AsOf)
}

func printFieldList(in *ast.FieldList) {
	fmt.Print(" [|Fields|] ", len(in.Fields))
}

func printSelectField(in *ast.SelectField) {
	fmt.Print(" [AsName] ", in.AsName, " [WildCard] ", in.WildCard)
}

func printColumnNameExpr(in *ast.ColumnNameExpr) {
	fmt.Print(" [Refer] ", in.Refer)
}

func printColumnName(in *ast.ColumnName) {
	fmt.Print(" [Schema] ", in.Schema, " [Table] ", in.Table, " [Name] ", in.Name)
}

func printAggregateFuncExpr(in *ast.AggregateFuncExpr) {
	fmt.Print(" [F] ", in.F, " [|Args|] ", len(in.Args), " [Distinct] ", in.Distinct, " [Order?] ", in.Order != nil)
}

func printWindowFuncExpr(in *ast.WindowFuncExpr) {
	fmt.Print(" [F] ", in.F, " [|Args|] ", len(in.Args), " [Distinct] ", in.Distinct,
		" [IgnoreNull] ", in.IgnoreNull, " [FromLast] ", in.FromLast)
}

func printWindowSpec(in *ast.WindowSpec) {
	fmt.Print(" [Name] ", in.Name, " [Ref] ", in.Ref, " [PartitionBy?] ", in.PartitionBy != nil,
		" [OrderBy?] ", in.OrderBy != nil, " [Frame?] ", in.Frame != nil, " [OnlyAlias] ", in.OnlyAlias)
}

func printPatternLikeExpr(in *ast.PatternLikeExpr) {
	if t, ok := (in.Expr).(*test_driver.ValueExpr); ok {
		fmt.Print(" [Expr] ")
		printValueExpr(t)
	}
	if t, ok := (in.Pattern).(*test_driver.ValueExpr); ok {
		fmt.Print(" [Pattern] ")
		printValueExpr(t)
	}
	fmt.Print(" [Not] ", in.Not, " [Escape] ", string(in.Escape))
}

func printPatternRegexpExpr(in *ast.PatternRegexpExpr) {
	if t, ok := (in.Expr).(*test_driver.ValueExpr); ok {
		fmt.Print(" [Expr] ")
		printValueExpr(t)
	}
	if t, ok := (in.Pattern).(*test_driver.ValueExpr); ok {
		fmt.Print(" [Pattern] ")
		printValueExpr(t)
	}
	fmt.Print(" [Not] ", in.Not, " [Re] ", in.Re, " [Sexpr] ", in.Sexpr)
}

func learnAST(sql string) (string, error) {
	p := parser.New()
	stmtNodes, _, err := p.Parse(sql, "", "")
	if err != nil {
		return "", errors.Wrap(err, "[learnAST]parse error")
	}
	if stmtNodes == nil || len(stmtNodes) == 0 {
		return "", errors.New("[learnAST]stmtNodes == nil || len(stmtNodes) == 0")
	}
	rootNode := &stmtNodes[0]
	v := &learnASTVisitor{depth: 0}
	(*rootNode).Accept(v)
	buf := new(bytes.Buffer)
	ctx := format.NewRestoreCtx(format.DefaultRestoreFlags, buf)
	err = (*rootNode).Restore(ctx)
	if err != nil {
		return "", errors.Wrap(err, "[learnAST]restore error")
	}
	return buf.String(), nil
}