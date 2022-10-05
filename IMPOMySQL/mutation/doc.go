// Package mutation: mutate a sql statement, including 3 stages:
//
// stage1. init, some of the sql features we can not handle, we need to convert them into valid format.
//
// stage2. mutation.
//
// stage3. execute the mutated sql.
//
// Note that a sql statement can have multiple mutations.
package mutation
