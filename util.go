package util

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func GetFile(prefix string, root string, request string) (http.File, os.FileInfo) {
	var file http.File
	var stat os.FileInfo
	var publicpath, publicfilepath string
	var err error

	fname := request[len(prefix):]

	if !strings.HasPrefix(fname, ".") {
		if publicpath, err = filepath.Abs(root); err == nil {
			if publicfilepath, err = filepath.Abs(root + "/" + fname); err == nil {
				if strings.HasPrefix(publicfilepath, publicpath) {
					if file, err = http.Dir(root).Open(fname); err == nil {
						if stat, err = file.Stat(); err == nil {
							if !stat.IsDir() {
								return file, stat
							}
						}
					}
				}
			}
		}
	}

	return nil, nil
}
