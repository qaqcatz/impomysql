package stage2

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/format"
	"github.com/pingcap/tidb/parser/opcode"
	_ "github.com/pingcap/tidb/parser/test_driver"
	"math/rand"
)

// all mutations
const (
	// *ast.SelectStmt: Distinct true -> false
	FixMDistinctU = "FixMDistinctU"
	// *ast.SelectStmt: Distinct false -> true
	FixMDistinctL = "FixMDistinctL"
	// *ast.BinaryOperationExpr, *ast.CompareSubqueryExpr: a {>|<|=} b -> a {>=|<=|>=} b
	FixMCmpOpU = "FixMCmpOpU"
	// *ast.BinaryOperationExpr, *ast.CompareSubqueryExpr: a {>=|<=} b -> a {>|<} b
	FixMCmpOpL = "FixMCmpOpL"
	// *ast.BinaryOperationExpr, *ast.CompareSubqueryExpr:
	//
	// a {>|>=} b -> (a) + 1 {>|>=} (b) + 0
	//
	// a {<|<=} b -> (a) + 0 {<|<=} (b) + 1
	//
	// may false positive
	FixMCmpU = "FixMCmpU"
	// *ast.BinaryOperationExpr, *ast.CompareSubqueryExpr:
	//
	// a {>|>=} b -> (a) + 0 {>|>=} (b) + 1
	//
	// a {<|<=} b -> (a) + 1 {<|<=} (b) + 0
	//
	// may false positive
	FixMCmpL = "FixMCmpL"
	// *ast.CompareSubqueryExpr: ALL true -> false
	FixMCmpSubU = "FixMCmpSubU"
	// *ast.CompareSubqueryExpr: ALL false -> true
	FixMCmpSubL = "FixMCmpSubL"
	// *ast.Join: NaturalJoin true -> false
	FixMNaturalJoinU = "FixMNaturalJoinU"
	// *ast.Join: NaturalJoin false -> true
	FixMNaturalJoinL = "FixMNaturalJoinL"
	// *ast.SelectStmt: AfterSetOperator UNION -> UNION ALL
	FixMUnionAllU    = "FixMUnionAllU"
	// *ast.SelectStmt: AfterSetOperator UNION ALL -> UNION
	FixMUnionAllL    = "FixMUnionAllL"
	// *ast.BetweenExpr:
	//   expr between l and r
	//   ->
	//   (expr) >= l and (expr) <= r
	//   -> FixMCmpOpU / FixMCmpU )
	RdMBetweenU = "RdMBetweenU"
	// *ast.BetweenExpr:
	//   expr between l and r
	//   ->
	//   (expr) >= l and (expr) <= r
	//   -> FixMCmpOpL / FixMCmpL )
	RdMBetweenL = "RdMBetweenL"
	// *ast.PatternInExpr: in(x,x,x) -> in(x,x,x,...)
	RdMInU = "RdMInU"
	// *ast.PatternInExpr: in(x,x,x,...) -> in(x,x,x)
	RdMInL = "RdMInL"
	// *ast.PatternLikeExpr: normal char -> '_'|'%',  '_' -> '%'
	RdMLikeU = "RdMLikeU"
	// *ast.PatternLikeExpr: '%' -> '_'
	RdMLikeL = "RdMLikeL"
	// *ast.PatternRegexpExpr: '^'|'$' -> '', normal char -> '.', '+'|'?' -> '*'
	RdMRegExpU = "RdMRegExpU"
	// *ast.PatternRegexpExpr: '*' -> '+'|'?'
	RdMRegExpL = "RdMRegExpL"
	// *ast.SelectStmt: WHERE xxx -> WHERE TRUE | WHERE (xxx) OR 1
	RdMWhereU = "RdMWhereU"
	// *ast.SelectStmt: WHERE xxx -> WHERE FALSE | WHERE (xxx) AND 0
	RdMWhereL = "RdMWhereL"
	// *ast.HavingClause: HAVING xxx -> HAVING TRUE | HAVING (xxx) OR 1
	RdMHavingU = "RdMHavingU"
	// *ast.HavingClause: HAVING xxx -> HAVING FALSE | HAVING (xxx) AND 0
	RdMHavingL = "RdMHavingL"
	// *ast.Join: join value select, same columnNum for NATURAL JOIN
	RdJoinU = "RdJoinU"
	// *ast.Join: remove Right, Tp = 0
	FixJoinL = "RdJoinL"
	// *ast.SetOprStmt: union value select, same columnNum
	RdUnionU = "RdUnionU"
	// *ast.SetOprStmt: remove SelectList.Selects[1:]
	FixUnionL = "RdUnionL"
)

