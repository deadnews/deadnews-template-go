package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name    string
		envVars map[string]string
		wantErr bool
		check   func(t *testing.T, cfg *Config)
	}{
		{
			name:    "missing SERVICE_DSN",
			envVars: map[string]string{},
			wantErr: true,
		},
		{
			name: "valid config",
			envVars: map[string]string{
				"SERVICE_DSN": "postgres://user:pass@localhost:5432/db",
			},
			check: func(t *testing.T, cfg *Config) {
				assert.Equal(t, "postgres://user:pass@localhost:5432/db", cfg.DSN)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("SERVICE_DSN", "")
			_ = os.Unsetenv("SERVICE_DSN")

			for k, v := range tt.envVars {
				t.Setenv(k, v)
			}

			cfg, err := LoadConfig()

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, cfg)
			if tt.check != nil {
				tt.check(t, cfg)
			}
		})
	}
}
