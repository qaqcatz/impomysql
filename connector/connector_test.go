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
	conn, err := NewConnector(testHost, testPort, testUsername, testPassword, "", "")
	if err != nil {
		t.Fatal(err.Error())
	}
	result := conn.ExecSQL("CREATE DATABASE IF NOT EXISTS TEST")
	if result.Err != nil {
		t.Fatal(result.Err.Error())
	} else {
		t.Log(result.ToString())
	}

	conn, err = NewConnector(testHost, testPort, testUsername, testPassword, testDBname, "")
	if err != nil {
		t.Fatal(err.Error())
	}

	t.Log(conn.ToString())

	result = conn.ExecSQL("DROP TABLE IF EXISTS T")
	if result.Err != nil {
		t.Fatal(result.Err.Error())
	} else {
		t.Log(result.ToString())
	}
	result = conn.ExecSQL("CREATE TABLE T(ID INT, NAME TEXT, X DOUBLE)")
	if result.Err != nil {
		t.Fatal(result.Err.Error())
	} else {
		t.Log(result.ToString())
	}

	for i := 0; i < 3; i++ {
		result := conn.ExecSQL("INSERT INTO T VALUES ("+strconv.Itoa(i)+", '"+string(rune(i+'A'))+"', -" + strconv.Itoa(i) + ")")
		if result.Err != nil {
			t.Fatal(result.Err.Error())
		} else {
			t.Log(result.ToString())
		}
	}

	result = conn.ExecSQL("SELECT 1+2, ID, NAME, X FROM T;")
	if result.Err != nil {
		t.Fatal(result.Err.Error())
	} else {
		t.Log(result.ToString())
	}

	result = conn.ExecSQL("select 9223372036854775807 + 1")
	if result.Err != nil {
		t.Fatal(result.Err.Error())
	} else {
		t.Log(result.ToString())
	}
	result = conn.ExecSQLS("select 9223372036854775807 + 1")
	if result.Err != nil {
		t.Log(result.Err.Error())
	} else {
		t.Fatal(result.ToString())
	}

	testSql := "SELECT (~DEGREES(0.9219647951826007)|FORMAT_BYTES(`f1`)), (~1^`f1`) FROM (SELECT (X^_UTF8MB4'do'-X) AS `f1` FROM (SELECT X FROM T) AS `t1`) AS `t2`;"

	result = conn.ExecSQLS(testSql)
	if result.Err != nil {
		t.Fatal(result.Err.Error())
	} else {
		t.Log(result.ToString())
	}
	outStr, errStr, err := conn.ExecSQLX(testSql, -1)
	t.Log("[out str]", outStr)
	t.Log("[err str]", errStr)
	if err != nil {
		t.Log(err.Error())
	} else {
		t.Fatal("must error!")
	}
}