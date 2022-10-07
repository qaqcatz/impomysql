package randgen

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestRandGen(t *testing.T) {
	sqls, err := RandGen(ZZDefault, YYDefault, 10, 123456)
	if err != nil {
		t.Fatal(err.Error())
	}
	data, err := json.Marshal(sqls)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = ioutil.WriteFile("./test.json", data, 0777)
	if err != nil {
		t.Fatal(err.Error())
	}
}
