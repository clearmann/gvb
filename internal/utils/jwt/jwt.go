package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	g "gvb/internal/global"
	"time"
)

var (
	ErrTokenExpired     = errors.New("token 已过期，请重新登录")
	ErrTokenNotValidYet = errors.New("token 无效，请重新登录")
	ErrTokenMalformed   = errors.New("token 不正确，请重新登录")
	ErrTokenNotValid    = errors.New("这不是一个 token ，请重新登录")
)

type MyClaim struct {
	UserId   int    `json:"user_id"`
	Username string `json:"username"`
	RoleId   []int  `json:"role_id"`
	jwt.RegisteredClaims
}

func GenToken(userId int, roleId []int) (string, error) {
	secret := []byte(g.Conf.JWT.Secret)
	expireHour := g.Conf.JWT.Expire
	claim := &MyClaim{
		UserId: userId,
		RoleId: roleId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    g.Conf.JWT.Issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireHour) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString(secret)
}
func ParseToken(tokenString string) (*MyClaim, error) {
	secret := []byte(g.Conf.JWT.Secret)
	jwtToken, err := jwt.ParseWithClaims(
		tokenString, &MyClaim{}, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})
	if err != nil {
		switch ve, ok := err.(*jwt.ValidationError); ok {
		case ve.Errors&jwt.ValidationErrorMalformed != 0:
			return nil, ErrTokenMalformed
		case ve.Errors&jwt.ValidationErrorExpired != 0:
			return nil, ErrTokenExpired
		case ve.Errors&jwt.ValidationErrorNotValidYet != 0:
			return nil, ErrTokenNotValidYet
		}
	}
	if claims, ok := jwtToken.Claims.(*MyClaim); ok && jwtToken.Valid {
		return claims, nil
	}
	return nil, ErrTokenNotValid
}
