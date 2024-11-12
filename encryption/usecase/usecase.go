package usecase

import (
	"bytes"
	"context"
	"crypto/des"
	"encoding/hex"
	"errors"
	"kami/domain"

	"golang.org/x/crypto/bcrypt"
)

type encryptionUsecase struct {
	desKey []byte
}

func NewEncryptionUsecase(desKey []byte) domain.EncryptionUsecase {
	return &encryptionUsecase{
		desKey: desKey,
	}
}

func (e *encryptionUsecase) DesEncrypt(ctx context.Context, text string) (string, error) {
	block, err := des.NewCipher(e.desKey)
	if err != nil {
		return "", nil
	}
	bs := block.BlockSize()
	src := zeroPadding([]byte(text), bs)
	if len(src)%bs != 0 {
		return "", errors.New("Need a multiple of the blocksize")
	}
	out := make([]byte, len(src))
	dst := out
	for len(src) > 0 {
		block.Encrypt(dst, src[:bs])
		src = src[bs:]
		dst = dst[bs:]
	}
	return hex.EncodeToString(out), nil
}

func (e *encryptionUsecase) DesDecrypt(ctx context.Context, text string) (string, error) {
	src, err := hex.DecodeString(text)
	if err != nil {
		return "", err
	}
	block, err := des.NewCipher(e.desKey)
	if err != nil {
		return "", err
	}
	out := make([]byte, len(src))
	dst := out
	bs := block.BlockSize()
	if len(src)%bs != 0 {
		return "", errors.New("crypto/cipher: input not full blocks")
	}
	for len(src) > 0 {
		block.Decrypt(dst, src[:bs])
		src = src[bs:]
		dst = dst[bs:]
	}
	out = zeroUnPadding(out)
	return string(out), nil
}

func zeroPadding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(cipherText, padtext...)
}

func zeroUnPadding(originData []byte) []byte {
	return bytes.TrimFunc(originData,
		func(r rune) bool {
			return r == rune(0)
		})
}

func (e *encryptionUsecase) HashPassword(ctx context.Context, password string) (string, error) {
	result, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(result), err
}

func (e *encryptionUsecase) CheckPwHash(ctx context.Context, password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
