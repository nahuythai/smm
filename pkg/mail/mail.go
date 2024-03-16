package mail

import (
	"github.com/go-mail/mail"
)

var global *service

var defaultOption = Option{
	Email:    "noreply@smm.com",
	Password: "",
	MailHost: "localhost",
	MailPort: 1025, //nolint:gomnd
}

type Option struct {
	Email    string
	Password string
	MailHost string
	MailPort int
}

type service struct {
	opt    Option
	dialer *mail.Dialer
}
type Service interface {
	InitGlobal()
	SendHTMLMail(receiver string, body string) error
}

func New(opts ...Option) Service {
	s := &service{
		opt: defaultOption,
	}
	if len(opts) > 0 {
		s.opt = opts[0]
	}
	s.dialer = mail.NewDialer(s.opt.MailHost, s.opt.MailPort, s.opt.Email, s.opt.Password)
	return s
}

func (s *service) InitGlobal() {
	global = s
}

func GetGlobal() *service {
	return global
}

func (s *service) SendHTMLMail(receiver string, body string) error {
	m := mail.NewMessage()

	m.SetHeader("From", s.opt.Email)

	m.SetHeader("To", receiver)

	m.SetHeader("Subject", "Email Verification")

	m.SetBody("text/html", body)
	if err := s.dialer.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
