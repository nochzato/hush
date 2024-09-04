package passutils

import (
	"testing"
	"testing/quick"

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

func TestGeneratePassword(t *testing.T) {
	f := func(length uint8) bool {
		if length < 6 {
			length = 6
		}

		if length > 255 {
			length = 255
		}

		password, err := GeneratePassword(int(length))
		if err != nil {
			t.Logf("error generating password: %v", err)
			return false
		}

		if len(password) != int(length) {
			t.Logf("password length mismatch. got: %d, want %d", len(password), length)
			return false
		}

		if err := CheckPasswordStrength(password); err != nil {
			t.Logf("generated password does not meet strength criteria: %v", err)
			return false
		}

		return true
	}

	config := &quick.Config{
		MaxCount: 1000,
	}

	if err := quick.Check(f, config); err != nil {
		t.Errorf("test failed: %v", err)
	}
}
