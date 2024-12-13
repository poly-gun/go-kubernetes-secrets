package secrets

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// Path represents a file system path as a string, which can be resolved to an absolute system path using the [Path.String] function.
type Path string

// String represents the Path in a fully-resolvable system path. If Path is nil or an empty string, the return value
// will be the caller's [os.Getwd]. String will return an error satisfying errors.Is(err, os.SyscallError) if the
// call to [os.Getwd] fails.
func (p *Path) String() (path string, e error) {
	if p == nil || string(*p) == "" {
		path, e = os.Getwd()
		if e != nil {
			e = fmt.Errorf("unable to determine current working directory for Path type: %w", e)

			return "", e
		}
	}

	path, e = filepath.Abs(path)
	if e != nil {
		e = fmt.Errorf("unable to determine absolute path for Path type: %w", e)
	}

	*p = Path(path)

	return path, e
}

// Options is the configuration structure optionally mutated via the [Settings] constructor used throughout the package.
type Options struct {
	Directory Path  // Directory represents a string Path; Directory is the target volume mount where the kubernetes secret is mounted, and must be a directory.
	FS        fs.FS // FS represents a [io/fs.FS] (filesystem). If specified, hidden files and directories that start with a dot are ignored. Defaults to nil.
}

// options represents a default constructor. The following function isn't publicly exposed as all package functions that
// use the [Settings] function argument will call this function.
func options() *Options {
	return &Options{
		FS: nil,
	}
}

// Settings represents a functional constructor for the [Options] type. Typical callers of Variadic won't need to perform
// nil checks as all implementations first construct an [Options] reference using packaged default(s).
type Settings func(o *Options)
