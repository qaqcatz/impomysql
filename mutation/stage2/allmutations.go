package stage2

// all mutations
const (
	// *ast.SelectStmt: Distinct true -> false
	FixMDistinctU = "FixMDistinctU"
	// *ast.SelectStmt: Distinct false -> true
	FixMDistinctL = "FixMDistinctL"

	// *ast.SelectStmt: AfterSetOperator UNION -> UNION ALL
	FixMUnionAllU = "FixMUnionAllU"
	// *ast.SelectStmt: AfterSetOperator UNION ALL -> UNION
	FixMUnionAllL = "FixMUnionAllL"

	// *ast.BinaryOperationExpr, *ast.CompareSubqueryExpr: a {>|<|=} b -> a {>=|<=|>=} b
	FixMCmpOpU = "FixMCmpOpU"
	// *ast.BinaryOperationExpr, *ast.CompareSubqueryExpr: a {>=|<=} b -> a {>|<} b
	FixMCmpOpL = "FixMCmpOpL"

	// *ast.PatternInExpr: in(x,x,x) -> in(x,x,x,null)
	FixMInNullU = "FixMInNullU"

	// *ast.SelectStmt: WHERE xxx -> WHERE 1
	FixMWhere1U = "FixMWhere1U"
	// *ast.SelectStmt: WHERE xxx -> WHERE 0
	FixMWhere0L = "FixMWhere0L"

	// *ast.SelectStmt: HAVING xxx -> HAVING 1
	FixMHaving1U = "FixMHaving1U"
	// *ast.SelectStmt: HAVING xxx -> HAVING 0
	FixMHaving0L = "FixMHaving0L"

	// *ast.Join: ON xxx -> ON 1
	FixMOn1U = "FixMOn1U"
	// *ast.Join: ON xxx -> ON 0
	FixMOn0L = "FixMOn0L"

	// *ast.SetOprSelectList: remove Selects[1:] for UNION ALL
	FixMRmUnionAllL = "FixMRmUnionAllL"

	// *ast.PatternLikeExpr: normal char -> '_'|'%',  '_' -> '%'
	RdMLikeU = "RdMLikeU"
	// *ast.PatternLikeExpr: '%' -> '_'
	RdMLikeL = "RdMLikeL"

	// *ast.PatternRegexpExpr: '^'|'$' -> '', normal char -> '.', '+'|'?' -> '*'
	RdMRegExpU = "RdMRegExpU"
	// *ast.PatternRegexpExpr: '*' -> '+'|'?'
	RdMRegExpL = "RdMRegExpL"
)

// discard the following mutation! false positive caused by type conversion:
// 1. --------------------------------------------------
// *ast.BinaryOperationExpr:
//
// a {>|>=} b -> (a) + 1 {>|>=} (b) + 0
//
// a {<|<=} b -> (a) + 0 {<|<=} (b) + 1
// FixMCmpU = "FixMCmpU"
// 2. --------------------------------------------------
// *ast.BinaryOperationExpr:
//
// a {>|>=} b -> (a) + 0 {>|>=} (b) + 1
//
// a {<|<=} b -> (a) + 1 {<|<=} (b) + 0
//
// FixMCmpL = "FixMCmpL"
// 3. --------------------------------------------------
// *ast.CompareSubqueryExpr: ALL true -> false
// FixMCmpSubU = "FixMCmpSubU"
// 4. --------------------------------------------------
// *ast.CompareSubqueryExpr: ALL false -> true
// FixMCmpSubL = "FixMCmpSubL"
// 5. --------------------------------------------------
// *ast.BetweenExpr:
//   expr between l and r
//   ->
//   (expr) >= l and (expr) <= r
//   -> FixMCmpU, 1 and and (expr) <= r, (expr) >= l and 1 )
// RdMBetweenU = "RdMBetweenU"
// 6. --------------------------------------------------
// *ast.BetweenExpr:
//   expr between l and r
//   ->
//   (expr) >= l and (expr) <= r
//   -> FixMCmpOpL / FixMCmpL )
// RdMBetweenL = "RdMBetweenL"
// 7. --------------------------------------------------
// *ast.PatternInExpr: in(x,x,x) -> in(x,x,x,...)
// RdMInU = "RdMInU"
// 8. --------------------------------------------------
// *ast.PatternInExpr: in(x,x,x,...) -> in(x,x,x)
// RdMInL = "RdMInL"