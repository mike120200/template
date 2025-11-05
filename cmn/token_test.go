package cmn

import (
	"github.com/spf13/viper"
	_ "reflect"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func TestCreateToken(t *testing.T) {
	// 设置一个临时的jwt secret
	viper.Set("safe.jwtSecret", "test-secret")

	claims := jwt.MapClaims{
		"user_id": 123,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token, err := CreateToken(claims)
	if err != nil {
		t.Fatalf("CreateToken() error = %v", err)
	}

	if token == "" {
		t.Errorf("CreateToken() got empty token")
	}
	t.Logf("token: %s", token)
}

func TestVerifyToken(t *testing.T) {
	viper.Set("safe.jwtSecret", "test-secret")
	userId := 123
	claims := jwt.MapClaims{
		"user_id": float64(userId), // 注意：通过JWT传输时，数字可能会被解析为float64
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token, err := CreateToken(claims)
	if err != nil {
		t.Fatalf("Failed to create token for testing: %v", err)
	}

	// 正常验证
	parsedClaims, err := VerifyToken(token)
	if err != nil {
		t.Fatalf("VerifyToken() with valid token error = %v", err)
	}

	if parsedClaims == nil {
		t.Fatal("VerifyToken() with valid token returned nil claims")
	}
	if gotUserId := parsedClaims.(jwt.MapClaims)["user_id"]; gotUserId != float64(userId) {
		t.Errorf("VerifyToken() got user_id = %v, want %v", gotUserId, float64(userId))
	}

	// 无效token
	_, err = VerifyToken("invalid-token")
	if err == nil {
		t.Error("VerifyToken() with invalid token should return an error")
	}

	// 空token
	_, err = VerifyToken("")
	if err == nil {
		t.Error("VerifyToken() with empty token should return an error")
	}

	// 使用错误的secret
	viper.Set("safe.jwtSecret", "wrong-secret")
	_, err = VerifyToken(token)
	if err == nil {
		t.Error("VerifyToken() with wrong secret should return an error")
	}
	// 恢复secret
	viper.Set("safe.jwtSecret", "test-secret")

	// 测试过期的token
	expiredClaims := jwt.MapClaims{
		"user_id": 456,
		"exp":     time.Now().Add(-time.Hour).Unix(),
	}
	expiredToken, _ := CreateToken(expiredClaims)
	_, err = VerifyToken(expiredToken)
	if err == nil {
		t.Error("VerifyToken() with expired token should return an error")
	}
}

func TestToken_NoSecret(t *testing.T) {
	_ = ViperInit(".conf_linux.json")
	_ = LoggerInit()
	viper.Set("safe.jwtSecret", "")

	// CreateToken with no secret
	_, err := CreateToken(jwt.MapClaims{"user_id": 1})
	if err == nil {
		t.Error("CreateToken() with no secret should return an error")
	}

	// VerifyToken with no secret
	_, err = VerifyToken("any-token")
	if err == nil {
		t.Error("VerifyToken() with no secret should return an error")
	}
}
