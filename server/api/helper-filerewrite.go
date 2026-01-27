package api

import (
	"errors"
	"io/fs"
	"net/http"
)

func FileRewrite(fs http.FileSystem, rewrite string) http.FileSystem {
	return fileRewrite{fs: fs, rewrite: rewrite}
}

type fileRewrite struct {
	fs http.FileSystem

	rewrite string
}

func (r fileRewrite) Open(file string) (http.File, error) {
	if result, err := r.fs.Open(file); !errors.Is(err, fs.ErrNotExist) || file == r.rewrite {
		return result, err
	}

	return r.fs.Open(r.rewrite)
}
