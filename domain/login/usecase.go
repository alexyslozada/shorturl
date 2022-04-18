package login

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"

	"github.com/alexyslozada/shorturl/model"
)

type Login struct {
	storage Storage
}

func New(s Storage) Login {
	return Login{storage: s}
}

func (l Login) Login(email, password string) (string, error) {
	user, err := l.storage.ByEmail(email)
	if err != nil {
		return "", fmt.Errorf("login.Login() %w", err)
	}
	if !user.Active {
		return "", fmt.Errorf("the user is not active")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", fmt.Errorf("password is wrong")
	}

	claims := model.JWTCustomClaims{
		UserID: user.ID,
		Email:  user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    "Alexys Lozada",
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(model.Secret))
	if err != nil {
		return "", fmt.Errorf("login.Login() %w", err)
	}

	return t, nil
}
