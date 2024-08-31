package hushcore

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/nochzato/hush/internal/whisper"
)

const (
	hushDirName        = ".hush"
	masterHashFileName = "master.hash"
	saltFileName       = "salt"
)

var getHushDir = defaultGetHushDir

func defaultGetHushDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, hushDirName), nil
}

func readSalt() (string, error) {
	hushDir, err := getHushDir()
	if err != nil {
		return "", err
	}

	saltFile := filepath.Join(hushDir, saltFileName)
	salt, err := os.ReadFile(saltFile)
	if err != nil {
		return "", fmt.Errorf("failed to read salt: %w", err)
	}

	return string(salt), err
}

func InitHush(masterPassword string) error {
	hushDir, err := getHushDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(hushDir, 0700); err != nil {
		return fmt.Errorf("failed to create hush directory: %w", err)
	}

	masterPasswordFile := filepath.Join(hushDir, masterHashFileName)
	saltFile := filepath.Join(hushDir, saltFileName)

	if _, err := os.Stat(masterPasswordFile); err == nil {
		return fmt.Errorf("hush is already initialized")
	}

	if err := whisper.CheckPasswordStrength(masterPassword); err != nil {
		return fmt.Errorf("master password is too weak: %w", err)
	}

	key, salt, err := whisper.DeriveKey(masterPassword)
	if err != nil {
		return fmt.Errorf("failed to derive key: %w", err)
	}

	encryptedMasterPassword, err := whisper.EncryptPassword(masterPassword, key)
	if err != nil {
		return fmt.Errorf("failed to encrypt master password: %w", err)
	}

	if err := os.WriteFile(masterPasswordFile, []byte(encryptedMasterPassword), 0600); err != nil {
		return fmt.Errorf("failed to save encrypted master password: %w", err)
	}

	if err := os.WriteFile(saltFile, []byte(salt), 0600); err != nil {
		return fmt.Errorf("failed to save salt: %w", err)
	}

	fmt.Println("Hush initialized successfully!")
	return nil
}

func AddPassword(name, password, masterPassword string) error {
	sanitizedName, err := sanitizeFileName(name)
	if err != nil {
		return fmt.Errorf("invalid filename: %w", err)
	}

	if err := whisper.CheckPasswordStrength(password); err != nil {
		return fmt.Errorf("password is too weak: %w", err)
	}

	salt, err := readSalt()
	if err != nil {
		return fmt.Errorf("failed to read salt: %w", err)
	}

	decryptedMasterPassword, err := validateMasterPassword(salt, masterPassword)
	if err != nil {
		return fmt.Errorf("error validating master password: %w", err)
	}

	encryptionKey, err := whisper.DeriveKeyWithSalt(decryptedMasterPassword, salt)
	if err != nil {
		return fmt.Errorf("failed to derive encryption key: %w", err)
	}

	encryptedPassword, err := whisper.EncryptPassword(password, encryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt password: %w", err)
	}

	err = savePassword(sanitizedName, encryptedPassword)
	if err != nil {
		return fmt.Errorf("failed to save password: %w", err)
	}

	return nil
}

func savePassword(name, encryptedPassword string) error {
	hushDir, err := getHushDir()
	if err != nil {
		return err
	}

	filePath := filepath.Join(hushDir, name+".hush")

	err = os.WriteFile(filePath, []byte(encryptedPassword), 0600)
	if err != nil {
		return fmt.Errorf("failed to write password file: %w", err)
	}
	return nil
}

func GetPassword(name, masterPassword string) (string, error) {
	sanitizedName, err := sanitizeFileName(name)
	if err != nil {
		return "", fmt.Errorf("invalid filename: %w", err)
	}

	hushDir, err := getHushDir()
	if err != nil {
		return "", err
	}

	encryptedPassword, err := readEncryptedPassword(hushDir, sanitizedName)
	if err != nil {
		return "", fmt.Errorf("failed to read password file: %w", err)
	}

	salt, err := readSalt()
	if err != nil {
		return "", fmt.Errorf("failed to read salt: %w", err)
	}

	decryptedMasterPassword, err := validateMasterPassword(salt, masterPassword)
	if err != nil {
		return "", fmt.Errorf("error validating master password: %w", err)
	}

	encryptionKey, err := whisper.DeriveKeyWithSalt(decryptedMasterPassword, salt)
	if err != nil {
		return "", fmt.Errorf("failed to derive encryption key: %w", err)
	}

	decryptedPassword, err := whisper.DecryptPassword(string(encryptedPassword), encryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt password: %w", err)
	}

	trimmedPassword := strings.TrimSpace(decryptedPassword)

	return trimmedPassword, nil
}

