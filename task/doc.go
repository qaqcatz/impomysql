// Package task: Basic task for finding logical bug:
//
// 1. random generate sql statements -- randgen.RandGen
//
// 2. initialize random sqls, filter parse error. -- stage1.Init
//
// 3. execute random sqls, filter execute error. -- connector.Connector
//
// 4. for each random sqls, try all of its mutation points, get mutated sqls. -- stage2.MutateAll
//
// 5. execute each mutated sqls, compare their results with the original result,
// detect logical bugs. -- connector.Connector, oracle.Check
//
// How to use: see Run
package task