// Candidate: (mutation name, U, candidate node, Flag).
//
// Flag: 1: positive, 0: negative.
//
// U: when positive, 1: upper mutation, 0: lower mutation.
//
// example:
//   [positive]
//     SELECT * FROM T WHERE X > 0;
//     [ upper mutation ] X > 0 -> X >= 0
//     The result set will expand
//   [negative]
//     SELECT * FROM T WHERE (X > 0) IS FALSE;
//     [upper mutation ] x > 0 -> X >= 0
//     The result set will shrink
//   [negative]
//     SELECT * FROM T WHERE (X > 0) IS FALSE;
//     [lower mutation ] x > 0 -> X > 1
//     The result set will expand
// Obviously you should use !(U ^ Flag) to calculate the effect of mutation
type Candidate struct {
	MutationName string // mutation name
	// 1: upper mutation, strings.HasSuffix(MutationName, "U"): true;
	// 0: lower mutation, strings.HasSuffix(MutationName, "L"): true
	U int
	Node ast.Node // candidate node
	Flag int // 1: positive, 0: negative
}

// MutateVisitor: visit the sub-AST according to randgen.YYImpo and obtain the candidate set of mutation points.
//
// Each mutation has its own name.
//
// about the prefix {FixM|RdM}(currently not working):
//
// - FixM means fixed mutation;
//
// - RdM means random mutation;
//
// about the suffix {U|L}:
//
// - U means upper mutation,
//
// - L means lower mutation.
//
// see:
//   FixMDistinctU
//	 FixMDistinctL
//	 FixMCmpOpU
//	 FixMCmpOpL
//	 FixMCmpU
//	 FixMCmpL
//	 FixMCmpSubU
//	 FixMCmpSubL
//	 FixMNaturalJoinU
//	 FixMNaturalJoinL
//	 FixMUnionAllU
//	 FixMUnionAllL
//	 RdMBetweenU
//	 RdMBetweenL
//	 RdMInU
//	 RdMInL
//	 RdMLikeU
//	 RdMLikeL
//	 RdMRegExpU
//	 RdMRegExpL
//	 RdMWhereU
//	 RdMWhereL
//	 RdMHavingU
//	 RdMHavingL
//	 RdJoinU
// 	 FixJoinL
// 	 RdUnionU
//	 FixUnionL
// To obtain the final oracle, you should combine  Candidate.Flag
type MutateVisitor struct {
	CandidatesT map[string][]*Candidate // mutation name : slice of *Candidate
	Candidates   map[ast.Node]int        // mutation name : slice of *Candidate
}

func CalCandidates(rootNode ast.Node) *MutateVisitor {
	v := &MutateVisitor{
		CandidatesT: make(map[string][]*Candidate),
		Candidates:   make(map[ast.Node]int)}
	v.visit(rootNode, 1)
	return v
}

// visit: top
func (v *MutateVisitor) visit(in ast.Node, flag int) {
	switch in.(type) {
	case *ast.SetOprStmt:
		v.visitSetOprStmt(in.(*ast.SetOprStmt), flag)
	case *ast.SelectStmt:
		v.visitSelect(in.(*ast.SelectStmt), flag)
	}
}

