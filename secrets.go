package secrets

import (
	"context"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/x-ethr/levels"
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

// Walk recursively traverses the specified directory and its subdirectories.
// It collects file paths, directory names, and file contents to build a Secrets map; ignores hidden files and directories that start with a dot.
//   - Returns an error if any occurred during the traversal.
func (s Secrets) Walk(ctx context.Context, directory string) error {
	e := filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !(strings.HasPrefix(d.Name(), ".")) {
			slog.Log(ctx, levels.Trace, "Secrets WalK", slog.String("path", path), slog.String("name", d.Name()), slog.Bool("directory", d.IsDir()))
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

	if e != nil {
		slog.WarnContext(ctx, "Error Walking Directory", slog.String("error", e.Error()))
		return e
	}

	return nil
}

// FS walks the specified file system and populates the Secrets map.
// It ignores hidden files and directories that start with a dot.
// - Returns an error if any occurred during the file system walk.
func (s Secrets) FS(ctx context.Context, filesystem fs.FS) error {
	e := fs.WalkDir(filesystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !(strings.HasPrefix(d.Name(), ".")) {
			slog.Log(ctx, levels.Trace, "Secrets FS Walk", slog.String("path", path), slog.String("name", d.Name()), slog.Bool("directory", d.IsDir()))
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

	if e != nil {
		slog.WarnContext(ctx, "Error Walking Filesystem", slog.String("error", e.Error()))
		return e
	}

	return nil
}

// New returns a new instance of the Secrets type.
// It initializes a Secrets map with an empty map for each secret.
func New() Secrets {
	return make(Secrets)
}

// Walk takes a context and a directory path, and returns a Secrets map and an error.
// It creates a new instance of the Secrets type, then calls the Walk method of that instance with the given context and directory.
// It returns the updated instance and any error that occurred during the Walk operation.
func Walk(ctx context.Context, directory string) (Secrets, error) {
	instance := New()
	e := instance.Walk(ctx, directory)
	return instance, e
}

// FS creates a new instance of Secrets and populates it by walking the provided file system using the given context.
// It returns the populated Secrets and any error encountered during the file system walk.
func FS(ctx context.Context, filesystem fs.FS) (Secrets, error) {
	instance := New()
	e := instance.FS(ctx, filesystem)
	return instance, e
}
