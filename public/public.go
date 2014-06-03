package public

import (
	"github.com/tsurai/jantar"
	"github.com/tsurai/jantar-middleware"
	"net/http"
	"strings"
)

type Public struct {
	jantar.Middleware
}

func (p *Public) Initialize() {

}

func (p *Public) Cleanup() {

}

func (p *Public) Call(respw http.ResponseWriter, req *http.Request) bool {
	if strings.HasPrefix(req.URL.Path, "/public/") {
		if file, stat := util.GetFile("/public/", "public", req.URL.Path); file != nil {
			http.ServeContent(respw, req, req.URL.Path, stat.ModTime(), file)
			file.Close()

			return false
		}
	}

	return true
}
