package secrets

import (
	"io/fs"
)

// Options is the configuration structure optionally mutated via the [Settings] constructor used throughout the package.
type Options struct {
	Directory string // Directory represents a string path; Directory is the target volume mount where the kubernetes secret is mounted, and must be a directory.
	FS        fs.FS  // FS represents a [io/fs.FS] (filesystem). If specified, hidden files and directories that start with a dot are ignored. Defaults to nil.
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
