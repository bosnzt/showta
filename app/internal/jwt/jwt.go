package jwt

import (
	"github.com/golang-jwt/jwt/v4"
	"showta.cc/app/system/conf"
	"showta.cc/app/system/msg"
	"time"
)

type AppClaims struct {
	Username string `json:"username"`
	PwdStamp int64  `json:"pwd_stamp"`
	jwt.RegisteredClaims
}

func GenToken(username string, pwdStamp int64) (tokenString string, err error) {
	claim := AppClaims{
		Username: username,
		PwdStamp: pwdStamp,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(conf.AppConf.Secure.TokenExpire))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err = token.SignedString([]byte(conf.AppConf.Secure.JwtSecret))
	return tokenString, err
}

func Secret() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(conf.AppConf.Secure.JwtSecret), nil
	}
}

func ParseToken(tokenss string) (*AppClaims, error) {
	token, err := jwt.ParseWithClaims(tokenss, &AppClaims{}, Secret())
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, msg.ErrTokenInvalid
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, msg.ErrTokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, msg.ErrTokenInvalid
			} else {
				return nil, msg.ErrTokenInvalid
			}
		}
	}
	if claims, ok := token.Claims.(*AppClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, msg.ErrTokenInvalid
}
