package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func CheckPassword(
    hashedPassword string,
    password string,
) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}