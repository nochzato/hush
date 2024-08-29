package hushcore

import (
	"os"
	"path/filepath"
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

func setupTestDir(t *testing.T) (string, func()) {
	t.Helper()

	tempDir, err := os.MkdirTemp("", "hush_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	originalGetHushDir := getHushDir
	getHushDir = func() (string, error) {
		return tempDir, nil
	}

	return tempDir, func() {
		os.RemoveAll(tempDir)
		func() { getHushDir = originalGetHushDir }()
	}
}

func TestAddPassword(t *testing.T) {
	tempDir, clean := setupTestDir(t)
	defer clean()

	masterPassword := "strongMasterPassword123!"
	err := InitHush(masterPassword)
	require.NoError(t, err)

	name := "testname"
	password := "testPassword123!"

	err = AddPassword(name, password, masterPassword)
	require.NoError(t, err)

	filePath := filepath.Join(tempDir, name+".hush")

	if _, err := os.Stat(filePath); err != nil {
		require.NoError(t, err)
	}
}

func TestAddAndGetPassword(t *testing.T) {
	_, clean := setupTestDir(t)
	defer clean()

	masterPassword := "strongMasterPassword123!"
	err := InitHush(masterPassword)
	require.NoError(t, err)

	name := "testname"
	password := "testPassword123!"

	err = AddPassword(name, password, masterPassword)
	require.NoError(t, err)

	retrievedPassword, err := GetPassword(name, masterPassword)
	require.NoError(t, err)

	require.Equal(t, retrievedPassword, password)
}
