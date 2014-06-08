package flash

import (
	"fmt"
	"github.com/tsurai/jantar"
	"net/http"
	"net/url"
	"reflect"
)

type Flasher struct {
	data  map[string]string
	respw http.ResponseWriter
}

func NewFlasher(respw http.ResponseWriter) *Flasher {
	return &Flasher{make(map[string]string), respw}
}

func (f *Flasher) Flash(name string, obj interface{}) {
	f.data[name] = fmt.Sprint(obj)
}

func (f *Flasher) FlashStruct(name string, obj interface{}) {
	t := reflect.TypeOf(obj)

	if t.Kind() == reflect.Struct || (t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct) {

		value := reflect.ValueOf(obj)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
			value = value.Elem()
		}

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if field.Tag.Get("jantar") != "noflash" {
				f.data[name+"."+field.Name] = fmt.Sprint(value.Field(i).Interface())
			}
		}
	} else {
		jantar.Log.Errord(jantar.JLData{"expected": "struct", "got": t.Kind()}, "failed to add struct flash. Invalid object type")
	}
}

func (f *Flasher) Save() {
	if len(f.data) > 0 {
		values := url.Values{}
		for key, val := range f.data {
			values.Add(key, val)
		}

		http.SetCookie(f.respw, &http.Cookie{Name: "JANTAR_FLASH", Value: values.Encode(), Secure: false, HttpOnly: true, Path: "/"})
	}
}
