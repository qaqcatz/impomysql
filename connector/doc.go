// Package connector: connect to mysql, auto create database if not exists,
// execute raw sql statements, return raw execution result or error.
//
// We also provide some useful functions:
//   - init database, see Connector.InitDB
//   - extract sqls from a sql file, see ExtractSQL
//   - init database with a sql file, see Connector.InitDBWithDDL, Connector.InitDBWithDDLPath
//   - connectorPool, see NewConnectorPool, ConnectorPool.WaitForFree, ConnectorPool.BackToPool
package connector
