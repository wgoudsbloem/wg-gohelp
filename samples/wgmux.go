package main

import (
	"fmt"
	"net/http"

	"github.com/wgoudsbloem/wg-gohelp/wgmux"
)

func main() {
	mux := wgmux.NewMux()
	mux.HandleFuncRouter("/api/:name/", func(w http.ResponseWriter, r *http.Request) {
		argmap := wgmux.GetArgMap(r)
		fmt.Fprintf(w, "hello %s\n", argmap[":name"])
	})
	panic(http.ListenAndServe(":8080", mux))
}