// visitSetOprStmt: top1
func (v *MutateVisitor) visitSetOprStmt(in *ast.SetOprStmt, flag int) {
	if in == nil {
		return
	}
	v.visitWithClause(in.SelectList.With, flag)
	v.visitSetOprSelectList(in.SelectList, flag)
}

func (v *MutateVisitor) visitSetOprSelectList(in *ast.SetOprSelectList, flag int) {
	// MySQL only has UNION [ALL]
	if in == nil {
		return
	}
	v.visitWithClause(in.With, flag)
	for _, sel := range in.Selects {
		switch sel.(type) {
		case *ast.SetOprSelectList:
			v.visitSetOprSelectList(sel.(*ast.SetOprSelectList), flag)
		case *ast.SelectStmt:
			v.visitSelect(sel.(*ast.SelectStmt), flag)
		}
	}
}

func (v *MutateVisitor) visitWithClause(in *ast.WithClause, flag int) {
	if in == nil {
		return
	}
	for _, cte := range in.CTEs {
		v.visitSubqueryExpr(cte.Query, flag)
	}
}

func (v *MutateVisitor) visitSubqueryExpr(in *ast.SubqueryExpr, flag int) {
	if in == nil {
		return
	}
	v.visitResultSetNode(in.Query, flag)
}

// visitResultSetNode: include
// SelectStmt, SubqueryExpr, TableSource, TableName, Join and SetOprStmt.
func (v *MutateVisitor) visitResultSetNode(in ast.ResultSetNode, flag int) {
	if in == nil {
		return
	}
	switch in.(type) {
	case *ast.SelectStmt:
		v.visitSelect(in.(*ast.SelectStmt), flag)
	case *ast.SubqueryExpr:
		v.visitSubqueryExpr(in.(*ast.SubqueryExpr), flag)
	case *ast.TableSource:
		v.visitTableSource(in.(*ast.TableSource), flag)
	case *ast.TableName:
		// skim
	case *ast.Join:
		v.visitJoin(in.(*ast.Join), flag)
	case *ast.SetOprStmt:
		v.visitSetOprStmt(in.(*ast.SetOprStmt), flag)
	}
}

func (v *MutateVisitor) visitTableSource(in *ast.TableSource, flag int) {
	if in == nil {
		return
	}
	v.visitResultSetNode(in.Source, flag)
}

// visitSelect: top1
func (v *MutateVisitor) visitSelect(in *ast.SelectStmt, flag int) {
	if in == nil {
		return
	}
	// distinct mutation
	v.Candidates[in] = flag
	// from
	v.visitTableRefClause(in.From, flag)
	// where
	v.visitExprNode(in.Where, flag)
	// having
	v.visitHavingClause(in.Having, flag)
	// with
	v.visitWithClause(in.With, flag)
}

func (v *MutateVisitor) visitTableRefClause(in *ast.TableRefsClause, flag int) {
	if in == nil {
		return
	}
	v.visitJoin(in.TableRefs, flag)
}

func (v *MutateVisitor) visitJoin(in *ast.Join, flag int) {
	if in == nil {
		return
	}
	// skim left | right join
	if in.Tp == ast.LeftJoin || in.Tp == ast.RightJoin {
		return
	}
	v.visitResultSetNode(in.Left, flag)
	v.visitResultSetNode(in.Right, flag)
	// on
	v.visitOnCondition(in.On, flag)
}

func (v *MutateVisitor) visitOnCondition(in *ast.OnCondition, flag int) {
	if in == nil {
		return
	}
	v.visitExprNode(in.Expr, flag)
}

func (v *MutateVisitor) visitHavingClause(in *ast.HavingClause, flag int) {
	if in == nil {
		return
	}
	v.visitExprNode(in.Expr, flag)
}

