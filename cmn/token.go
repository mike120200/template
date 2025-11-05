package cmn

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func CreateToken(claims jwt.MapClaims) (string, error) {
	key := viper.GetString("safe.jwtSecret")
	if key == "" {
		appErr := NewAppError(CommonError, "jwtSecret is empty")
		Logger().Error(appErr.Error())
		return "", appErr
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	keyByte := []byte(key)
	// 使用密钥对JWT进行签名
	tokenString, err := token.SignedString(keyByte)
	if err != nil {
		appErr := NewAppError(CommonError, err.Error())
		Logger().Error("token.SignedString error: " + appErr.Error())
		return "", appErr
	}
	return tokenString, nil
}

// VerifyToken 验证jwt token
func VerifyToken(tokenString string) (jwt.Claims, error) {
	logger := zap.L()
	if logger == nil {
		fmt.Println("create logger failed, please check zap logger")
		return nil, nil
	}
	if tokenString == "" {
		err := NewAppError(CommonError, "token is empty")
		logger.Error(err.Error())
		return nil, err
	}
	key := viper.GetString("safe.jwtSecret")
	if key == "" {
		err := NewAppError(CommonError, "jwtSecret is empty")
		logger.Error(err.Error())
		return nil, err
	}
	// 解析token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.SigningMethodHS256 != token.Method {
			logger.Error("invalid signing method")
			return nil, errors.New("invalid signing method")
		}

		return []byte(key), nil
	})
	if err != nil {
		appErr := NewAppError(CommonError, err.Error())
		logger.Error("jwt.Parse error: " + appErr.Error())
		return nil, err
	}
	if !token.Valid {
		appErr := NewAppError(CommonError, "invalid tokenString"+tokenString)
		logger.Error(appErr.Error())
		return nil, appErr
	}
	// 获取token中的claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		appErr := NewAppError(CommonError, "invalid JWT claims")
		logger.Error(appErr.Error())
		return nil, appErr
	}

	// 开发期间打印解析出的openId和sessionKey

	return claims, nil

}
