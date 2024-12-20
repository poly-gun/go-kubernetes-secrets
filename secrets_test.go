package secrets_test

import (
	"context"
	"embed"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/poly-gun/go-kubernetes-secrets"
)

//go:embed test-data
var filesystem embed.FS

func Test(t *testing.T) {
	ctx := context.Background()

	t.Run("New", func(t *testing.T) {
		instance := secrets.New()
		if instance == nil {
			t.Fatalf("New() Returned Nil Map")
		}
	})

	t.Run("Secrets-FS", func(t *testing.T) {
		instance := secrets.New()
		t.Run("Base", func(t *testing.T) {
			e := instance.WalkWithContext(ctx, func(o *secrets.Options) {
				o.FS = filesystem
			})

			if e != nil {
				t.Fatalf("FS() Returned an Error: %v", e)
			}

			for secret, keys := range instance {
				for key := range keys {
					t.Logf("Secret: %s, Key: %s", secret, key)
				}
			}
		})

		t.Run("Old-Secret(s)", func(t *testing.T) {
			for secret, keys := range instance {
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

	t.Run("Secrets-Directory", func(t *testing.T) {
		instance := secrets.New()

		cwd, e := os.Getwd()
		if e != nil {
			t.Fatalf("os.Getwd() returned %v", e)
		}

		target := filepath.Join(cwd, "test-data")

		t.Run("Base", func(t *testing.T) {
			if e := instance.WalkWithContext(ctx, func(o *secrets.Options) {
				o.Directory = target
			}); e != nil {
				t.Fatalf("Walk() Returned an Error: %v", e)
			}

			for secret, keys := range instance {
				for key := range keys {
					t.Logf("Secret: %s, Key: %s", secret, key)
				}
			}
		})

		t.Run("Old-Secret(s)", func(t *testing.T) {
			for secret, keys := range instance {
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
