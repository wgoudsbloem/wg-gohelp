<h3>Simple mux that extends the default http.ServeMux</h3>

[![Build Status](https://travis-ci.org/wgoudsbloem/wg-gohelp.svg?branch=master)](https://travis-ci.org/wgoudsbloem/wg-gohelp)
[![Coverage Status](https://coveralls.io/repos/github/wgoudsbloem/wg-gohelp/badge.png?branch=master)](https://coveralls.io/github/wgoudsbloem/wg-gohelp?branch=master)

install package :  
```go get github.com/wgoudsbloem/wg-gohelp/wgmux```

simple basic usage:
```GO

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

```