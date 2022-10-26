// Package stage1: The implication oracle cannot handle these features, remove them.
//
// 1. remove aggregate functions
//
// 2. remove window functions.
//
// 3. remove LEFT|RIGHT JOIN
//
// 4. remove Limit
//
// Note that:
//
// (1) The transformed sql may fail to execute.
//
// (2) we only Support SELECT statement.
//
// How to use: see Init, InitAndExec
package stage1
