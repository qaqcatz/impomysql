// Package stage2: mutate a sql statement.
//
// 1. visit the sub-AST and obtain the candidate set of  mutation points.
//
// 2. you can choose any mutation point to mutate, each mutation has no side effects.
//
// How to use: If you want to choose mutation points yourself, see CalCandidates and ImpoMutate / ImpoMutateAndExec.
// If you want to try all of the mutation points, see MutateAll.
//
// all mutations:
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
package stage2
