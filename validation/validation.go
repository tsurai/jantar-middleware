package validation

import (
	"github.com/tsurai/jantar"
	"github.com/tsurai/jantar/context"
	"net/http"
	"net/url"
)

type Validation struct {
	jantar.Middleware
}

func (v *Validation) Initialize() {
	tm := jantar.GetModule(jantar.ModuleTemplateManager).(*jantar.TemplateManager)

	tm.AddTmplFunc("errors", func(args map[string]interface{}) map[string][]string {
		if errorMap, ok := args["_errors"]; ok {
			return errorMap.(map[string][]string)
		}
		return nil
	})

	tm.AddTmplFunc("hasError", func(args map[string]interface{}, key string) bool {
		if errorMap, ok := args["_errors"]; ok {
			_, ok := errorMap.(map[string][]string)[key]
			return ok
		}
		return false
	})
}

func (v *Validation) Cleanup() {

}

func (v *Validation) Call(respw http.ResponseWriter, req *http.Request) bool {
	// fetch flash from cookie
	if cookie, err := req.Cookie("JANTAR_ERRORS"); err == nil {
		errorMap := make(map[string][]string)

		if m, err := url.ParseQuery(cookie.Value); err == nil {
			errorMap = m
		}

		renderArgs := context.Get(req, "renderArgs").(map[string]interface{})
		renderArgs["_errors"] = errorMap

		// delete cookie
		cookie.MaxAge = -9999
		http.SetCookie(respw, cookie)
	}

	return true
}
