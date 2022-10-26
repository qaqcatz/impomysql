// Package mutation: mutate a sql statement, a mutated sql will be created through the following stages:
//
// stage1. init, some of the sql features we can not handle, we need to convert them into valid format.
//
// stage2. mutation.
// Note that a sql statement can have multiple mutations.
//
// You should execute these mutated sqls yourself, and use the implication oracle
// to detect logical bugs. see oracle.Check
package mutation