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
	// 設置過期時間（24小時）
	expirationTime := time.Now().Add(24 * time.Hour)

	// 創建聲明（轉換為 MapClaims 以便使用 CreateToken）
	claims := jwt.MapClaims{
		"user_id":  userId,
		"username": username,
		"exp":      expirationTime.Unix(),
		"iat":      time.Now().Unix(),
		"iss":      "my_template",
	}

	// 使用 CreateToken 生成 token
	return CreateToken(claims)
}

// ParseToken 解析 JWT token
func ParseToken(tokenString string) (*Claims, error) {
	// 使用 VerifyToken 驗證並解析 token
	claims, err := VerifyToken(tokenString)
	if err != nil {
		// 檢查是否為過期錯誤
		var ve *jwt.ValidationError
		if errors.As(err, &ve) && ve.Errors&jwt.ValidationErrorExpired != 0 {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	// 將 MapClaims 轉換為 Claims 結構體
	mapClaims, ok := claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	// 檢查 token 是否過期（通過 exp 字段）
	if exp, ok := mapClaims["exp"].(float64); ok {
		if int64(exp) < time.Now().Unix() {
			return nil, ErrExpiredToken
		}
	}

	result := &Claims{}
	if userId, ok := mapClaims["user_id"].(string); ok {
		result.UserId = userId
	}
	if username, ok := mapClaims["username"].(string); ok {
		result.Username = username
	}
	if exp, ok := mapClaims["exp"].(float64); ok {
		result.ExpiresAt = int64(exp)
	}
	if iat, ok := mapClaims["iat"].(float64); ok {
		result.IssuedAt = int64(iat)
	}
	if iss, ok := mapClaims["iss"].(string); ok {
		result.Issuer = iss
	}

	return result, nil
}
