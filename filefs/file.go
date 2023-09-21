// Package filefs wraps os.DirFS to provide a local filesystem for file:// URLs.
package filefs

import (
	"io/fs"
	"net/url"
	"os"
	"path/filepath"

	"github.com/helmwave/go-fsimpl"
	"github.com/helmwave/go-fsimpl/internal"
)

type fileFS struct {
	root string
}

// New returns a filesystem (an fs.FS) for the tree of files rooted at the
// directory root. This filesystem is suitable for use with the 'file:' URL
// scheme, and interacts with the local filesystem.
//
// This is effectively a wrapper for os.DirFS.
func New(u *url.URL) (fs.FS, error) {
	rootPath := pathForDirFS(u)

	return &fileFS{root: rootPath}, nil
}

// return the correct filesystem path for the given URL. Supports Windows paths
// and UNCs as well
func pathForDirFS(u *url.URL) string {
	if u.Path == "" {
		return ""
	}

	rootPath := u.Path
	if len(rootPath) >= 3 {
		if rootPath[0] == '/' && rootPath[2] == ':' {
			rootPath = rootPath[1:]
		}
	}

	// a file:// URL with a host part should be interpreted as a UNC
	switch u.Host {
	case ".":
		rootPath = "//./" + rootPath
	case "":
		// nothin'
	default:
		rootPath = "//" + u.Host + rootPath
	}

	return rootPath
}

// FS is used to register this filesystem with an fsimpl.FSMux
//
//nolint:gochecknoglobals
var FS = fsimpl.FSProviderFunc(New, "file")

var (
	_ fs.FS                = (*fileFS)(nil)
	_ fs.ReadDirFS         = (*fileFS)(nil)
	_ fs.ReadFileFS        = (*fileFS)(nil)
	_ fs.StatFS            = (*fileFS)(nil)
	_ fs.GlobFS            = (*fileFS)(nil)
	_ fs.SubFS             = (*fileFS)(nil)
	_ internal.WriteableFS = (*fileFS)(nil)
)

func (f *fileFS) dirFS() fs.FS {
	return os.DirFS(f.root)
}

func (f *fileFS) join(path string) string {
	return filepath.Join(f.root, path)
}

func (f *fileFS) Open(name string) (fs.File, error) {
	return f.dirFS().Open(name)
}

func (f *fileFS) ReadFile(name string) ([]byte, error) {
	return fs.ReadFile(f.dirFS(), name)
}

func (f *fileFS) ReadDir(name string) ([]fs.DirEntry, error) {
	return fs.ReadDir(f.dirFS(), name)
}

func (f *fileFS) Stat(name string) (fs.FileInfo, error) {
	return fs.Stat(f.dirFS(), name)
}

func (f *fileFS) Glob(name string) ([]string, error) {
	return fs.Glob(f.dirFS(), name)
}

func (f *fileFS) Sub(name string) (fs.FS, error) {
	return fs.Sub(f.dirFS(), name)
}

func (f *fileFS) OpenFile(name string, flag int, perm fs.FileMode) (internal.WriteableFile, error) {
	return os.OpenFile(f.join(name), flag, perm)
}

func (f *fileFS) Create(name string) (internal.WriteableFile, error) {
	return os.Create(f.join(name))
}

func (f *fileFS) Mkdir(name string, perm fs.FileMode) error {
	return os.Mkdir(f.join(name), perm)
}

func (f *fileFS) MkdirAll(path string, perm fs.FileMode) error {
	return os.MkdirAll(f.join(path), perm)
}

func (f *fileFS) Remove(name string) error {
	return os.Remove(name)
}

func (f *fileFS) RemoveAll(path string) error {
	return os.RemoveAll(f.join(path))
}

func (f *fileFS) Rename(oldpath, newpath string) error {
	return os.Rename(f.join(oldpath), f.join(newpath))
}
