package wgmux

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUrlMatcherOk1(t *testing.T) {
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
	t.Logf("\n\tin:\t%+v\n\tm:\t%+v\n", in, m)
}

func TestUrlMatcherOk2(t *testing.T) {
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
	t.Logf("\n\tin:\t%+v\n\tm:\t%+v\n", in, m)
}

func TestUrlMatcherOk3(t *testing.T) {
	in := "/make/audi?param=test"
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
	t.Logf("\n\tin:\t%+v\n\tm:\t%+v\n", in, m)
}

func TestUrlMatcherFail(t *testing.T) {
	in := "/audi/model/a4"
	m := "/make/:make/model/:model"
	want1 := false
	got1, _ := urlMatcher(in, m)
	if want1 != got1 {
		t.Errorf("\nwant:\t%+v\ngot:\t%+v", want1, got1)
	}
	t.Logf("\n\tin:\t%+v\n\tm:\t%+v\n", in, m)
}

func TestUrlMatcherFail2(t *testing.T) {
	in := "make/audi/"
	m := "/make/:make/model/:model"
	want1 := false
	got1, _ := urlMatcher(in, m)
	if want1 != got1 {
		t.Errorf("\nwant:\t%+v\ngot:\t%+v", want1, got1)
	}
	t.Logf("\n\tin:\t%+v\n\tm:\t%+v\n", in, m)
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
	eq(in2, args[":arg1"], t)
}

func TestContextHandlerOK(t *testing.T) {
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
	eq(http.StatusOK, rr.Result().StatusCode, t)
}

func TestContextHandlerOK2(t *testing.T) {
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
	eq(http.StatusOK, rr.Result().StatusCode, t)
}

func TestContextHandlerFail(t *testing.T) {
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
	eq(http.StatusOK, rr.Result().StatusCode, t)
}

func TestGetArgMapOK(t *testing.T) {
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
	eq(value, gotMap[key], t)
}

func TestGetArgMapFail(t *testing.T) {
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

func TestServerIntegrationOK1(t *testing.T) {
	m := NewMux()
	h := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	m.HandleFuncRouter("/test/", h)
	ts := httptest.NewServer(m)
	defer ts.Close()
	req, err := http.NewRequest(http.MethodGet, ts.URL+"/test/", nil)
	if err != nil {
		t.Error(err)
	}
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}
	eq(http.StatusOK, r.StatusCode, t)
}

func TestServerIntegrationOK2(t *testing.T) {
	m := NewMux()
	h := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		eq("123", r.URL.Query().Get("abc"), t)
		eq(http.MethodPost, r.Method, t)
	}
	m.HandleFuncRouter("/test/", h)
	ts := httptest.NewServer(m)
	defer ts.Close()
	req, err := http.NewRequest(http.MethodPost, ts.URL+"/test?abc=123", nil)
	if err != nil {
		t.Error(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}
	eq(http.StatusOK, resp.StatusCode, t)
}

func eq(want, got interface{}, t *testing.T) {
	if want != got {
		t.Errorf("\nwant:\t%+v\ngot:\t%+v", want, got)
	}
}
