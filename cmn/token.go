package cmn

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Claims JWT 聲明結構
type Claims struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// 定義錯誤
var (
	ErrInvalidToken = errors.New("无效token")
	ErrExpiredToken = errors.New("token已过期")
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

// GenerateToken 生成 JWT token
func GenerateToken(userId, username string) (string, error) {
	key := viper.GetString("safe.jwtSecret")
	if key == "" {
		// 使用默認密鑰（僅用於開發環境）
		key = "default-secret-key-change-in-production"
	}

	// 設置過期時間（24小時）
	expirationTime := time.Now().Add(24 * time.Hour)

	// 創建聲明
	claims := &Claims{
		UserId:   userId,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "my_template",
		},
	}

	// 創建 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 簽名 token
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		Logger().Error("生成token失敗", zap.Error(err))
		return "", err
	}

	return tokenString, nil
}

// ParseToken 解析 JWT token
func ParseToken(tokenString string) (*Claims, error) {
	key := viper.GetString("safe.jwtSecret")
	if key == "" {
		// 使用默認密鑰（僅用於開發環境）
		key = "default-secret-key-change-in-production"
	}

	// 解析 token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 驗證簽名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(key), nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, ErrExpiredToken
			}
		}
		Logger().Error("解析token失敗", zap.Error(err))
		return nil, ErrInvalidToken
	}

	// 獲取聲明
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}
