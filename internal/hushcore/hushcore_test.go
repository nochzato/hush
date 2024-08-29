package hushcore

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSanitizeFileName(t *testing.T) {
	tc := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "valid name",
			input:    "myaccount",
			expected: "myaccount",
			wantErr:  false,
		},
		{
			name:     "name with spaces",
			input:    "my account",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "name with special chars",
			input:    "my&account",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "too long name",
			input:    string(make([]byte, 256)),
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			result, err := sanitizeFileName(tt.input)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, result, tt.expected)
		})
	}
}

func TestAddAndGetPassword(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "hush_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	originalGetHushDir := getHushDir
	getHushDir = func() (string, error) {
		return tempDir, nil
	}
	defer func() { getHushDir = originalGetHushDir }()

	masterPassword := "strongMasterPassword123!"
	err = InitHush(masterPassword)
	require.NoError(t, err)

	name := "testname"
	password := "testPassword123!"

	err = AddPassword(name, password, masterPassword)
	require.NoError(t, err)

	retrievedPassword, err := GetPassword(name, masterPassword)
	require.NoError(t, err)

	require.Equal(t, retrievedPassword, password)
}
