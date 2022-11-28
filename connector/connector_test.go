package connector

import (
	"strconv"
	"testing"
)

// sudo docker run -itd --name test -p 13306:3306 -e MYSQL_ROOT_PASSWORD=123456 mysql
const (
	testHost = "127.0.0.1"
	testPort = 13306
	testUsername = "root"
	testPassword = "123456"
	testDBname = "TEST"
)

func TestConnector_ExecSQL(t *testing.T) {
	conn, err := NewConnector(testHost, testPort, testUsername, testPassword, testDBname)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	// create table
	result := conn.ExecSQL("DROP TABLE IF EXISTS T")
	if result.Err != nil {
		t.Fatalf("%+v", result.Err)
	} else {
		t.Log(result.ToString())
	}
	result = conn.ExecSQL("CREATE TABLE T(ID INT, NAME TEXT, X DOUBLE)")
	if result.Err != nil {
		t.Fatalf("%+v", result.Err)
	} else {
		t.Log(result.ToString())
	}

	for i := 0; i < 3; i++ {
		result := conn.ExecSQL("INSERT INTO T VALUES ("+strconv.Itoa(i)+", '"+string(rune(i+'A'))+"', -" + strconv.Itoa(i) + ")")
		if result.Err != nil {
			t.Fatalf("%+v", result.Err)
		} else {
			t.Log(result.ToString())
		}
	}

	// normal
	result = conn.ExecSQL("SELECT 1+2, ID, NAME, X FROM T;")
	if result.Err != nil {
		t.Fatalf("%+v", result.Err)
	} else {
		t.Log(result.ToString())
	}

	// error
	result = conn.ExecSQL("select 9223372036854775807 + 1")
	if result.Err != nil {
		t.Logf("%+v", result.Err)
	} else {
		t.Fatal("must error!")
	}

	testSql := "SELECT (~DEGREES(0.9219647951826007)|FORMAT_BYTES(`f1`)), (~1^`f1`) FROM (SELECT (X^_UTF8MB4'do'-X) AS `f1` FROM (SELECT X FROM T) AS `t1`) AS `t2`;"

	result = conn.ExecSQL(testSql)
	if result.Err != nil {
		t.Logf("%+v", result.Err)
	} else {
		t.Fatal("must error!")
	}

	errCode, err := result.GetErrorCode()
	if err == nil {
		t.Log("error code = " ,errCode)
	} else {
		t.Fatalf("%+v", err)
	}

	// result cmp
	result1 := conn.ExecSQL("SELECT ID, NAME, X FROM T;")
	if result1.Err != nil {
		t.Fatalf("%+v", result1.Err)
	}
	result2 := conn.ExecSQL("SELECT ID, NAME, X FROM T WHERE ID != 1;")
	if result2.Err != nil {
		t.Fatalf("%+v", result2.Err)
	}
	result3 := conn.ExecSQL("SELECT ID, NAME, X FROM T WHERE ID != 2;")
	if result3.Err != nil {
		t.Fatalf("%+v", result3.Err)
	}

	cmp, err := result1.CMP(result1)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	if cmp != 0 {
		t.Fatal("must 0")
	}

	cmp, err = result1.CMP(result2)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	if cmp != 1 {
		t.Fatal("must 1")
	}

	cmp, err = result2.CMP(result1)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	if cmp != -1 {
		t.Fatal("must -1")
	}

	cmp, err = result3.CMP(result2)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	if cmp != 2 {
		t.Fatal("must 2")
	}
}