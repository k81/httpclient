package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/k81/log"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		if r.Form.Get("hello") == "world" {
			fmt.Fprintf(w, "hello world")
		} else {
			fmt.Fprintf(w, "bad hello")
		}
	}))

	ctx := context.TODO()
	client := New(Timeout(time.Second*5), DisableRedirect)

	query := url.Values{}
	query.Add("hello", "world")

	result, err := client.Get(ctx, server.URL, "", SetQuery(query))
	require.NoError(t, err)
	require.Equal(t, "hello world", result)
}

func TestPost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		if r.Form.Get("a") == "1" && r.Form.Get("b") == "2" {
			fmt.Fprintf(w, "hello world")
		} else {
			fmt.Fprintf(w, "bad hello")
		}
	}))
	ctx := context.TODO()
	client := New(Timeout(time.Second*5), DisableRedirect)

	form := url.Values{}
	form.Add("a", "1")
	form.Add("b", "2")
	result, err := client.Post(ctx, server.URL, form.Encode(), SetTypeForm())
	require.NoError(t, err)
	require.Equal(t, "hello world", result)
}

func TestJSONPost(t *testing.T) {
	type Hello struct {
		Hello string `json:"hello"`
	}

	type HelloResult struct {
		ErrNo  int    `json:"errno"`
		ErrMsg string `json:"errmsg"`
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		h := &Hello{}
		err = json.Unmarshal(data, h)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if h.Hello == "world" {
			fmt.Fprintf(w, `{"errno":0, "errmsg":"hello world"}`)
		} else {
			fmt.Fprintf(w, `{"errno":1, "errmsg":"bad hello"}`)
		}

	}))

	ctx := context.TODO()
	client := NewJSON(Timeout(time.Second*5), DisableRedirect)

	hello := &Hello{
		Hello: "world",
	}

	result := &HelloResult{}

	err := client.Post(ctx, server.URL, hello, result, SetTypeJSON())
	require.NoError(t, err)
	require.Equal(t, 0, result.ErrNo)
	require.Equal(t, "hello world", result.ErrMsg)
}

func TestLogContextFunc(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		if r.Form.Get("hello") == "world" {
			fmt.Fprintf(w, "hello world")
		} else {
			fmt.Fprintf(w, "bad hello")
		}
	}))

	ctx := context.TODO()
	client := New(Timeout(time.Second*5), DisableRedirect)
	client.SetLogContextFunc(func(ctx context.Context, req *http.Request) context.Context {
		return log.WithContext(ctx, "log_method", req.Method)
	})

	query := url.Values{}
	query.Add("hello", "world")

	result, err := client.Get(ctx, server.URL, "", SetQuery(query))
	require.NoError(t, err)
	require.Equal(t, "hello world", result)
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
