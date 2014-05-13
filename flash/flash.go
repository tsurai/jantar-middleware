package flash

import (
	"github.com/tsurai/jantar"
	"github.com/tsurai/jantar/context"
	"net/http"
	"net/url"
)

type Flash struct {
	jantar.Middleware
}

func (f *Flash) Initialize() {
	tm := jantar.GetModule(jantar.ModuleTemplateManager).(*jantar.TemplateManager)
	tm.AddTmplFunc("flash", func(args map[string]interface{}, key string) string {
		if flashMap, ok := args["_flash"]; ok {
			return flashMap.(map[string]string)[key]
		}
		return ""
	})
}

func (f *Flash) Cleanup() {

}

func (f *Flash) Call(respw http.ResponseWriter, req *http.Request) bool {
	// fetch flash from cookie
	if cookie, err := req.Cookie("JANTAR_FLASH"); err == nil {
		flashMap := make(map[string]string)

		if m, err := url.ParseQuery(cookie.Value); err == nil {
			for key, val := range m {
				flashMap[key] = val[0]
			}
		}

		renderArgs := context.Get(req, "renderArgs").(map[string]interface{})
		renderArgs["_flash"] = flashMap

		// delete cookie
		cookie.MaxAge = -9999
		http.SetCookie(respw, cookie)
	}

	return true
}
