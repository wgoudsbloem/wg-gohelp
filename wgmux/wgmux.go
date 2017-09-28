package wgmux

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

type handler func(http.ResponseWriter, *http.Request)

type Argmap map[string]string

type argType string

var Args argType = "args"

type muxer struct {
	handlers map[string]handler
	*http.ServeMux
}

func NewMuxer() muxer {
	m := muxer{make(map[string]handler), http.NewServeMux()}
	m.HandleFunc("/", m.mainHandler)
	return m
}

func (m *muxer) Handle(path string, h handler) {
	m.handlers[path] = h
}

func (m *muxer) Get(path string) (h handler, args Argmap) {
	for k, v := range m.handlers {
		if ok, _args := urlMatcher(path, k); ok {
			h = v
			args = _args
			return
		}
	}
	return
}

func (m *muxer) mainHandler(w http.ResponseWriter, r *http.Request) {
	fn, _args := m.Get(r.URL.Path)
	if fn == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "not found :(")
		return
	}
	ctx := context.WithValue(r.Context(), Args, _args)
	fn(w, r.WithContext(ctx))
}

func urlMatcher(in string, match string) (b bool, m map[string]string) {
	if in[len(in)-1] != '/' {
		in += "/"
	}
	ins := strings.Split(in, "/")
	if match[len(match)-1] != '/' {
		match += "/"
	}
	matches := strings.Split(match, "/")
	if len(ins) != len(matches) {
		b = false
		return
	}
	m = make(Argmap)
	for i := 0; i < len(ins); i++ {
		if ins[i] != matches[i] {
			if matches[i][0] == ':' {
				m[matches[i]] = string(ins[i])
			} else {
				b = false
				return
			}
		}
	}
	b = true
	return
}
