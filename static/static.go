package static

import (
	"github.com/tsurai/jantar"
	"github.com/tsurai/jantar-middleware"
	"net/http"
)

type Static struct {
	jantar.Middleware
}

func (s *Static) Initialize() {

}

func (s *Static) Cleanup() {

}

func (s *Static) Call(respw http.ResponseWriter, req *http.Request) bool {
	target := req.URL.Path

	if req.URL.Path == "/" {
		target = "/index.html"
	}

	if file, stat := util.GetFile("/", "views/_static", target); file != nil {
		http.ServeContent(respw, req, req.URL.Path, stat.ModTime(), file)
		file.Close()

		return false
	}

	return true
}
