// Package stage2: mutate a sql statement.
//
// 1. visit the sub-AST according to randgen.YYDefault and obtain candidate mutation points.
//
// 2. randomly select a mutation point according to the random seed to mutate.
package stage2
