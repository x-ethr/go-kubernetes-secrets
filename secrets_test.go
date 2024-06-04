package secrets

import (
	"context"
	"embed"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

//go:embed test-data
var filesystem embed.FS

func Test(t *testing.T) {
	ctx := context.Background()

	t.Run("New", func(t *testing.T) {
		secrets := New()

		if secrets == nil {
			t.Fatalf("New() Returned Nil Map")
		}
	})

	t.Run("Secrets-FS", func(t *testing.T) {
		secrets := New()

		t.Run("Base", func(t *testing.T) {
			if e := secrets.FS(ctx, filesystem); e != nil {
				t.Fatalf("FS() Returned an Error: %v", e)
			}

			for secret, keys := range secrets {
				for key := range keys {
					t.Logf("Secret: %s, Key: %s", secret, key)
				}
			}
		})

		t.Run("Old-Secret(s)", func(t *testing.T) {
			if e := secrets.FS(ctx, filesystem); e != nil {
				t.Fatalf("FS() Returned an Error: %v", e)
			}

			for secret, keys := range secrets {
				for key := range keys {
					t.Logf("Secret: %s, Key: %s", secret, key)
					value := keys[key]
					t.Logf("Secret: %s, Key: %s, Value: %s", secret, key, value)
					if strings.HasPrefix(string(value), "old") {
						t.Fatalf("..data Value Assigned to Secret, Value")
					}
				}
			}
		})
	})

	t.Run("Secrets-Walk", func(t *testing.T) {
		secrets := New()

		cwd, e := os.Getwd()
		if e != nil {
			t.Fatalf("os.Getwd() returned %v", e)
		}

		target := filepath.Join(cwd, "test-data")

		t.Run("Base", func(t *testing.T) {
			if e := secrets.Walk(ctx, target); e != nil {
				t.Fatalf("Walk() Returned an Error: %v", e)
			}

			for secret, keys := range secrets {
				for key := range keys {
					t.Logf("Secret: %s, Key: %s", secret, key)
				}
			}
		})

		t.Run("Old-Secret(s)", func(t *testing.T) {
			if e := secrets.Walk(ctx, target); e != nil {
				t.Fatalf("FS() Returned an Error: %v", e)
			}

			for secret, keys := range secrets {
				for key := range keys {
					t.Logf("Secret: %s, Key: %s", secret, key)
					value := keys[key]
					t.Logf("Secret: %s, Key: %s, Value: %s", secret, key, value)
					if strings.HasPrefix(string(value), "old") {
						t.Fatalf("..data Value Assigned to Secret, Value")
					}
				}
			}
		})
	})
}
