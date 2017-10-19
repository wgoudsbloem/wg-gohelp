package main

import (
	"fmt"
	"net/http"
	"wg-gohelp/wgmux"
)

func main() {
	m := wgmux.NewMux()
	m.HandleFuncRouterMethod(http.MethodGet, "/test/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "working")
	})
	panic(http.ListenAndServe(":8081", m))
}
