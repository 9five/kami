package domain

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type KamiUser struct {
	gorm.Model
	Email    string    `json:"email"`
	Phone    string    `json:"phone"`
	Password string    `json:"password"`
	Status   string    `json:"status"`
	Points   int64     `json:"points"`
	Gender   string    `json:"gender"`
	Birthday time.Time `json:"birthday"`
	Name     string    `json:"name"`
	Career   string    `json:"career"`
}

type KamiUserRepository interface {
	New(ctx context.Context, user *KamiUser) (*KamiUser, error)
	Get(ctx context.Context, user *KamiUser, optsWhere ...map[string]interface{}) (*KamiUser, error)
	Update(ctx context.Context, user *KamiUser) error
	NewLog(ctx context.Context, log *KamiUserLog) (*KamiUserLog, error)
	GetLog(ctx context.Context, log *KamiUserLog) (*KamiUserLog, error)
	UpdateLog(ctx context.Context, log *KamiUserLog) error
}

type KamiUserUsercase interface {
	NewKamiUser(ctx context.Context, user *KamiUser) (*KamiUser, error)
	GetKamiUser(ctx context.Context, user *KamiUser) (*KamiUser, error)
	UpdateKamiUser(ctx context.Context, user *KamiUser) error
	GetKamiUserLog(ctx context.Context, log *KamiUserLog) (*KamiUserLog, error)
	UpdateKamiUserLog(ctx context.Context, log *KamiUserLog) error
	LoginKamiUser(ctx context.Context, input *KamiUser) (*KamiUser, error)
	CheckKamiUserLog(ctx context.Context, log *KamiUserLog) error
	GenerateToken(ctx context.Context, user *KamiUser) (string, error)
	UpdateUserInfo(ctx context.Context, user *KamiUser, userInput *KamiUserInput) error
}

type KamiUserOutput struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Status   string `json:"status"`
	Points   int64  `json:"points"`
	Gender   string `json:"gender"`
	Birthday string `json:"birthday"`
	Name     string `json:"name"`
	Career   string `json:"career"`
}

type KamiUserLog struct {
	gorm.Model
	Phone         string    `json:"phone"`
	AuthTime      time.Time `json:"auth_time"`
	AuthFrequency int64     `json:"auth_frequency"`
}

type VerificationInput struct {
	ForgotPw bool   `json:"forgot_pw"`
	Token    string `json:"token"`
	Code     string `json:"code"`
}

type KamiUserInput struct {
	Email    string `json:"email"`
	Gender   string `json:"gender"`
	Birthday string `json:"birthday"`
	Name     string `json:"name"`
	Career   string `json:"career"`
}
