<h3>Simple mux that extends the default http.ServeMux</h3>

install package :  
```go get github.com/wgoudsbloem/wg-gohelp/wgmux```

simple basic usage:
```GO

package main

import (
    "fmt"
    "net/http"
    "wgmux/wg-gohelp/wgmux"
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