package wgmux

import "testing"

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
