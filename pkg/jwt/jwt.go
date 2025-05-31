package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	Jit string `json:"jti,omitempty"` // JWT ID
	Sub string `json:"sub,omitempty"` // Subject (user ID)
	Iat int64  `json:"iat,omitempty"` // Issued at
	Exp int64  `json:"exp,omitempty"` // Expiration time
}

func GenerateHS256JWT(claims JWTClaims, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"jti": claims.Jit,
		"sub": claims.Sub,
		"iat": claims.Iat,
		"exp": claims.Exp,
	})

	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ParseHS256JWT(tokenString string, secret string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func ValidateHS256JWT(token *jwt.Token) (bool, error) {
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if _, ok := claims["jti"]; !ok {
			return false, fmt.Errorf("%w: missing jti claim", jwt.ErrTokenRequiredClaimMissing)
		}
		if _, ok := claims["sub"]; !ok {
			return false, fmt.Errorf("%w: missing sub claim", jwt.ErrTokenRequiredClaimMissing)
		}
		if _, ok := claims["iat"]; !ok {
			return false, fmt.Errorf("%w: missing iat claim", jwt.ErrTokenRequiredClaimMissing)
		}
		if _, ok := claims["exp"]; !ok {
			return false, fmt.Errorf("%w: missing exp claim", jwt.ErrTokenRequiredClaimMissing)
		}
		return true, nil
	}
	return false, jwt.ErrTokenInvalidClaims
}

func ExtractClaims(token *jwt.Token) (JWTClaims, error) {
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		jwtClaims := JWTClaims{
			Jit: claims["jti"].(string),
			Sub: claims["sub"].(string),
			Iat: int64(claims["iat"].(float64)),
			Exp: int64(claims["exp"].(float64)),
		}
		return jwtClaims, nil
	}
	return JWTClaims{}, jwt.ErrTokenInvalidClaims
}
