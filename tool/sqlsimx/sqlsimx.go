package sqlsimx

import (
	"github.com/qaqcatz/impomysql/connector"
	"io/ioutil"
)

// do not adjust the order!
var SimDMLFuncs = []func(sql string, result *connector.Result, conn *connector.Connector) (string, error) {
	rmBinOpTrue,
	rmBinOpFalse,
	rmFields,
	rmCharset,
}

// SqlSimX: more powerful, flexible sql simplification tool.
// Provide your sql in input file (only support 1 sql),
// we will try to remove each node in original/mutated sql,
// simplify if the result does not change.
// Then write the simplified sql to output file.
// Note that you need to prepare the ddl yourself.
func SqlSimX(inputPath string, outputPath string,
	host string, port int, username string, password string, dbname string) {
	// read input sql
	data, err := ioutil.ReadFile(inputPath)
	if err != nil {
		panic("[SqlSimX]read input sql error: " + err.Error())
	}
	sql := string(data)

	// create connector
	conn, err := connector.NewConnector(host, port, username, password, dbname)
	if err != nil {
		panic("[SqlSimX]create connector error: " + err.Error())
	}

	// first execute the original sql
	res := conn.ExecSQL(sql)
	if res.Err != nil {
		panic("[SqlSimX]exec sql error: " + res.Err.Error())
	}

	// simplify
	for _, simDMLFunc := range SimDMLFuncs {
		sql, err = simDMLFunc(sql, res, conn)
		if err != nil {
			panic("[SqlSimX]sim func error: " + err.Error())
		}
	}

	// write to output file
	err = ioutil.WriteFile(outputPath, []byte(sql), 0777)
	if err != nil {
		panic("[SqlSimX]write output sql error: " + err.Error())
	}
}