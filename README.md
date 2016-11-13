# cameljson
cameljson provides a convenient way to adapt your json HTTP response in camel case style.

## Why?
How often do you need to add ```json``` tag to your struct fields only to change the name in camel case style? According to Google JSON Style Guide (https://google.github.io/styleguide/jsoncstyleguide.xml) the camel case format is must. cameljson is a minimal and lightweight middleware that does the adapting for you.

## Usage
Your handlers should implement the standard http.Handler interface

```go
type Handler interface {
        ServeHTTP(ResponseWriter, *Request)
}
```

This complete example shows the full power of cameljson.

```go
package main

import (
    "net/http"

    "github.com/FeNoMeNa/cameljson"
)

func myApp(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte(`{"Name":"Ivan"}`))
}

func main() {
    myHandler := http.HandlerFunc(myApp)
    mux := cameljson.Middleware(myHandler)

    http.ListenAndServe(":8000", mux)
}
```

Here myApp handler send ```{"Name":"Ivan"}``` json and then cameljson middleware adapts it to ```{"name":"Ivan"}``` (camel case).
