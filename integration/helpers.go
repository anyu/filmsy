package cmd

import (
	"net/http"
	"net/url"
	"github.com/onsi/gomega/ghttp"
)
func setUpMovieDBServer(path string, f func(q url.Values) url.Values) *ghttp.Server {
	dbServer := ghttp.NewServer()
	dbServer.RouteToHandler(http.MethodGet, path, func(w http.ResponseWriter, req *http.Request) {
		q := req.URL.Query()
		f(q)
		req.URL.RawQuery = q.Encode()
		w.Write([]byte(`{"some-info"}`))
	})

	return dbServer
}
