//go:build awsDev
// +build awsDev

package config

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	DB            *gorm.DB
	JwtKey        []byte
	DesKey        []byte
	AllowOrigin   []string
	AwsBucket     string
	WeibyVendorId string
	WeibyApiKey   string
)

func init() {
	DB = newDBConfig()
	JwtKey = []byte("zesnSXqjiDLctJjW5czrwRFFNzknN+v4mIR6pCmXF+kUkMyBXyATDeqreCeV+C8jtFkAC3XIAQ4DC5aIqmjPj1WDd0FJGZDU0D686mkKPR6SlXhZkKuilfI0mGQ+cg1f1rqCPOwHjWCmSLVHvHBvac0Y+3p8j9gi1KLsdX5Q2/BpBhbUXZGP9TUmwwHZ2Li1E5dA3jB0EWQO5HfGEazKjJSJ/5IQl+wkKKagPQ==")
	DesKey = []byte("w!z%C*F-")
	AllowOrigin = []string{
		"https://membership.kamikami.co",
		"https://membership.kamikami.co:3001",
		"http://localhost:3000",
		"http://localhost:8000",
	}
	AwsBucket, _ = os.LookupEnv("AWS_BUCKET")
	WeibyVendorId = "5d30ab"
	WeibyApiKey = "7AnzZ1jKqG16vXalfX4qS7pm6BqGr1fBjtRsS3Ld"
}

func newDBConfig() *gorm.DB {
	dsn := "host=172.31.21.139 user=kamiDev password=capsulekami dbname=kamiDev port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{TablePrefix: "kami."},
	})
	if err != nil {
		panic(err.Error())
	}
	return db
}

func GetS3Client(ctx context.Context) (s3Client *s3.Client, err error) {
	awsID, _ := os.LookupEnv("AWS_ID")
	awsSecret, _ := os.LookupEnv("AWS_SECRET")
	awsRegion, _ := os.LookupEnv("AWS_REGION")

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(awsRegion))
	if err != nil {
		return
	}

	cfg.Credentials = aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(awsID, awsSecret, ""))

	s3Client = s3.NewFromConfig(cfg)
	return
}