// visitExprNode: main
func (v *MutateVisitor) visitExprNode(in ast.ExprNode, flag int) {
	if in == nil {
		return
	}
	switch in.(type) {
	//case ast.FuncNode:
	//case ast.ValueExpr:
	case *ast.BetweenExpr:
		v.visitBetweenExpr(in.(*ast.BetweenExpr), flag)
	case *ast.BinaryOperationExpr:
		v.visitBinaryOperationExpr(in.(*ast.BinaryOperationExpr), flag)
	case *ast.CaseExpr:
		// skim
	case *ast.SubqueryExpr:
		v.visitSubqueryExpr(in.(*ast.SubqueryExpr), flag)
	case *ast.CompareSubqueryExpr:
		v.visitCompareSubqueryExpr(in.(*ast.CompareSubqueryExpr), flag)
	case *ast.TableNameExpr:
		// skim
	case *ast.ColumnNameExpr:
		// skim
	case *ast.DefaultExpr:
		// skim
	case *ast.ExistsSubqueryExpr:
		v.visitExistsSubqueryExpr(in.(*ast.ExistsSubqueryExpr), flag)
	case *ast.PatternInExpr:
		v.visitPatternInExpr(in.(*ast.PatternInExpr), flag)
	case *ast.IsNullExpr:
		v.visitIsNullExpr(in.(*ast.IsNullExpr), flag)
	case *ast.IsTruthExpr:
		v.visitIsTruthExpr(in.(*ast.IsTruthExpr), flag)
	case *ast.PatternLikeExpr:
		v.visitPatternLikeExpr(in.(*ast.PatternLikeExpr), flag)
	//case ast.ParamMarkerExpr:
	case *ast.ParenthesesExpr:
		v.visitParenthesesExpr(in.(*ast.ParenthesesExpr), flag)
	case *ast.PositionExpr:
		// skim
	case *ast.PatternRegexpExpr:
		v.visitPatternRegexpExpr(in.(*ast.PatternRegexpExpr), flag)
	case *ast.RowExpr:
		// skim
	case *ast.UnaryOperationExpr:
		v.visitUnaryOperationExpr(in.(*ast.UnaryOperationExpr), flag)
	case *ast.ValuesExpr:
		// skim
	case *ast.VariableExpr:
		// skim
	case *ast.MaxValueExpr:
		// skim
		// https://dev.mysql.com/doc/refman/8.0/en/partitioning-range.html
	case *ast.MatchAgainst:
		// skim, todo
	case *ast.SetCollationExpr:
		// skim
	case *ast.FuncCallExpr:
		v.visitFuncCallExpr(in.(*ast.FuncCallExpr), flag)
	case *ast.FuncCastExpr:
		v.visitFuncCastExpr(in.(*ast.FuncCastExpr), flag)
	case *ast.TrimDirectionExpr:
		v.visitTrimDirectionExpr(in.(*ast.TrimDirectionExpr), flag)
	case *ast.AggregateFuncExpr:
		// skim
	case *ast.WindowFuncExpr:
		// skim
	case *ast.TimeUnitExpr:
		// skim
	case *ast.GetFormatSelectorExpr:
		// skim
	}
}

// visitBetweenExpr: between mutation
func (v *MutateVisitor) visitBetweenExpr(in *ast.BetweenExpr, flag int) {
	if in.Not {
		flag = flag ^ 1
	}
	// between mutation
	v.Candidates[in] = flag
}

