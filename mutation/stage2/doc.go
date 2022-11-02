// Package stage2: mutate a sql statement.
//
// 1. visit the sub-AST and obtain the candidate set of  mutation points.
//
// 2. you can choose any mutation point to mutate, each mutation has no side effects.
//
// How to use: If you want to choose mutation points yourself, see CalCandidates and ImpoMutate / ImpoMutateAndExec.
// If you want to try all of the mutation points, see MutateAll / MutateAllAndExec.
//
// all mutations:
//   FixMDistinctU
//	 FixMDistinctL
//	 FixMCmpOpU
//	 FixMCmpOpL
//	 FixMUnionAllU
//	 FixMUnionAllL
//   FixMInNullU
//	 FixMWhere1U
//	 FixMWhere0L
//	 FixMHaving1U
//	 FixMHaving0L
//	 FixMOn1U
//	 FixMOn0L
//	 FixMRmUnionAllL
//	 RdMLikeU
//	 RdMLikeL
//	 RdMRegExpU
//	 RdMRegExpL
package stage2
