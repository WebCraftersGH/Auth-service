package usecase

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/WebCraftersGH/Auth-service/internal/contracts"
	"github.com/WebCraftersGH/Auth-service/internal/domain"
)

type mailSVC struct {
	mailsRepo contracts.MailsRepo
	smtpHost  string
	smtpPort  int
	smtpUser  string
	smtpPass  string
	smtpFrom  string
	logger    contracts.ILogger
}

func NewMailSVC(
	mailsRepo contracts.MailsRepo,
	smtpHost string,
	smtpPort int,
	smtpUser string,
	smtpPass string,
	smtpFrom string,
	logger contracts.ILogger,
) *mailSVC {
	return &mailSVC{
		mailsRepo: mailsRepo,
		smtpHost:  smtpHost,
		smtpPort:  smtpPort,
		smtpUser:  smtpUser,
		smtpPass:  smtpPass,
		smtpFrom:  smtpFrom,
		logger:    logger,
	}
}

func (s *mailSVC) SendMail(ctx context.Context, toEmail string, code domain.OTP) error {
	subject := "Your OTP code"
	body := fmt.Sprintf("Your OTP code is: %s\nIt will expire at %s\n", code.Value, code.ExpiresAt.Format("2006-01-02 15:04:05 MST"))
	raw := buildRawMail(s.smtpFrom, toEmail, subject, body)

	addr := fmt.Sprintf("%s:%d", s.smtpHost, s.smtpPort)
	var auth smtp.Auth
	if strings.TrimSpace(s.smtpUser) != "" {
		auth = smtp.PlainAuth("", s.smtpUser, s.smtpPass, s.smtpHost)
	}

	if err := smtp.SendMail(addr, auth, s.smtpFrom, []string{toEmail}, []byte(raw)); err != nil {
		return err
	}

	if err := s.mailsRepo.SaveMail(ctx, domain.Mail{Value: raw, ToEmail: toEmail}); err != nil {
		s.logger.Warnf("mail saved to SMTP, but failed to persist mail to db: %v", err)
	}

	return nil
}

func buildRawMail(from, to, subject, body string) string {
	headers := []string{
		fmt.Sprintf("From: %s", from),
		fmt.Sprintf("To: %s", to),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"",
		body,
	}
	return strings.Join(headers, "\r\n")
}
