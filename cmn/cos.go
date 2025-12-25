package cmn

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/spf13/viper"
	"github.com/tencentyun/cos-go-sdk-v5"
)

var cosInstance *cos.Client

func InitCos() {
	bucketName := viper.GetString("cos.bucket_name")
	if bucketName == "" {
		panic("bucket_name can not be empty")
	}
	region := viper.GetString("cos.region")
	if region == "" {
		panic("region can not be empty")
	}
	CosSecretID := viper.GetString("cos.secret_id")
	if CosSecretID == "" {
		panic("secret_id can not be empty")
	}
	CosSecretKey := viper.GetString("cos.secret_key")
	if CosSecretKey == "" {
		panic("secret_key can not be empty")
	}
	CosTimeOut := viper.GetInt("cos.timeout")
	if CosTimeOut <= 30 {
		CosTimeOut = 30
	}
	CosTimeOutDuration := time.Duration(CosTimeOut) * time.Second
	u, _ := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", bucketName, region))
	fmt.Printf("cos url: %s\n", u.String())
	cosInstance = cos.NewClient(
		&cos.BaseURL{BucketURL: u},
		&http.Client{
			Timeout: CosTimeOutDuration,
			Transport: &cos.AuthorizationTransport{
				SecretID:  CosSecretID,
				SecretKey: CosSecretKey,
				Transport: &http.Transport{
					MaxIdleConns:        100,
					MaxIdleConnsPerHost: 100,
					IdleConnTimeout:     90 * time.Second,
				},
			},
		},
	)
}

func Cos() *cos.Client {
	if cosInstance == nil {
		panic("cos instance is nil")
	}
	return cosInstance
}
