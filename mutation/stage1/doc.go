// Package stage1. check and init, we cannot support all features, we will try to remove these unsupported features.
// For features that cannot be removed, we will throw errors.
//
// 1. remove aggregate functions(and GROUP BY).
//
// 2. remove window functions.
//
// 3. remove LEFT|RIGHT JOIN
//
// 4. remove Limit
//
// 5. remove uncertain functions
//
// Note that:
//
// (1) we only support SELECT statement.
//
// (2) make sure your sql has no side-effects, such as assign operations, SELECT into.
//
// (3) The transformed sql may fail to execute.
//
// How to use: see Init, InitAndExec
package stage1
