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
	m.HandleFunc("/", m.mainHandler)
	return m
}

func (m *mux) HandleFuncRouter(path string, h handler) {
	m.handlers[path] = h
}

func (m *mux) get(path string) (h handler, args map[string]string) {
	for k, v := range m.handlers {
		if ok, _args := urlMatcher(path, k); ok {
			h = v
			args = _args
			return
		}
	}
	return
}

func (m *mux) mainHandler(w http.ResponseWriter, r *http.Request) {
	fn, _args := m.get(r.URL.Path)
	if fn == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "not found :(")
		return
	}
	ctx := context.WithValue(r.Context(), argString, _args)
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
func GetArgMap(r *http.Request) map[string]string {
	return r.Context().Value(argString).(map[string]string)
}
