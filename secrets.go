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

type Secret string // Secret represents the kubernetes secret. On a pod's filesystem, [Secret] value represents the directory where the volume was mounted.
type Key string    // Key represents a kubernetes secret's key. On a pod's filesystem, [Key] represents a file's name.

type Value []byte // Value represents a kubernetes secret's value. On a pod's filesystem, [Value] represents the [Key] file's binary contents.
func (v Value) String() string {
	return string(v)
}

type Secrets map[Secret]map[Key]Value

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

			s[secret][key] = value
		}

		return nil
	})

	if e != nil {
		slog.WarnContext(ctx, "Error Walking Directory", slog.String("error", e.Error()))
		return e
	}

	return nil
}

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

			s[secret][key] = value
		}

		return nil
	})

	if e != nil {
		slog.WarnContext(ctx, "Error Walking Filesystem", slog.String("error", e.Error()))
		return e
	}

	return nil
}

func New() Secrets {
	return make(Secrets)
}
