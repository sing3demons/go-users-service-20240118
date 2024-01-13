package security

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sing3demons/users/model"
)

type RegisteredClaims struct {
	jwt.RegisteredClaims
	UserName string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
}

/*
	 -> Generate key
	    mkdir -p cert
		openssl genrsa -out cert/id_rsa 4096
		openssl rsa -in cert/id_rsa -pubout -out cert/id_rsa.pub
*/
func GenerateToken(user model.User) (token string, err error) {
	privateKey, err := getSecretPrivateKeyFromEnv()
	if err != nil {
		return "", err
	}

	rsa, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return "", err
	}

	claims := &RegisteredClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID.Hex(),
			Issuer:    os.Getenv("ISSUER"),
			Audience:  jwt.ClaimStrings{},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
	}

	if os.Getenv("AUDIENCE") != "" {
		audience := strings.Split(os.Getenv("AUDIENCE"), ",")
		claims.Audience = append(claims.Audience, audience...)
	}

	if user.Email != "" {
		claims.Email = user.Email
	}

	if user.Username != "" {
		claims.UserName = user.Username
	}

	return jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(rsa)
}

func ValidateToken(token string) (jwt.MapClaims, error) {
	publicKey, err := getSecretPublicKeyFromEnv()
	if err != nil {
		return nil, err
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		return nil, fmt.Errorf("validate: parse key: %w", err)
	}

	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return nil, fmt.Errorf("validate: invalid")
	}

	return claims, nil
}
func getSecretPrivateKeyFromEnv() (privateKey []byte, err error) {
	private := os.Getenv("PRIVATE_KEY")
	if private == "" {
		b, err := os.ReadFile("cert/id_rsa")
		if err != nil {
			return nil, err
		}
		private = base64.StdEncoding.EncodeToString(b)
	}
	privateKey, err = base64.StdEncoding.DecodeString(private)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func getSecretPublicKeyFromEnv() (publicKey []byte, err error) {
	public := os.Getenv("PUBLIC_KEY")
	if public == "" {
		b, err := os.ReadFile("cert/id_rsa.pub")
		if err != nil {
			return nil, err
		}
		public = base64.StdEncoding.EncodeToString(b)
	}
	publicKey, err = base64.StdEncoding.DecodeString(public)
	if err != nil {
		return nil, err
	}
	return publicKey, nil
}
