package hash

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	// Hash the password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(password, hash string) bool {
	// Compare the password with the hash
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}