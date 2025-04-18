package jwt

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestJwtAuth_Encode(t *testing.T) {
	priKey, pubKey := loadKeyPair()
	jwtAuth := NewJwtAuth(priKey, pubKey)

	tcs := []struct {
		name         string
		customClaims jwt.MapClaims
		wantErr      error
	}{
		{
			name:         "basic",
			customClaims: jwt.MapClaims{},
			wantErr:      nil,
		}, {
			name: "with biz id",
			customClaims: jwt.MapClaims{
				BizIdParamName: float64(100000000),
			},
			wantErr: nil,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			token, err := jwtAuth.Encode(tc.customClaims)
			if err != nil {
				assert.Equal(t, tc.wantErr, err)
				return
			}

			assert.NotEmpty(t, token)

			claims, err := jwtAuth.Decode(token)
			assert.NoError(t, err)

			assert.NotEmpty(t, claims["iat"])
			assert.NotEmpty(t, claims["exp"])

			assert.Equal(t, "jotice", claims["iss"])
		})
	}
}

func TestJwtAuth_Decode(t *testing.T) {
	priKey, pubKey := loadKeyPair()
	jwtAuth := NewJwtAuth(priKey, pubKey)

	tcs := []struct {
		name      string
		tokenFunc func(t *testing.T) string
		wantErr   error
	}{
		{
			name: "validate",
			tokenFunc: func(t *testing.T) string {
				validClaims := jwt.MapClaims{
					"uid":  "100000001",
					"role": "admin",
				}

				validToken, err := jwtAuth.Encode(validClaims)
				assert.NoError(t, err)
				return validToken
			},
			wantErr: nil,
		}, {
			name: "expired",
			tokenFunc: func(t *testing.T) string {
				expiredClaims := jwt.MapClaims{
					"exp": time.Now().Add(-time.Second).Unix(),
				}
				expiredToken, err := jwtAuth.Encode(expiredClaims)
				assert.NoError(t, err)

				return expiredToken
			},
			wantErr: jwt.ErrTokenExpired,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			token := tc.tokenFunc(t)
			c, err := jwtAuth.Decode(token)

			if err != nil {
				assert.True(t, errors.Is(err, tc.wantErr))
				return
			}

			assert.NotNil(t, c)
			assert.Equal(t, "100000001", c["uid"])
			assert.Equal(t, "admin", c["role"])

		})
	}
}

var (
	priPem = `-----BEGIN PRIVATE KEY-----
MC4CAQAwBQYDK2VwBCIEIHXEAUN6Lp8Hdq8P0Mcv9mjIG1sgPWBf1Mh+OKP5HXvC
-----END PRIVATE KEY-----`
	pubPem = `-----BEGIN PUBLIC KEY-----
MCowBQYDK2VwAyEAxCSxEyY/+A7T7EtXF7AHw4Zfklh/QdjG8fxfRFYZgY8=
-----END PUBLIC KEY-----`
)

func loadKeyPair() (ed25519.PrivateKey, ed25519.PublicKey) {
	priKeyBlock, _ := pem.Decode([]byte(priPem))
	if priKeyBlock == nil {
		panic("failed to decode private key PEM")
	}
	// PEM 块本身只标注 PUBLIC KEY（通用标签），没有专门标注 ED25519 PUBLIC KEY
	// 所需要先用 x509 包统一使用 ParsePKCS8PrivateKey / ParsePKIXPublicKey 处理所有标准公钥格式
	// 再类型断言转成 ed25519.PrivateKey / ed25519.PublicKey
	priKey, err := x509.ParsePKCS8PrivateKey(priKeyBlock.Bytes)
	if err != nil {
		panic(err)
	}

	pubKeyBlock, _ := pem.Decode([]byte(pubPem))
	if pubKeyBlock == nil {
		panic("failed to decode public key PEM")
	}
	publicKey, err := x509.ParsePKIXPublicKey(pubKeyBlock.Bytes)
	if err != nil {
		panic(err)
	}

	return priKey.(ed25519.PrivateKey), publicKey.(ed25519.PublicKey)
}
