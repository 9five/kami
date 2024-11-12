package domain

import "context"

type Encryption struct {
}

type EncryptionRepository interface {
}

type EncryptionUsecase interface {
	DesEncrypt(ctx context.Context, text string) (string, error)
	DesDecrypt(ctx context.Context, text string) (string, error)
	HashPassword(ctx context.Context, password string) (string, error)
	CheckPwHash(ctx context.Context, password, hash string) bool
}
