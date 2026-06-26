package jwt

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/duke-git/lancet/v2/condition"
	jwt "github.com/golang-jwt/jwt/v5"
	"strings"
	"time"
)

var jwtSecret = []byte("bcm-9c0b56e0-0f8a-43ab-856b-9bc5f37e2e7e")

type Claims struct {
	Username string                 `json:"username"`
	UserId   uint64                 `json:"user_id"`
	AppId    string                 `json:"app_id"`
	Level    string                 `json:"level"`  // 是否是超级用户 1 表示 是，0 表示 不是
	Extra    map[string]interface{} `json:"extra"`  // 任意附加信息，一般用于存储权限等数据
	Access   bool                   `json:"access"` // 该用户是否可以访问系统，该字段仅用于前台访问控制，后端程序只赋值，不做处理
	jwt.MapClaims
}

func (c *Claims) VipIsExpired() bool {
	if val, ok := c.Extra["vipIsExpired"]; ok {
		if v, ok2 := val.(bool); ok2 {
			return v
		}
		return false
	}
	return true
}

func GeneratePassword(s, salt string) string {
	b := []byte(s)
	h := md5.New()
	h.Write(b)
	h.Write([]byte(salt))
	return hex.EncodeToString(h.Sum(nil))
}

func GenerateToken(c Claims) (string, error) {
	nowTime := time.Now()                     //当前时间
	expireTime := nowTime.Add(24 * time.Hour) //有效时间

	claims := Claims{
		Username: c.Username,
		UserId:   c.UserId,
		AppId:    c.AppId,
		Access:   c.Access,
		Level:    c.Level,
		Extra:    c.Extra,
		MapClaims: jwt.MapClaims{
			"exp":    expireTime.Unix(),
			"Issuer": "YiGeChengZi",
			"Extra":  c.Extra,
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

// GenerateExpiredToken 生成一个过期的token
func GenerateExpiredToken(c Claims) (string, error) {
	claims := Claims{
		Username: c.Username,
		UserId:   c.UserId,
		AppId:    c.AppId,
		MapClaims: jwt.MapClaims{
			"exp":    time.Now().Unix(),
			"Issuer": "YiGeChengZi",
			"Extra":  map[string]interface{}{},
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

func ParseToken(token string) (*Claims, error) {
	// token 以Bearer开头需要截取掉
	_token := condition.TernaryOperator(strings.HasPrefix(token, "Bearer"), strings.Replace(token, "Bearer ", "", 1), token)
	tokenClaims, _ := jwt.ParseWithClaims(_token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, errors.New("登录已经失效")
}
