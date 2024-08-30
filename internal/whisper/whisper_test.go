package whisper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckPasswordStrength(t *testing.T) {
	tc := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid",
			password: "strongPass!123",
			wantErr:  false,
		},
		{
			name:     "less then 6",
			password: "Le$1",
			wantErr:  true,
		},
		{
			name:     "no lowercase",
			password: "STRONGPASS!123",
			wantErr:  true,
		},
		{
			name:     "no uppercase",
			password: "strongpass!123",
			wantErr:  true,
		},
		{
			name:     "no numbers",
			password: "strongPass!",
			wantErr:  true,
		},
		{
			name:     "no special",
			password: "strongPass123",
			wantErr:  true,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPasswordStrength(tt.password)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
