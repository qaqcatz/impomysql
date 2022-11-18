// Package mutation: mutate a sql statement, a mutated sql will be created through the following stages:
//
// stage1. check and init, we cannot support all features, we will try to remove these unsupported features.
// For features that cannot be removed, we will throw errors.
//
// stage2. mutation. Note that a sql statement can have multiple mutations.
//
// You can use the implication oracle to detect logical bugs.
package mutation