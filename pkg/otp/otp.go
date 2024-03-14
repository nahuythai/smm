package otp

import (
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

var global *service

var defaultOption = Option{
	Issuer: "SMM",
}

type Option struct {
	Issuer string
	Digits int
	Period int
}

type service struct {
	issuer string
	digits int
	period int
}
type Service interface {
	InitGlobal()
}

func New(opts ...Option) Service {
	opt := defaultOption
	if len(opts) > 0 {
		if opts[0].Digits > 0 {
			opt.Digits = opts[0].Digits
			opt.Period = opts[0].Period
		}
	}
	return &service{
		issuer: opt.Issuer,
		digits: opt.Digits,
		period: opt.Period,
	}
}

func (s *service) InitGlobal() {
	global = s
}

func GetGlobal() *service {
	return global
}

func (s *service) GeneratePassCode(secretKey string) (string, error) {
	passcode, err := totp.GenerateCodeCustom(secretKey, time.Now(), totp.ValidateOpts{
		Period:    60,
		Skew:      1,
		Digits:    otp.Digits(s.digits),
		Algorithm: otp.AlgorithmSHA1,
	})
	return passcode, err
}

func (s *service) Validate(passcode, secretKey string) bool {
	return totp.Validate(passcode, secretKey)
}

func (s *service) GenerateSecretKey(username string) (*otp.Key, error) {
	return totp.Generate(totp.GenerateOpts{
		Issuer:      s.issuer,
		AccountName: username,
		Period:      60,
		Digits:      otp.Digits(s.digits),
		Algorithm:   otp.AlgorithmSHA1,
	})
}
