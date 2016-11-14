package cameljson

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestMiddleware(t *testing.T) {
	cases := []struct {
		in   []byte
		want []byte
	}{
		{
			[]byte(`{"NamE":"Ivan","AGE":26,"_sex":"male"}`),
			[]byte(`{"namE":"Ivan","age":26,"_sex":"male"}`),
		},
		{
			[]byte(`[{"FIRST_NAME":"Ivan"},{"age":26}]`),
			[]byte(`[{"first_name":"Ivan"},{"age":26}]`),
		},
		{
			[]byte(`{"MaP":{"PaM":"amp"},array:[1,2,3]}`),
			[]byte(`{"maP":{"paM":"amp"},array:[1,2,3]}`),
		},
	}

	for _, c := range cases {
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "http://example.com/", nil)

		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write(c.in)
		})

		Middleware(h).ServeHTTP(rec, req)

		got := rec.Body.Bytes()

		if !jsonEqual(got, c.want) {
			t.Errorf("expected %s", c.want)
			t.Errorf("     got %s", got)
		}
	}
}

func TestMiddlewareWithBasicString(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{
			"invalid_json",
			"invalid_json",
		},
		{
			"",
			"",
		},
	}

	for _, c := range cases {
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "http://example.com/", nil)

		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(c.in))
		})

		Middleware(h).ServeHTTP(rec, req)

		got := rec.Body.String()

		if got != c.want {
			t.Errorf("expected '%s'", c.want)
			t.Errorf("     got '%s'", got)
		}
	}
}

func jsonEqual(a, b []byte) bool {
	var i1, i2 interface{}

	json.Unmarshal(a, &i1)
	json.Unmarshal(b, &i2)

	return reflect.DeepEqual(i1, i2)
}
