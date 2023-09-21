package fsimpl

import (
	"io"
	"io/fs"
)

type WriteableFile interface {
	fs.File
	io.WriteCloser
}

type WriteableFS interface {
	fs.FS

	OpenFile(string, int, fs.FileMode) (WriteableFile, error)
	Create(string) (WriteableFile, error)
	Mkdir(string, fs.FileMode) error
	MkdirAll(string, fs.FileMode) error
	Remove(string) error
	RemoveAll(string) error
	Rename(string, string) error
}
