package wgmux

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

type handler func(http.ResponseWriter, *http.Request)

type argType string

var argString argType = "args"

type mux struct {
	handlers map[string]handler
	*http.ServeMux
}

// NewMux will return a mux with all http.ServeMux methods
// and an additional router based func: HandleFuncRouter
func NewMux() mux {
	m := mux{make(map[string]handler), http.NewServeMux()}
	m.HandleFunc("/", m.contextHandler)
	return m
}

func (m *mux) HandleFuncRouter(path string, h handler) {
	m.handlers[path] = h
}

// handlerWithArgs the handler matched to the path and a map with path aruments
// passed path string
func (m *mux) handlerWithArgs(path string) (h handler, args map[string]string) {
	for k, v := range m.handlers {
		if ok, _args := urlMatcher(path, k); ok {
			h = v
			args = _args
			return
		}
	}
	return
}

func (m *mux) contextHandler(w http.ResponseWriter, r *http.Request) {
	fn, _args := m.handlerWithArgs(r.URL.Path)
	if fn == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "not found :(")
		return
	}
	ctx := context.WithValue(r.Context(), argString, _args)
	fn(w, r.WithContext(ctx))
}

func urlMatcher(in, match string) (b bool, m map[string]string) {
	// strip off everything after a possible ? and add a /
	for i, v := range in {
		if v == '?' {
			in = in[:i]
		}
	}
	if in[len(in)-1] != '/' {
		in += "/"
	}
	// add a trailing slash if it doesn't have one to match up the in
	if match[len(match)-1] != '/' {
		match += "/"
	}
	matches := strings.Split(match, "/")
	ins := strings.Split(in, "/")
	if len(ins) != len(matches) {
		b = false
		return
	}
	m = make(map[string]string)
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

// GetArgMap retrieves the path arguments in a map[string]string
func GetArgMap(r *http.Request) (map[string]string, bool) {
	ctx := r.Context().Value(argString)
	if ctx == nil {
		return nil, false
	}
	m, ok := ctx.(map[string]string)
	return m, ok
}
