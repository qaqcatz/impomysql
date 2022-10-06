package testsqls

import "testing"

func TestReadSQLFile(t *testing.T) {
	data, err, filepath := ReadSQLFile(SQLFileTest)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(string(data))
	t.Log(filepath)
}
