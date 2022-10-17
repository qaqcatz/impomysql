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
	conn, err := NewConnector(testHost, testPort, testUsername, testPassword, "")
	if err != nil {
		t.Fatal(err.Error())
	}
	result := conn.ExecSQL("CREATE DATABASE IF NOT EXISTS TEST")
	if result.Err != nil {
		t.Fatal(result.Err.Error())
	} else {
		t.Log(result.ToString())
	}

	conn, err = NewConnector(testHost, testPort, testUsername, testPassword, testDBname)
	if err != nil {
		t.Fatal(err.Error())
	}
	result = conn.ExecSQL("DROP TABLE IF EXISTS T")
	if result.Err != nil {
		t.Fatal(result.Err.Error())
	} else {
		t.Log(result.ToString())
	}
	result = conn.ExecSQL("CREATE TABLE T(ID INT, NAME TEXT)")
	if result.Err != nil {
		t.Fatal(result.Err.Error())
	} else {
		t.Log(result.ToString())
	}

	for i := 0; i < 3; i++ {
		result := conn.ExecSQL("INSERT INTO T VALUES ("+strconv.Itoa(i)+", '"+string(rune(i+'A'))+"')")
		if result.Err != nil {
			t.Fatal(result.Err.Error())
		} else {
			t.Log(result.ToString())
		}
	}

	result = conn.ExecSQL("SELECT 1+2, ID, NAME FROM T;")
	if result.Err != nil {
		t.Fatal(result.Err.Error())
	} else {
		t.Log(result.ToString())
	}
}