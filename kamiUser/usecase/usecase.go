package usecase

import (
	"context"
	"errors"
	"fmt"
	"kami/domain"
	"time"
	_ "time/tzdata"

	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

type kamiUserUsecase struct {
	jwtKey       []byte
	desKey       []byte
	kamiUserRepo domain.KamiUserRepository
}

func NewKamiUserUsecase(jwtKey []byte, kamiUserRepo domain.KamiUserRepository) domain.KamiUserUsercase {
	return &kamiUserUsecase{
		jwtKey:       jwtKey,
		kamiUserRepo: kamiUserRepo,
	}
}

func (k *kamiUserUsecase) NewKamiUser(ctx context.Context, user *domain.KamiUser) (result *domain.KamiUser, err error) {
	return k.kamiUserRepo.New(ctx, user)
}

func (k *kamiUserUsecase) GetKamiUser(ctx context.Context, user *domain.KamiUser) (result *domain.KamiUser, err error) {
	return k.kamiUserRepo.Get(ctx, user)
}

func (k *kamiUserUsecase) UpdateKamiUser(ctx context.Context, user *domain.KamiUser) error {
	return k.kamiUserRepo.Update(ctx, user)
}

func (k *kamiUserUsecase) GetKamiUserLog(ctx context.Context, log *domain.KamiUserLog) (*domain.KamiUserLog, error) {
	return k.kamiUserRepo.GetLog(ctx, log)
}

func (k *kamiUserUsecase) UpdateKamiUserLog(ctx context.Context, log *domain.KamiUserLog) error {
	return k.kamiUserRepo.UpdateLog(ctx, log)
}

func (k *kamiUserUsecase) LoginKamiUser(ctx context.Context, input *domain.KamiUser) (*domain.KamiUser, error) {
	user, err := k.kamiUserRepo.Get(ctx, &domain.KamiUser{Phone: input.Phone})
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if user.Model.ID == 0 {
		user, err = k.kamiUserRepo.New(ctx, input)
		if err != nil {
			return nil, err
		}
	} else if user.Password == "" {
		user.Password = input.Password
		if err = k.kamiUserRepo.Update(ctx, user); err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (k *kamiUserUsecase) CheckKamiUserLog(ctx context.Context, log *domain.KamiUserLog) error {
	logRecord, err := k.kamiUserRepo.GetLog(ctx, log)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	loc, _ := time.LoadLocation("Asia/Taipei")
	if logRecord.Model.ID != 0 {
		if logRecord.AuthFrequency >= 3 && time.Now().In(loc).Sub(logRecord.AuthTime).Hours() < 24 {
			return errors.New("the login activity is abnormal, please try again in 24 hours or contact KAMIKAMI customer service")
		}

		if time.Now().In(loc).Sub(logRecord.AuthTime).Seconds() < 60 {
			return errors.New("the sending time of the two verification codes must exceed 60 seconds")
		}
	} else {
		logRecord.Phone = log.Phone
		if logRecord, err = k.kamiUserRepo.NewLog(ctx, logRecord); err != nil {
			return err
		}
	}

	logRecord.AuthTime = time.Now().In(loc)
	logRecord.AuthFrequency += 1
	if err := k.kamiUserRepo.UpdateLog(ctx, logRecord); err != nil {
		return err
	}

	return nil
}

func (k *kamiUserUsecase) GenerateToken(ctx context.Context, user *domain.KamiUser) (string, error) {
	expiresAt := time.Now().AddDate(0, 0, 7).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, domain.PrivateClaims{
		StandardClaims: jwt.StandardClaims{
			Subject:   fmt.Sprintf("%s", user.Phone),
			ExpiresAt: expiresAt,
		},
		UID: user.Model.ID,
	})
	tokenString, err := token.SignedString(k.jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (k *kamiUserUsecase) UpdateUserInfo(ctx context.Context, user *domain.KamiUser, userInput *domain.KamiUserInput) error {
	if user.Birthday.IsZero() {
		bd, err := time.Parse("2006-01-02", userInput.Birthday)
		if err != nil {
			return err
		}
		user.Birthday = bd
	}

	if user.Email == "" {
		user.Email = userInput.Email
	}

	user.Gender = userInput.Gender
	user.Name = userInput.Name
	user.Career = userInput.Career

	return k.kamiUserRepo.Update(ctx, user)
}
