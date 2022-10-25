package stage2

import (
	"bytes"
	"errors"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/format"
	"github.com/pingcap/tidb/parser/opcode"
	_ "github.com/pingcap/tidb/parser/test_driver"
	"log"
	"strings"
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
	// *ast.BinaryOperationExpr:
	//
	// a {>|>=} b -> (a) + 1 {>|>=} (b) + 0
	//
	// a {<|<=} b -> (a) + 0 {<|<=} (b) + 1
	//
	// may false positive, skim
	FixMCmpU = "FixMCmpU"
	// *ast.BinaryOperationExpr:
	//
	// a {>|>=} b -> (a) + 0 {>|>=} (b) + 1
	//
	// a {<|<=} b -> (a) + 1 {<|<=} (b) + 0
	//
	// may false positive, skim
	FixMCmpL = "FixMCmpL"
	// *ast.CompareSubqueryExpr: ALL true -> false
	FixMCmpSubU = "FixMCmpSubU"
	// *ast.CompareSubqueryExpr: ALL false -> true
	FixMCmpSubL = "FixMCmpSubL"
	// *ast.SelectStmt: AfterSetOperator UNION -> UNION ALL
	FixMUnionAllU = "FixMUnionAllU"
	// *ast.SelectStmt: AfterSetOperator UNION ALL -> UNION
	FixMUnionAllL = "FixMUnionAllL"
	// *ast.BetweenExpr:
	//   expr between l and r
	//   ->
	//   (expr) >= l and (expr) <= r
	//   -> FixMCmpU, 1 and and (expr) <= r, (expr) >= l and 1 )
	// may false positive, skim
	RdMBetweenU = "RdMBetweenU"
	// *ast.BetweenExpr:
	//   expr between l and r
	//   ->
	//   (expr) >= l and (expr) <= r
	//   -> FixMCmpOpL / FixMCmpL )
	// may false positive, skim
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
	// *ast.SelectStmt: HAVING xxx -> HAVING TRUE | HAVING (xxx) OR 1
	RdMHavingU = "RdMHavingU"
	// *ast.SelectStmt: HAVING xxx -> HAVING FALSE | HAVING (xxx) AND 0
	RdMHavingL = "RdMHavingL"
	// *ast.Join: ON xxx -> ON TRUE | ON (xxx) OR 1
	RdMOnU = "RdMOnU"
	// *ast.Join: ON xxx -> ON FALSE | ON (xxx) AND 0
	RdMOnL = "RdMOnL"
	// *ast.SetOprSelectList: remove Selects[1:]
	FixMUnionL = "RdUnionL"
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
	U    int
	Node ast.Node // candidate node
	Flag int      // 1: positive, 0: negative
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
//	 RdMOnU
//	 RdMOnL
//	 FixMUnionL
// function:
//   visitxxx: calculate flag, call miningxxx
//   miningxxx: call addxxx
//   addxxx: calculate mutation u/l
type MutateVisitor struct {
	Candidates map[string][]*Candidate // mutation name : slice of *Candidate
}

// CalCandidates: visit the sub-AST according to randgen.YYImpo and obtain the candidate set of mutation points.
func CalCandidates(rootNode ast.Node) *MutateVisitor {
	v := &MutateVisitor{
		Candidates: make(map[string][]*Candidate)}
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

// visitSetOprSelectList: miningSetOprSelectList
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

	v.miningSetOprSelectList(in, flag)
}

func (v *MutateVisitor) visitWithClause(in *ast.WithClause, flag int) {
	if in == nil {
		return
	}
	// cannot support recursive WITH
	if in.IsRecursive {
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

// visitSelect: top1, miningSelectStmt
func (v *MutateVisitor) visitSelect(in *ast.SelectStmt, flag int) {
	if in == nil {
		return
	}

	// from
	v.visitTableRefClause(in.From, flag)
	// where
	v.visitExprNode(in.Where, flag)
	// having
	v.visitHavingClause(in.Having, flag)
	// with
	v.visitWithClause(in.With, flag)

	v.miningSelectStmt(in, flag)
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

	v.miningJoin(in, flag)
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

// visitBetweenExpr: miningBetweenExpr
func (v *MutateVisitor) visitBetweenExpr(in *ast.BetweenExpr, flag int) {
	if in == nil {
		return
	}
	if in.Not {
		flag = flag ^ 1
	}
	// after if in.Not
	v.miningBetweenExpr(in, flag)
}

// visitBinaryOperationExpr: important bridge, miningBinaryOperationExpr
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
		// cmp mutation, see miningBinaryOperationExpr
	case opcode.LE:
		// cmp mutation, see miningBinaryOperationExpr
	case opcode.EQ:
		// cmp mutation, see miningBinaryOperationExpr
	case opcode.NE:
		// cmp mutation, see miningBinaryOperationExpr
	case opcode.LT:
		// cmp mutation, see miningBinaryOperationExpr
	case opcode.GT:
		// cmp mutation, see miningBinaryOperationExpr
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

	v.miningBinaryOperationExpr(in, flag)
}

// visitCompareSubqueryExpr: miningCompareSubqueryExpr
func (v *MutateVisitor) visitCompareSubqueryExpr(in *ast.CompareSubqueryExpr, flag int) {
	if in == nil {
		return
	}
	// before All
	v.miningCompareSubqueryExpr(in, flag)
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

// visitPatternInExpr: miningPatternInExpr
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
		// after in.Not
		v.miningPatternInExpr(in, flag)
	}
}

func (v *MutateVisitor) visitIsNullExpr(in *ast.IsNullExpr, flag int) {
	if in == nil {
		return
	}
	// skim
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

// visitPatternLikeExpr: miningPatternLikeExpr
func (v *MutateVisitor) visitPatternLikeExpr(in *ast.PatternLikeExpr, flag int) {
	if in == nil {
		return
	}
	if in.Not {
		flag = flag ^ 1
	}
	// after in.Not
	v.miningPatternLikeExpr(in, flag)
}

// visitParenthesesExpr: ()
func (v *MutateVisitor) visitParenthesesExpr(in *ast.ParenthesesExpr, flag int) {
	if in == nil {
		return
	}
	v.visitExprNode(in.Expr, flag)
}

// visitPatternRegexpExpr: miningPatternRegexpExpr
func (v *MutateVisitor) visitPatternRegexpExpr(in *ast.PatternRegexpExpr, flag int) {
	if in == nil {
		return
	}
	if in.Not {
		flag = flag ^ 1
	}
	// after in.Not
	v.miningPatternRegexpExpr(in, flag)
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

func (v *MutateVisitor) visitFuncCallExpr(in *ast.FuncCallExpr, flag int) {
	if in == nil {
		return
	}
	// skim func call
}

func (v *MutateVisitor) visitFuncCastExpr(in *ast.FuncCastExpr, flag int) {
	if in == nil {
		return
	}
	// skim cast
}

func (v *MutateVisitor) visitTrimDirectionExpr(in *ast.TrimDirectionExpr, flag int) {
	if in == nil {
		return
	}
	// skim trim
}

func (v *MutateVisitor) miningSetOprSelectList(in *ast.SetOprSelectList, flag int) {
	// FixMUnionL
	v.addFixMUnionL(in, flag)
}

func (v *MutateVisitor) miningSelectStmt(in *ast.SelectStmt, flag int) {
	// FixMDistinctU
	v.addFixMDistinctU(in, flag)
	// FixMDistinctL
	v.addFixMDistinctL(in, flag)
	// FixMUnionAllU
	v.addFixMUnionAllU(in, flag)
	// FixMUnionAllL
	v.addFixMUnionAllL(in, flag)
	// RdMWhereU
	v.addRdMWhereU(in, flag)
	// RdMWhereL
	v.addRdMWhereL(in, flag)
    // RdMHavingU
	v.addRdMHavingU(in, flag)
	// RdMHavingL
	v.addRdMHavingL(in, flag)
}

func (v *MutateVisitor) miningJoin(in *ast.Join, flag int) {
	// RdMOnU
	v.addRdMOnU(in, flag)
	// RdMOnL
	v.addRdMOnL(in, flag)
}

func (v *MutateVisitor) miningBinaryOperationExpr(in *ast.BinaryOperationExpr, flag int) {
	// FixMCmpOpU
	v.addFixMCmpOpU(in, flag)
	// FixMCmpOpL
	v.addFixMCmpOpL(in, flag)
	// FixMCmpU
	//v.addFixMCmpU(in, flag)
	// FixMCmpL
	//v.addFixMCmpL(in, flag)
}

func (v *MutateVisitor) miningCompareSubqueryExpr(in *ast.CompareSubqueryExpr, flag int) {
	// FixMCmpOpU
	v.addFixMCmpOpU(in, flag)
	// FixMCmpOpL
	v.addFixMCmpOpL(in, flag)
	// FixMCmpSubU
	v.addFixMCmpSubU(in, flag)
	// FixMCmpSubL
	v.addFixMCmpSubL(in, flag)
}

func (v *MutateVisitor) miningBetweenExpr(in *ast.BetweenExpr, flag int) {
	// RdMBetweenU
	//v.addRdMBetweenU(in, flag)
	// RdMBetweenL
	//v.addRdMBetweenL(in, flag)
}

func (v *MutateVisitor) miningPatternInExpr(in *ast.PatternInExpr, flag int) {
	// RdMInU
	v.addRdMInU(in, flag)
	// RdMInL
	v.addRdMInL(in, flag)
}

func (v *MutateVisitor) miningPatternLikeExpr(in *ast.PatternLikeExpr, flag int) {
	// RdMLikeU
	v.addRdMLikeU(in, flag)
	// RdMLikeL
	v.addRdMLikeL(in, flag)
}

func (v *MutateVisitor) miningPatternRegexpExpr(in *ast.PatternRegexpExpr, flag int) {
	// RdMRegExpU
	v.addRdMRegExpU(in, flag)
	// RdMRegExpL
	v.addRdMRegExpL(in, flag)
}

func (v *MutateVisitor) addCandidate(mutationName string, u int, in ast.Node, flag int) {
	if strings.HasSuffix(mutationName, "U") && u == 0 {
		log.Fatal("strings.HasSuffix(mutationName, \"U\") && u == 0")
	}
	if strings.HasSuffix(mutationName, "L") && u != 0 {
		log.Fatal("strings.HasSuffix(mutationName, \"L\") && u != 0")
	}
	var ls []*Candidate = nil
	ok := false
	if ls, ok = v.Candidates[mutationName]; !ok {
		ls = make([]*Candidate, 0)
	}
	ls = append(ls, &Candidate{
		MutationName: mutationName,
		U:            u,
		Node:         in,
		Flag:         flag,
	})
	v.Candidates[mutationName] = ls
}

// ImpoMutate: you can choose any candidate to mutate, each mutation has no side effects.
func ImpoMutate(rootNode ast.Node, candidate *Candidate, seed int64) ([]byte, error) {
	var sql []byte = nil
	var err error = nil
	switch candidate.MutationName {
	case FixMDistinctU:
		sql, err = doFixMDistinctU(rootNode, candidate.Node)
	case FixMDistinctL:
		sql, err = doFixMDistinctL(rootNode, candidate.Node)
	case FixMCmpOpU:
		sql, err = doFixMCmpOpU(rootNode, candidate.Node)
	case FixMCmpOpL:
		sql, err = doFixMCmpOpL(rootNode, candidate.Node)
	//case FixMCmpU:
	//	sql, err = doFixMCmpU(rootNode, candidate.Node)
	//case FixMCmpL:
	//	sql, err = doFixMCmpL(rootNode, candidate.Node)
	case FixMCmpSubU:
		sql, err = doFixMCmpSubU(rootNode, candidate.Node)
	case FixMCmpSubL:
		sql, err = doFixMCmpSubL(rootNode, candidate.Node)
	case FixMUnionAllU:
		sql, err = doFixMUnionAllU(rootNode, candidate.Node)
	case FixMUnionAllL:
		sql, err = doFixMUnionAllL(rootNode, candidate.Node)
	//case RdMBetweenU:
	//	sql, err = doRdMBetweenU(rootNode, candidate.Node, seed)
	//case RdMBetweenL:
	//	sql, err = doRdMBetweenL(rootNode, candidate.Node, seed)
	case RdMInU:
		sql, err = doRdMInU(rootNode, candidate.Node, seed)
	case RdMInL:
		sql, err = doRdMInL(rootNode, candidate.Node, seed)
	case RdMLikeU:
		sql, err = doRdMLikeU(rootNode, candidate.Node, seed)
	case RdMLikeL:
		sql, err = doRdMLikeL(rootNode, candidate.Node, seed)
	case RdMRegExpU:
		sql, err = doRdMRegExpU(rootNode, candidate.Node, seed)
	case RdMRegExpL:
		sql, err = doRdMRegExpL(rootNode, candidate.Node, seed)
	case RdMWhereU:
		sql, err = doRdMWhereU(rootNode, candidate.Node, seed)
	case RdMWhereL:
		sql, err = doRdMWhereL(rootNode, candidate.Node, seed)
	case RdMHavingU:
		sql, err = doRdMHavingU(rootNode, candidate.Node, seed)
	case RdMHavingL:
		sql, err = doRdMHavingL(rootNode, candidate.Node, seed)
	case RdMOnU:
		sql, err = doRdMOnU(rootNode, candidate.Node, seed)
	case RdMOnL:
		sql, err = doRdMOnL(rootNode, candidate.Node, seed)
	case FixMUnionL:
		sql, err = doFixMUnionL(rootNode, candidate.Node)
	}
	if err != nil {
		return nil, errors.New("ImpoMutate: " +  err.Error())
	}
	return sql, nil
}

func restore(rootNode ast.Node) ([]byte, error) {
	buf := new(bytes.Buffer)
	ctx := format.NewRestoreCtx(format.DefaultRestoreFlags, buf)
	err := rootNode.Restore(ctx)
	if err != nil {
		return nil, errors.New("restore error: " + err.Error())
	}
	return buf.Bytes(), nil
}
