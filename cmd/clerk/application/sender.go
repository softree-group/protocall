package application

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"strings"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
)

type SenderConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string
}

type Sender struct {
	client *sasl.Client
	config *SenderConfig
}

func NewSender(c *SenderConfig) *Sender {
	client := sasl.NewPlainClient("", c.Username, c.Password)
	return &Sender{
		client: &client,
		config: c,
	}
}

const protocallTmpl = `
<p>{.}></p>`

func (s *Sender) renderTemplate(link string) ([]byte, error) {
	protoTmpl, err := template.New("base").Parse(protocallTmpl)
	if err != nil {
		return nil, err
	}
	msg := bytes.Buffer{}
	if err := protoTmpl.Execute(&msg, link); err != nil {
		return nil, err
	}
	return msg.Bytes(), nil
}

func (s *Sender) SendProtocol(ctx context.Context, linkToDownload string, users []string) error {
	// msg, err := s.renderTemplate(linkToDownload)
	// if err != nil {
	// 	return err
	// }

	if err := smtp.SendMail(
		fmt.Sprintf("%v:%v", s.config.Host, s.config.Port),
		*s.client,
		s.config.Username,
		users,
		strings.NewReader(linkToDownload),
	); err != nil {
		return err
	}

	return nil
}
