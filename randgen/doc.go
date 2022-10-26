// Package randgen uses go-randgen(https://github.com/pingcap/go-randgen) to generator random ddl / dml statements.
//
// Note that the tidb/parser version in go-randgen's dependencies is conflict with us,
// so we use shell commands to deal with it. (with the help of runtime.Caller() it works)
//
// How to use: see RandGen, RandGenAndExecDDL
package randgen