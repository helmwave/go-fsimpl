// Package autofs provides the ability to look up all filesystems supported by
// this module. Using this package will compile a great many dependencies into
// the resulting binary, so unless you need to support all supported filesystems,
// use fsimpl.NewMux instead.
package autofs

import (
	"io/fs"
	"sync"

	"github.com/helmwave/go-fsimpl"
	"github.com/helmwave/go-fsimpl/awssmfs"
	"github.com/helmwave/go-fsimpl/awssmpfs"
	"github.com/helmwave/go-fsimpl/blobfs"
	"github.com/helmwave/go-fsimpl/consulfs"
	"github.com/helmwave/go-fsimpl/filefs"
	"github.com/helmwave/go-fsimpl/gitfs"
	"github.com/helmwave/go-fsimpl/httpfs"
	"github.com/helmwave/go-fsimpl/vaultfs"
)

//nolint:gochecknoglobals
var (
	mux     fsimpl.FSMux
	muxInit sync.Once
)

// Lookup returns an appropriate filesystem for the given URL.
// If a filesystem can't be found for the provided URL's scheme, an error will
// be returned.
func Lookup(u string) (fs.FS, error) {
	muxInit.Do(func() {
		mux = fsimpl.NewMux()
		mux.Add(awssmfs.FS)
		mux.Add(awssmpfs.FS)
		mux.Add(blobfs.FS)
		mux.Add(consulfs.FS)
		mux.Add(filefs.FS)
		mux.Add(gitfs.FS)
		mux.Add(httpfs.FS)
		mux.Add(vaultfs.FS)
	})

	return mux.Lookup(u)
}
