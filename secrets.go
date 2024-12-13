package secrets

import (
	"context"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

// Secret represents the kubernetes secret. On a pod's filesystem, [Secret] value represents the directory where the volume was mounted.
type Secret string

// Key represents a kubernetes secret's key. On a pod's filesystem, [Key] represents a file's name.
type Key string

// Value represents a kubernetes secret's value. On a pod's filesystem, [Value] represents the [Key] file's contents.
type Value string

func (v Value) Bytes() []byte {
	return []byte(v)
}

// Secrets represents a map[string]map[string][]byte mapping of [Secret] -> [Key] -> [Value].
type Secrets map[Secret]map[Key]Value

// walk recursively traverses directories or file systems based on configuration settings and populates the Secrets map, ignoring hidden files and directories.
func (s Secrets) walk(ctx context.Context, settings ...Settings) (e error) {
	var o = options()
	for _, configuration := range settings {
		configuration(o)
	}

	switch {
	case o.FS != nil:
		slog.DebugContext(ctx, "Walking System via FS Configuration", slog.Any("configuration", o))

		e = fs.WalkDir(o.FS, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !(strings.HasPrefix(d.Name(), ".")) {
				slog.DebugContext(ctx, "FS Secrets Walk", slog.String("path", path), slog.String("name", d.Name()), slog.Bool("directory", d.IsDir()))
				if d.IsDir() {
					secret := Secret(d.Name())
					s[secret] = make(map[Key]Value)
					return nil
				}

				key := Key(d.Name())
				secret := Secret(filepath.Base(filepath.Dir(path)))
				if strings.HasPrefix(string(secret), ".") {
					// --> avoid ..data and .symbolic-link directories
					secret = Secret(filepath.Base(filepath.Dir(filepath.Dir(path))))
				}

				value, exception := os.ReadFile(path)
				if exception != nil {
					return exception
				}

				s[secret][key] = Value(value)
			}

			return nil
		})
	default:
		slog.DebugContext(ctx, "Walking System via Directory Configuration", slog.Any("configuration", o))

		var directory string
		directory, e = o.Directory.String()

		if e != nil {
			slog.ErrorContext(ctx, "Unable to Compute Directory String Literal", slog.String("error", e.Error()), slog.Any("configuration", o))

			return e
		}

		e = filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !(strings.HasPrefix(d.Name(), ".")) {
				slog.DebugContext(ctx, "Directory Secrets WalK", slog.String("path", path), slog.String("name", d.Name()), slog.Bool("directory", d.IsDir()))
				if d.IsDir() {
					secret := Secret(d.Name())
					s[secret] = make(map[Key]Value)
					return nil
				}

				key := Key(d.Name())
				secret := Secret(filepath.Base(filepath.Dir(path)))
				if strings.HasPrefix(string(secret), ".") {
					// --> avoid ..data and .symbolic-link directories
					secret = Secret(filepath.Base(filepath.Dir(filepath.Dir(path))))
				}

				value, exception := os.ReadFile(path)
				if exception != nil {
					return exception
				}

				s[secret][key] = Value(value)
			}

			return nil
		})
	}

	if e != nil {
		slog.ErrorContext(ctx, "Unable to Walk Directory", slog.String("error", e.Error()), slog.Any("configuration", o))
		return e
	}

	return nil
}

// Walk recursively traverses the specified directory and its subdirectories.
// It collects file paths, directory names, and file contents to build a Secrets map; ignores hidden files and directories that start with a dot.
//   - Returns an error if any occurred during the traversal.
func (s Secrets) Walk(settings ...Settings) error {
	ctx := context.Background()

	return s.walk(ctx, settings...)
}

// WalkWithContext traverses the filesystem or directories using the provided context and settings, populating Secrets. Returns an error if traversal fails.
func (s Secrets) WalkWithContext(ctx context.Context, settings ...Settings) error {
	return s.walk(ctx, settings...)
}

// New returns a new instance of the Secrets type.
// It initializes a Secrets map with an empty map for each secret.
func New() Secrets {
	return make(Secrets)
}
