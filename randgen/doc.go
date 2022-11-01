// Package randgen uses go-randgen(https://github.com/pingcap/go-randgen) to generator random ddl / dml statements.
//
// Note that the tidb/parser version in go-randgen's dependencies is conflict with us,
// so we use shell commands to deal with it, see https://github.com/qaqcatz/gorandgensh.
// You should compile it first, and tell us the path of the executable file.
// default path: under this package + gorandgensh. (with the help of runtime.Caller())
//
// How to use: see RandGen, RandGenAndExecDDL
package randgen