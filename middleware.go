package cameljson

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
)

func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		w := &response{origin: resp}
		h.ServeHTTP(w, req)
		w.Flush()
	})
}

type response struct {
	buffer bytes.Buffer
	origin http.ResponseWriter
}

func (w *response) Header() http.Header {
	return w.origin.Header()
}

func (w *response) Write(p []byte) (n int, err error) {
	return w.buffer.Write(p)
}

func (w *response) WriteHeader(code int) {
	w.origin.WriteHeader(code)
}

func (w *response) Flush() error {
	var o interface{}
	var data = w.buffer.Bytes()

	if err := json.Unmarshal(data, &o); err != nil {
		_, werr := w.buffer.WriteTo(w.origin)
		return werr
	}

	adaptFieldNames(o)

	encoder := json.NewEncoder(w.origin)
	return encoder.Encode(o)
}

func adaptFieldNames(o interface{}) {
	if f, ok := o.(map[string]interface{}); ok {
		for k, v := range f {
			delete(f, k)
			name := strings.ToLower(k[:1]) + k[1:]
			f[name] = v
			adaptFieldNames(v)
		}
	}

	if f, ok := o.([]interface{}); ok {
		for _, v := range f {
			adaptFieldNames(v)
		}
	}
}
