package jwt

import (
	"crypto/ed25519"
	"fmt"
	"strings"
	"time"

	"maps"

	"github.com/golang-jwt/jwt/v5"
)

const BIZ_ID_PARAM_NAME = "biz_id"

type InterceptorBuilder struct {
	priKey ed25519.PrivateKey
	pubKey ed25519.PublicKey
}

func NewJwtAuth(privateKey ed25519.PrivateKey, publicKey ed25519.PublicKey) *InterceptorBuilder {
	return &InterceptorBuilder{
		priKey: privateKey,
		pubKey: publicKey,
	}
}

func (b *InterceptorBuilder) Encode(customClaims jwt.MapClaims) (string, error) {
	claims := jwt.MapClaims{
		"iat": time.Now().Unix(),
		"iss": "jotice",
	}

	maps.Copy(claims, customClaims)

	if _, ok := claims["exp"]; !ok {
		claims["exp"] = time.Now().Add(24 * time.Hour).Unix()
	}

	token := jwt.NewWithClaims(&jwt.SigningMethodEd25519{}, claims)
	return token.SignedString(b.priKey)
}

func (b *InterceptorBuilder) Decode(tokenStr string) (jwt.MapClaims, error) {
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("unsupport sign method: %v", t.Header["alg"])
		}

		return b.pubKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("fail to decode token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