// visitBinaryOperationExpr: important bridge, cmp mutation
func (v *MutateVisitor) visitBinaryOperationExpr(in *ast.BinaryOperationExpr, flag int) {
	if in == nil {
		return
	}
	switch in.Op {
	case opcode.LogicAnd:
		v.visitExprNode(in.L, flag)
		v.visitExprNode(in.R, flag)
	case opcode.LeftShift:
		// numeric skim
	case opcode.RightShift:
		// numeric skim
	case opcode.LogicOr:
		v.visitExprNode(in.L, flag)
		v.visitExprNode(in.R, flag)
	case opcode.GE:
		// cmp mutation
		v.Candidates[in] = flag
	case opcode.LE:
		// cmp mutation
		v.Candidates[in] = flag
	case opcode.EQ:
		// cmp mutation
		v.Candidates[in] = flag
	case opcode.NE:
		// cmp mutation
		v.Candidates[in] = flag
	case opcode.LT:
		// cmp mutation
		v.Candidates[in] = flag
	case opcode.GT:
		// cmp mutation
		v.Candidates[in] = flag
	case opcode.Plus:
		// numeric skim
	case opcode.Minus:
		// numeric skim
	case opcode.And:
		// numeric skim
	case opcode.Or:
		// numeric skim
	case opcode.Mod:
		// numeric skim
	case opcode.Xor:
		// numeric skim
	case opcode.Div:
		// numeric skim
	case opcode.Mul:
		// numeric skim
	//case opcode.Not:
	//case opcode.Not2:
	//case opcode.BitNeg:
	case opcode.IntDiv:
		// numeric skim
	case opcode.LogicXor:
		// skim
	case opcode.NullEQ:
		// skim
		//case opcode.In:
		//case opcode.Like:
		//case opcode.Case:
		//case opcode.Regexp:
		//case opcode.IsNull:
		//case opcode.IsTruth:
		//case opcode.IsFalsity:
	}
}

// visitCompareSubqueryExpr: cmp mutation
func (v *MutateVisitor) visitCompareSubqueryExpr(in *ast.CompareSubqueryExpr, flag int) {
	if in == nil {
		return
	}
	// cmp mutation
	v.Candidates[in] = flag
	// in.all false: ANY, in.all true: ALL
	if in.All {
		flag = flag ^ 1
	}
	switch (in.R).(type) {
	case *ast.SubqueryExpr:
		v.visitSubqueryExpr((in.R).(*ast.SubqueryExpr), flag)
	}
}

func (v *MutateVisitor) visitExistsSubqueryExpr(in *ast.ExistsSubqueryExpr, flag int) {
	if in == nil {
		return
	}
	if in.Not {
		flag = flag ^ 1
	}
	switch (in.Sel).(type) {
	case *ast.SubqueryExpr:
		v.visitSubqueryExpr((in.Sel).(*ast.SubqueryExpr), flag)
	}
}

// visitPatternInExpr: in mutation
func (v *MutateVisitor) visitPatternInExpr(in *ast.PatternInExpr, flag int) {
	if in == nil {
		return
	}
	if in.Not {
		flag = flag ^ 1
	}
	// IN (XXX,XXX,XXX) OR IN (SUBQUERY)?
	switch (in.Sel).(type) {
	case *ast.SubqueryExpr:
		v.visitSubqueryExpr((in.Sel).(*ast.SubqueryExpr), flag)
	default:
		// nil
		// in mutation
		v.Candidates[in] = flag
	}
}

// visitIsNullExpr: skim
func (v *MutateVisitor) visitIsNullExpr(in *ast.IsNullExpr, flag int) {
	if in == nil {
		return
	}
}

func (v *MutateVisitor) visitIsTruthExpr(in *ast.IsTruthExpr, flag int) {
	if in == nil {
		return
	}
	if in.Not {
		// IS NOT
		flag = flag ^ 1
	}
	if in.True <= 0 {
		// FALSE
		flag = flag ^ 1
	}
	v.visitExprNode(in.Expr, flag)
}

// visitPatternLikeExpr: like mutation
func (v *MutateVisitor) visitPatternLikeExpr(in *ast.PatternLikeExpr, flag int) {
	if in == nil {
		return
	}
	if in.Not {
		flag = flag ^ 1
	}
	// like mutation
	v.Candidates[in] = flag
}

