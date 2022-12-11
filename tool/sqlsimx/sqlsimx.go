package sqlsimx

import (
	"github.com/qaqcatz/impomysql/connector"
	"io/ioutil"
)

// do not adjust the order!
var SimDMLFuncs = []func(sql string, result *connector.Result, conn *connector.Connector) (string, error) {
	rmBinOp01,
	rmFields,
}

// SqlSimX: more powerful, flexible sql simplification tool.
//
// You can only provide one sql statement in `inputDMLPath`!
//   If you use `dml`, we will try to remove each node in your sql statement, simplify if the result does not change.
//   If you use `ddl`, we will remove unused tables, columns (only consider `CREATE TABLE` and `INSERT INTO VALUES`, may error).
// Then write the simplified sql to `outputPath`.
func SqlSimX(opt string, inputDMLPath string, inputDDLPath string, outputPath string,
	host string, port int, username string, password string, dbname string, specDmlFunc string) {
	// create connector
	conn, err := connector.NewConnector(host, port, username, password, dbname)
	if err != nil {
		panic("[SqlSimX]create connector error: " + err.Error())
	}

	// read and init input ddl
	ddls, err := connector.ExtractSqlFromPath(inputDDLPath)
	if err != nil {
		panic("[SqlSimX]read ddl error: " + err.Error())
	}
	err = conn.InitDBWithDDL(ddls)
	if err != nil {
		panic("[SqlSimX]init ddl error: " + err.Error())
	}

	// read input dml
	data, err := ioutil.ReadFile(inputDMLPath)
	if err != nil {
		panic("[SqlSimX]read input dml error: " + err.Error())
	}
	dml := string(data)

	// first execute the dml
	result := conn.ExecSQL(dml)
	if result.Err != nil {
		panic("[SqlSimX]exec dml error: " + result.Err.Error())
	}

	// simplify
	switch opt {
	case "dml":
		for _, simDMLFunc := range SimDMLFuncs {
			if specDmlFunc != "" && getFunctionName(simDMLFunc) != specDmlFunc {
				continue
			}
			dml, err = simDMLFunc(dml, result, conn)
			if err != nil {
				panic("[SqlSimX]sim dml error: " + err.Error())
			}
		}

		// write to output file
		err = ioutil.WriteFile(outputPath, []byte(dml), 0777)
		if err != nil {
			panic("[SqlSimX]write output dml error: " + err.Error())
		}
	case "ddl":
		newDDLs, err := rmtbcol(dml, ddls)
		if err != nil {
			panic("[SqlSimX]sim ddl error: " + err.Error())
		}

		// write to output file
		err = ioutil.WriteFile(outputPath, []byte(newDDLs), 0777)
		if err != nil {
			panic("[SqlSimX]write output ddl error: " + err.Error())
		}
	default:
		panic("[SqlSimX]please use dml, ddl")
	}
}