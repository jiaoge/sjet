package engine

import (
	"io"
	"os"
	"path/filepath"

	"github.com/CloudyKit/jet/v6"
	"github.com/sirupsen/logrus"
)

// 自定义文件系统loader，满足和Mem相同的操作行为
// OSFileSystemLoader implements Loader interface using OS file system (os.File).
type OSFileSystemLoader struct {
	dir    string
	Views  *jet.Set
	ECache *ECache
}

// NewOSFileSystemLoader returns an initialized OSFileSystemLoader.
func NewOSFileSystemLoader(dirPath string, views *jet.Set, ecache *ECache) *OSFileSystemLoader {
	return &OSFileSystemLoader{
		dir:    filepath.FromSlash(dirPath),
		Views:  views,
		ECache: ecache,
	}
}

// Exists returns true if a file is found under the template path after converting it to a file path
// using the OS's path seperator and joining it with the loader's directory path.
func (l *OSFileSystemLoader) Exists(templatePath string) bool {
	templatePath = filepath.Join(l.dir, filepath.FromSlash(templatePath))
	stat, err := os.Stat(templatePath)
	if err == nil && !stat.IsDir() {
		return true
	}
	return false
}

// Open returns the result of `os.Open()` on the file located using the same logic as Exists().
func (l *OSFileSystemLoader) Open(templatePath string) (io.ReadCloser, error) {
	return os.Open(filepath.Join(l.dir, filepath.FromSlash(templatePath)))
}

func (l *OSFileSystemLoader) Set(templatePath string, contents string) {
	t, err := l.Views.Parse(templatePath, contents)
	if err != nil {
		logrus.Error(err)
		return
	}
	l.ECache.Put(templatePath, t)
}

func (l *OSFileSystemLoader) Delete(templatePath string) {
	l.ECache.Del(templatePath)
}
