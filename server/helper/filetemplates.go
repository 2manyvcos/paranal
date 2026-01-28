package helper

import (
	"bytes"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

func FileTemplates(targetFS fs.FS, data any, targets ...string) http.FileSystem {
	for i, target := range targets {
		targets[i] = strings.TrimPrefix(target, "/")
	}

	httpFS := http.FS(targetFS)
	tmpls, err := template.New("").ParseFS(targetFS, targets...)
	if err != nil {
		return fileTemplate{fs: httpFS, targets: targets, tmplErr: err}
	}
	return fileTemplate{fs: httpFS, targets: targets, tmpls: tmpls, tmplData: data, cache: make(map[string]cacheEntry, len(targets))}
}

type fileTemplate struct {
	fs      http.FileSystem
	targets []string

	tmpls    *template.Template
	tmplData any
	tmplErr  error

	cache map[string]cacheEntry
}

type cacheEntry struct {
	name  string
	data  []byte
	error error
}

func (t fileTemplate) Open(file string) (http.File, error) {
	if file == "/" {
		file = "."
	} else {
		file = strings.TrimPrefix(file, "/")
	}

	if !slices.Contains(t.targets, file) {
		return t.fs.Open(file)
	}

	if t.tmplErr != nil {
		return nil, t.tmplErr
	}

	entry, ok := t.cache[file]
	if !ok {
		var b bytes.Buffer
		if err := t.tmpls.ExecuteTemplate(&b, file, t.tmplData); err != nil {
			entry.error = err
		} else {
			entry.name = filepath.Base(file)
			entry.data = b.Bytes()
		}
		t.cache[file] = entry
	}

	if entry.error != nil {
		return nil, entry.error
	}
	return bufferFile{Reader: bytes.NewReader(entry.data), name: entry.name}, nil
}

type bufferFile struct {
	*bytes.Reader
	name string
}

func (b bufferFile) Close() error                             { return nil }
func (b bufferFile) Readdir(count int) ([]os.FileInfo, error) { return nil, nil }
func (b bufferFile) Stat() (os.FileInfo, error) {
	return bufferFileInfo{name: b.name, size: b.Size()}, nil
}

type bufferFileInfo struct {
	name string
	size int64
}

func (i bufferFileInfo) Name() string       { return i.name }
func (i bufferFileInfo) Size() int64        { return i.size }
func (i bufferFileInfo) Mode() os.FileMode  { return 0444 }
func (i bufferFileInfo) ModTime() time.Time { return time.Now() }
func (i bufferFileInfo) IsDir() bool        { return false }
func (i bufferFileInfo) Sys() any           { return nil }
