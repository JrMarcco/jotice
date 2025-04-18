package jwt

import (
	"context"
	"crypto/ed25519"
	"errors"
	"fmt"
	"strings"
	"time"

	"maps"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const BizIdParamName = "biz_id"

type BizIdKey = struct{}

type InterceptorBuilder struct {
	priKey ed25519.PrivateKey
	pubKey ed25519.PublicKey
}

// NewJwtAuth create a new jwt auth interceptor
func NewJwtAuth(privateKey ed25519.PrivateKey, publicKey ed25519.PublicKey) *InterceptorBuilder {
	return &InterceptorBuilder{
		priKey: privateKey,
		pubKey: publicKey,
	}
}

// Encode encode a custom claim to a jwt token
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

// Decode decode a jwt token to a custom claim
func (b *InterceptorBuilder) Decode(tokenStr string) (jwt.MapClaims, error) {
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("unsupport sign algorithm: %v", t.Header["alg"])
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

// Build crate a grpc interceptor for jwt auth
func (b *InterceptorBuilder) Build() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		authHeaders := md.Get("Authorization")
		if len(authHeaders) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization token")
		}
		tokenStr := authHeaders[0]

		mc, err := b.Decode(tokenStr)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				return nil, status.Error(codes.Unauthenticated, "token expired")
			}
			if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
				return nil, status.Error(codes.Unauthenticated, "invalid signature")
			}
			return nil, status.Errorf(codes.Unauthenticated, "invalid token: %s", err.Error())
		}

		if val, ok := mc[BizIdParamName]; ok {
			bizId := int64(val.(float64))
			ctx = context.WithValue(ctx, BizIdKey{}, bizId)
		}
		return handler(ctx, req)
	}
}
