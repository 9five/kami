package usecase

import (
	"errors"
	"fmt"
	"kami/domain"

	"github.com/twilio/twilio-go"
	twilioVerifyV2 "github.com/twilio/twilio-go/rest/verify/v2"
)

type twilioServiceUsecase struct {
	twilioRepo domain.TwilioServiceRepository
}

func NewTwilioServiceUsecase(twilioRepo domain.TwilioServiceRepository) domain.TwilioServiceUsecase {
	return &twilioServiceUsecase{
		twilioRepo: twilioRepo,
	}
}

func (t *twilioServiceUsecase) SendVerificationSMS(phone string) error {
	twilioConfig, err := t.twilioRepo.GetTwilioConfig()
	if err != nil {
		return err
	}

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: twilioConfig.TwilioSID,
		Password: twilioConfig.TwilioToken,
	})

	params := &twilioVerifyV2.CreateVerificationParams{}
	params.SetTo(fmt.Sprintf(twilioConfig.TwilioVerifyToCountryCode+"%s", phone))
	params.SetChannel(twilioConfig.TwilioVerifyChannel)

	_, err = client.VerifyV2.CreateVerification(twilioConfig.TwilioVerifyServiceSID, params)
	return err
}

func (t *twilioServiceUsecase) VerificationCheck(phone string, code string) error {
	twilioConfig, err := t.twilioRepo.GetTwilioConfig()
	if err != nil {
		return err
	}

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: twilioConfig.TwilioSID,
		Password: twilioConfig.TwilioToken,
	})

	params := &twilioVerifyV2.CreateVerificationCheckParams{}
	params.SetCode(code)
	params.SetTo(fmt.Sprintf(twilioConfig.TwilioVerifyToCountryCode+"%s", phone))

	resp, err := client.VerifyV2.CreateVerificationCheck(twilioConfig.TwilioVerifyServiceSID, params)
	if err != nil {
		return err
	}

	if fmt.Sprint(resp.Valid) == "false" {
		return errors.New("verification code incorrect")
	}
	return nil
}
