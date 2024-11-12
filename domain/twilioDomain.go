package domain

import ()

type TwilioService struct {
	TwilioSID                 string
	TwilioToken               string
	TwilioFromPhone           string
	TwilioVerifyServiceSID    string
	TwilioVerifyToCountryCode string
	TwilioVerifyChannel       string
}

type TwilioServiceRepository interface {
	GetTwilioConfig() (*TwilioService, error)
}

type TwilioServiceUsecase interface {
	SendVerificationSMS(phone string) error
	VerificationCheck(phone string, code string) error
}
