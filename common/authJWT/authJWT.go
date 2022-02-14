package authJWT

import (
	// "strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/lemon997/lemonMall/global"
)

type Claims struct {
	Version    int8   `json:"version"`
	CustomerId int64  `json:"customerid"`
	LoginName  string `json:"loginname"`
	//playload存储ID,username,version,version用于区别是否有更改密码
	jwt.StandardClaims
}

func GetJWTSecret() []byte {
	return []byte(global.JWTSetting.Secret)
}

func GenerateToken(name string, id int64, version int8) (string, error) {
	//生成token
	nowTime := time.Now()
	expireTime := nowTime.Add(global.JWTSetting.Expire)
	claims := Claims{
		LoginName:  name,
		CustomerId: id,
		Version:    version,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    global.JWTSetting.Issuer,
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(GetJWTSecret())
	return token, err
}

func ParseToken(tokenString string) (*Claims, error) {
	//解析Token
	tokenClaims, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return GetJWTSecret(), nil
	})
	if err != nil {
		return nil, err
	}
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
