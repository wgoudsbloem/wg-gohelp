package wgmux

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

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

func TestHandlerWithArgs(t *testing.T) {
	in1 := "test1"
	in2 := "test2"
	in3 := "test3"
	in := "/" + in1 + "/" + in2 + "/" + in3
	hdin := func(w http.ResponseWriter, r *http.Request) {}
	mx := NewMux()
	mx.HandleFuncRouter("/test1/:arg1/test3", hdin)
	hndlr, args := mx.handlerWithArgs(in)
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

func TestContextHandler(t *testing.T) {
	testContextHandlerOK(t)
	testContextHandlerFail(t)
}

func testContextHandlerOK(t *testing.T) {
	in1 := "test1"
	in2 := "test2"
	in3 := "test3"
	in := "/" + in1 + "/" + in2 + "/" + in3
	hdin := func(w http.ResponseWriter, r *http.Request) {
		gotMap, ok := GetArgMap(r)
		want := in2
		got := gotMap[":arg1"]
		if !ok || want != got {
			t.Errorf("\nwant:\t%+v\ngot:\t%+v", want, got)
		}
	}
	mx := NewMux()
	mx.HandleFuncRouter("/test1/:arg1/test3", hdin)
	req, err := http.NewRequest("GET", in, nil)
	if err != nil {
		t.Error(err)
	}
	rr := httptest.NewRecorder()
	mx.contextHandler(rr, req)
	want := http.StatusOK
	got := rr.Result().StatusCode
	if want != got {
		t.Errorf("\nwant:\t%+v\ngot:\t%+v", want, got)
	}
}

func testContextHandlerFail(t *testing.T) {
	in1 := "testA"
	in2 := "testB"
	in3 := "testC"
	in := "/" + in1 + "/" + in2 + "/" + in3
	hdin := func(w http.ResponseWriter, r *http.Request) {

	}
	mx := NewMux()
	mx.HandleFuncRouter("/test1/:arg1/test3", hdin)
	req, err := http.NewRequest("GET", in, nil)
	if err != nil {
		t.Error(err)
	}
	rr := httptest.NewRecorder()
	mx.contextHandler(rr, req)
	want := http.StatusNotFound
	got := rr.Result().StatusCode
	if want != got {
		t.Errorf("\nwant:\t%+v\ngot:\t%+v", want, got)
	}
}

func TestGetArgMap(t *testing.T) {
	testGetArgMapOK(t)
	testGetArgMapFail(t)
}

func testGetArgMapOK(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err)
	}
	m := make(map[string]string)
	key := "key"
	value := "value"
	m[key] = value
	ctx := context.WithValue(req.Context(), argString, m)
	gotMap, ok := GetArgMap(req.WithContext(ctx))
	if !ok {
		t.Errorf("\nwant:\t%+v\ngot:\t%+v", m, gotMap)
	}
	got := gotMap[key]
	want := value
	if want != got {
		t.Errorf("\nwant:\t%+v\ngot:\t%+v", want, got)
	}
}

func testGetArgMapFail(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err)
	}
	m := make(map[string]string)
	gotMap, ok := GetArgMap(req)
	if ok {
		t.Errorf("\nwant:\t%+v\ngot:\t%+v", m, gotMap)
	}
}