// visitParenthesesExpr: ()
func (v *MutateVisitor) visitParenthesesExpr(in *ast.ParenthesesExpr, flag int) {
	if in == nil {
		return
	}
	v.visitExprNode(in.Expr, flag)
}

// visitPatternRegexpExpr: regexp mutation
func (v *MutateVisitor) visitPatternRegexpExpr(in *ast.PatternRegexpExpr, flag int) {
	if in == nil {
		return
	}
	if in.Not {
		flag = flag ^ 1
	}
	// regexp mutation
	v.Candidates[in] = flag
}

// visitUnaryOperationExpr: important bridge
func (v *MutateVisitor) visitUnaryOperationExpr(in *ast.UnaryOperationExpr, flag int) {
	if in == nil {
		return
	}
	switch in.Op {
	case opcode.Plus:
		// numeric skim
	case opcode.Minus:
		// numeric skim
	case opcode.Not:
		flag = flag ^ 1
		v.visitExprNode(in.V, flag)
	case opcode.Not2:
		flag = flag ^ 1
		v.visitExprNode(in.V, flag)
	case opcode.BitNeg:
		// numeric skim
	}
}

// visitFuncCallExpr: todo
func (v *MutateVisitor) visitFuncCallExpr(in *ast.FuncCallExpr, flag int) {
	if in == nil {
		return
	}
}

// visitFuncCastExpr: todo
func (v *MutateVisitor) visitFuncCastExpr(in *ast.FuncCastExpr, flag int) {
	if in == nil {
		return
	}
}

// visitTrimDirectionExpr: todo
func (v *MutateVisitor) visitTrimDirectionExpr(in *ast.TrimDirectionExpr, flag int) {
	if in == nil {
		return
	}
}

func (v *MutateVisitor) miningSelectStmt(in *ast.SelectStmt, flag int) {


}

func (v *MutateVisitor) addCandidate(mutationName string, u int, in ast.Node, flag int) {
	var ls []*Candidate = nil
	ok := false
	if ls, ok = v.CandidatesT[mutationName]; !ok {
		ls = make([]*Candidate, 0)
		v.CandidatesT[mutationName] = ls
	}
	ls = append(ls, &Candidate{
		MutationName: mutationName,
		U: u,
		Node: in,
		Flag: flag,
	})
}

// ImpoMutate: you can choose any mutation point to mutate, each mutation has no side effects.
func ImpoMutate(rootNode ast.Node, v *MutateVisitor, seed int64) {
	candidates := v.Candidates
	rand.Seed(seed)
	idx := rand.Intn(len(candidates))
	var candidate ast.Node = nil
	var flag int = 0
	i := 0
	for candidate, flag = range candidates {
		if i == idx {
			break
		}
		i++
	}
	fmt.Println(candidate, flag)
}

// Stage2:
//
// 1. visit the sub-AST according to randgen.YYImpo and obtain the candidate set of  mutation points.
//
// 2. you can choose any mutation point to mutate, each mutation has no side effects.
func Stage2(sql string, seed int64) (string, error) {
	// 1
	p := parser.New()
	stmtNodes, _, err := p.Parse(sql, "", "")
	if err != nil {
		return "", errors.New("Stage2: p.Parse() error: " + err.Error())
	}
	if stmtNodes == nil || len(stmtNodes) == 0 {
		return "", errors.New("Stage1: stmtNodes == nil || len(stmtNodes) == 0 ")
	}
	rootNode := &stmtNodes[0]

	v := CalCandidates(*rootNode)

	// 2
	ImpoMutate(*rootNode, v, seed)
	buf := new(bytes.Buffer)
	ctx := format.NewRestoreCtx(format.DefaultRestoreFlags, buf)
	err = (*rootNode).Restore(ctx)
	if err != nil {
		return "", errors.New("Stage2: (*rootNode).Restore() error: " + err.Error())
	}
	return buf.String(), nil
}
