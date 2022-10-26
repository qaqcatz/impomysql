// Package sqlsexecutor: read .sql file(MySQL) file or sqls []string, parse each sql to ast, execute them, get the results.
//
// Note that: Comments in .sql file cannot have ';'.
//
// How to use: see SQLSExecutor, SQLSExecutor.Exec, NewSQLSExecutor, NewSQLSExecutorB, NewSQLSExecutorS
//
// Update: We found that sqlsexecutor is not flexible enough, now it is only used for testing.
package sqlsexecutor
