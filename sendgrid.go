package gpms

import (
	"fmt"
	"net/http"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGrid struct {
	baseConfig *config
}

func (s SendGrid) withConfig(c *config) Provider {
	s.baseConfig = c
	return s
}

func (s SendGrid) Send(subject string, opts ...SendOption) error {
	cfg := newSendConfig(opts...)
	m := mail.NewV3Mail()
	m.Subject = subject
	from := s.formatEmail(s.baseConfig.staticFrom)
	if f := cfg.from; f != nil {
		from = s.formatEmail(f)
	}
	m.SetFrom(from)
	m.SetTemplateID(cfg.templateID)
	m.AddPersonalizations(&mail.Personalization{
		To:                  append(s.formatEmails(s.baseConfig.staticTo), s.formatEmails(cfg.to)...),
		From:                from,
		CC:                  append(s.formatEmails(s.baseConfig.staticCC), s.formatEmails(cfg.cc)...),
		BCC:                 append(s.formatEmails(s.baseConfig.staticBCC), s.formatEmails(cfg.bcc)...),
		DynamicTemplateData: cfg.vars,
	})
	client := sendgrid.NewSendClient(s.baseConfig.APIKey)
	resp, err := client.Send(m)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("unexpected error from server with code %d: %s", resp.StatusCode, resp.Body)
	}
	return nil
}

func (s SendGrid) formatEmails(e []*Email) []*mail.Email {
	emails := make([]*mail.Email, len(e))
	for i, m := range e {
		emails[i] = s.formatEmail(m)
	}
	return emails
}

func (s SendGrid) formatEmail(e *Email) *mail.Email {
	return &mail.Email{
		Name:    e.Name,
		Address: e.Address,
	}
}
