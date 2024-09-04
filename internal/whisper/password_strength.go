package whisper

import "unicode"

type PasswordStrengthError struct{ Message string }

func (e *PasswordStrengthError) Error() string {
	return e.Message
}

func CheckPasswordStrength(password string) error {
	var (
		hasMinLength = false
		hasUpper     = false
		hasLower     = false
		hasNumber    = false
	)

	if len(password) >= 6 {
		hasMinLength = true
	}

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		}
	}

	if !hasMinLength {
		return &PasswordStrengthError{"Password must be at least 6 characters long"}
	}
	if !hasUpper {
		return &PasswordStrengthError{"Password must contain at least one uppercase letter"}
	}
	if !hasLower {
		return &PasswordStrengthError{"Password must contain at least one lowercase letter"}
	}
	if !hasNumber {
		return &PasswordStrengthError{"Password must contain at least one number"}
	}

	return nil
}
