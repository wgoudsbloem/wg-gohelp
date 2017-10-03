package wgmux

import "testing"
import "net/http"

func TestUrlMatcher(t *testing.T) {
	testUrlMatcherOk1(t)
	testUrlMatcherOk2(t)
	testUrlMatcherFail(t)
	testUrlMatcherFail2(t)
}

func testUrlMatcherOk1(t *testing.T) {
	in := "/make/audi/model/a4"
	m := "/make/:make/model/:model"
	want1 := true
	got1, got2 := urlMatcher(in, m)
	if want1 != got1 {
		t.Errorf("\nwant:\t%+v\ngot:\t%+v", want1, got1)
	}
	want2 := map[string]string{}
	want2[":make"] = "audi"
	want2[":model"] = "a4"
	if got2[":make"] != want2[":make"] {
		t.Errorf("\nwant:\t%+v\ngot:\t%+v", want2, got2)
	}
}

func testUrlMatcherOk2(t *testing.T) {
	in := "/make/audi/"
	m := "/make/:make/"
	want1 := true
	got1, got2 := urlMatcher(in, m)
	if want1 != got1 {
		t.Errorf("\nwant:\t%+v\ngot:\t%+v", want1, got1)
	}
	want2 := map[string]string{}
	want2[":make"] = "audi"
	want2[":model"] = "a4"
	if got2[":make"] != want2[":make"] {
		t.Errorf("\nwant:\t%+v\ngot:\t%+v", want2, got2)
	}
}

func testUrlMatcherFail(t *testing.T) {
	in := "/audi/model/a4"
	m := "/make/:make/model/:model"
	want1 := false
	got1, _ := urlMatcher(in, m)
	if want1 != got1 {
		t.Errorf("\nwant:\t%+v\ngot:\t%+v", want1, got1)
	}
}

func testUrlMatcherFail2(t *testing.T) {
	in := "make/audi/"
	m := "/make/:make/model/:model"
	want1 := false
	got1, _ := urlMatcher(in, m)
	if want1 != got1 {
		t.Errorf("\nwant:\t%+v\ngot:\t%+v", want1, got1)
	}
}

func TestGet(t *testing.T) {
	in1 := "test1"
	in2 := "test2"
	in3 := "test3"
	in := "/" + in1 + "/" + in2 + "/" + in3
	hdin := func(w http.ResponseWriter, r *http.Request) {}
	mx := NewMux()
	mx.HandleFuncRouter("/test1/:arg1/test3", hdin)
	hndlr, args := mx.get(in)
	if hndlr == nil {
		t.Error("handler cannot be nil")
	}
	if args == nil {
		t.Error("args can not be nil")
	}
	want := in2
	got := args[":arg1"]
	if want != got {
		t.Errorf("\nwant:\t%+v\ngot:\t%+v", want, got)
	}
}
