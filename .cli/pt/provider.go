package pt

import (
	"context"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/pkg6/dlego"
)

type Provider struct {
	Config *Config
	logger dlego.ILogger
}

func (p *Provider) WithLogger(logger dlego.ILogger) {
	p.logger = logger
}

func (p *Provider) Deploy(ctx context.Context, certificate *certificate.Resource) error {
	return nil
}
