package testsqls

import (
	"testing"
)

func TestReadSQLFile(t *testing.T) {
	data, err, filepath := ReadSQLFile(SQLFileTest)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(string(data))
	t.Log(filepath)
}

func TestGetMariaDBConnector(t *testing.T) {
	// mariadb
	err := EnsureOtherDBTEST(MariaDB)
	if err != nil {
		t.Fatal(err.Error())
	}
	conn, err := GetOtherDBConnector(MariaDB)
	if err != nil {
		t.Fatal(err.Error())
	}
	result := conn.ExecSQL("SELECT 1;")
	t.Log(result.ToString())
	// tidb
	err = EnsureOtherDBTEST(TiDB)
	if err != nil {
		t.Fatal(err.Error())
	}
	conn, err = GetOtherDBConnector(TiDB)
	if err != nil {
		t.Fatal(err.Error())
	}
	result = conn.ExecSQL("SELECT 1;")
	t.Log(result.ToString())
	// oceanbase
	err = EnsureOtherDBTEST(OceanBase)
	if err != nil {
		t.Fatal(err.Error())
	}
	conn, err = GetOtherDBConnector(OceanBase)
	if err != nil {
		t.Fatal(err.Error())
	}
	result = conn.ExecSQL("SELECT 1;")
	t.Log(result.ToString())
}