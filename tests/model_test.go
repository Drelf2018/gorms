package tests_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Drelf2018/gorms"
	"github.com/Drelf2018/gorms/tests"
)

func TestFind(t *testing.T) {
	us, ok := gorms.Find[tests.User, bool]()
	if !ok {
		t.Fail()
	}
	b, _ := json.Marshal(us)
	println(string(b))
}

func TestPreloads(t *testing.T) {
	time.Sleep(time.Second * 2)
	us, err := gorms.Preloads[tests.User, error]()
	if err != nil {
		t.Fatal(err)
	}
	b, _ := json.Marshal(us)
	println(string(b))
}
