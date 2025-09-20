package config_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"pkgstats-cli/internal/config"
)

func TestLoad(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		content string
		want    *config.Config
		wantErr bool
	}{
		{
			name: "valid config",
			content: `
blocklist:
  packages:
    - "package-*"
  mirrors:
    - "mirror.com"
`,
			want: &config.Config{
				Blocklist: struct {
					Packages []string `mapstructure:"packages"`
					Mirrors  []string `mapstructure:"mirrors"`
				}{
					Packages: []string{"package-*"},
					Mirrors:  []string{"mirror.com"},
				},
			},
			wantErr: false,
		},
		{
			name:    "empty config",
			content: "",
			want:    &config.Config{},
			wantErr: false,
		},
		{
			name: "malformed config",
			content: `
blocklist:
  packages: ["one"
`,
			want:    nil,
			wantErr: true,
		},
		{
			name: "config with unknown field",
			content: `
blocklist:
  packages:
    - "package-*"
  mirrors:
    - "mirror.com"
unknown_field: "value"
`,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tmpDir := t.TempDir()
			configFile := filepath.Join(tmpDir, "pkgstats.yaml")

			if err := os.WriteFile(configFile, []byte(tt.content), 0o600); err != nil {
				t.Fatal(err)
			}

			got, err := config.Load(configFile)

			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load() = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("missing config file", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		configFile := filepath.Join(tmpDir, "nonexistent.yaml")

		got, err := config.Load(configFile)
		if err != nil {
			t.Errorf("Load() error = %v, wantErr false", err)
		}

		if !reflect.DeepEqual(got, &config.Config{}) {
			t.Errorf("Load() = %v, want %v", got, &config.Config{})
		}
	})
}