func validateMasterPassword(salt, masterPassword string) (decryptedMasterPassword string, err error) {
	encryptedMasterPassword, err := readEncryptedMasterPassword()
	if err != nil {
		return "", fmt.Errorf("failed to read encrypted master password: %w", err)
	}

	key, err := whisper.DeriveKeyWithSalt(masterPassword, salt)
	if err != nil {
		return "", fmt.Errorf("failed to derive key: %w", err)
	}

	decryptedMasterPassword, err = whisper.DecryptPassword(encryptedMasterPassword, key)
	if err != nil {
		return "", fmt.Errorf("incorrect master password")
	}

	return decryptedMasterPassword, nil
}

func readEncryptedPassword(hushDir, name string) ([]byte, error) {
	filePath := filepath.Join(hushDir, name+".hush")

	encryptedPassword, err := os.ReadFile(filePath)
	return encryptedPassword, err
}

func RemovePassword(name, masterPassword string) error {
	hushDir, err := getHushDir()
	if err != nil {
		return err
	}

	_, err = readEncryptedPassword(hushDir, name)
	if err != nil {
		return fmt.Errorf("failed to read password file: %w", err)
	}

	salt, err := readSalt()
	if err != nil {
		return fmt.Errorf("failed to read salt: %w", err)
	}

	_, err = validateMasterPassword(salt, masterPassword)
	if err != nil {
		return fmt.Errorf("error validating master password: %w", err)
	}

	filePath := filepath.Join(hushDir, name+".hush")
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete password file: %w", err)
	}

	return nil
}

func readEncryptedMasterPassword() (string, error) {
	hushDir, err := getHushDir()
	if err != nil {
		return "", err
	}

	masterPasswordFile := filepath.Join(hushDir, masterHashFileName)
	encryptedPassword, err := os.ReadFile(masterPasswordFile)
	if err != nil {
		return "", fmt.Errorf("failed to read encrypted master password: %w", err)
	}

	return string(encryptedPassword), nil
}

func ImplodeHush(masterPassword string) error {
	encryptedMasterPassword, err := readEncryptedMasterPassword()
	if err != nil {
		return fmt.Errorf("failed to read stored master password hash: %w", err)
	}

	salt, err := readSalt()
	if err != nil {
		return fmt.Errorf("failed to read salt: %w", err)
	}

	key, err := whisper.DeriveKeyWithSalt(masterPassword, salt)
	if err != nil {
		return fmt.Errorf("failed to derive key: %w", err)
	}

	_, err = whisper.DecryptPassword(encryptedMasterPassword, key)
	if err != nil {
		return fmt.Errorf("incorrect master password")
	}

	hushDir, err := getHushDir()
	if err != nil {
		return fmt.Errorf("failed to get hush directory: %w", err)
	}

	err = os.RemoveAll(hushDir)
	if err != nil {
		return fmt.Errorf("failed to delete hush directory: %w", err)
	}

	return nil
}

func sanitizeFileName(name string) (string, error) {
	name = strings.TrimSpace(name)

	if name == "" {
		return "", fmt.Errorf("filename cannot be empty")
	}

	if len(name) > 255 {
		return "", fmt.Errorf("filename is too long (max 255 characters)")
	}

	validChars := regexp.MustCompile(`^[a-zA-Z0-9)\-\.]+$`)
	if !validChars.MatchString(name) {
		return "", fmt.Errorf("filename contains invalid characters (only alphanumeric, underscore, hyphen, and dot are allowed)")
	}

	if strings.HasPrefix(name, ".") || strings.HasSuffix(name, ".") {
		return "", fmt.Errorf("filename cannot start or end with a dot")
	}

	return name, nil
}
