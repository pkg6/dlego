package local

import (
	"context"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/pkg/errors"
	"github.com/pkg6/dlego"
)

type Provider struct {
	Config *Config

	logger dlego.ILogger
}

func (p *Provider) SetLogger(logger dlego.ILogger) {
	p.logger = logger
}

func (p *Provider) Deploy(ctx context.Context, certificate *certificate.Resource) error {
	// 执行前置命令
	if p.Config.PreCommand != "" {
		stdout, stderr, err := ExecCommand(p.Config.PreCommand, "")
		if err != nil {
			return errors.Wrapf(err, "failed to execute pre-command, stdout: %s, stderr: %s", stdout, stderr)
		}
		p.logger.LogF("pre-command executed %s", stdout)
	}
	if err := CopyFile(p.Config.CertPath, certificate.Certificate); err != nil {
		return errors.Wrap(err, "failed to save certificate file")
	}
	p.logger.LogF("certificate file saved")
	if err := CopyFile(p.Config.KeyPath, certificate.PrivateKey); err != nil {
		return errors.Wrap(err, "failed to save private key file")
	}
	p.logger.LogF("private key file saved")
	// 执行后置命令
	if p.Config.PostCommand != "" {
		stdout, stderr, err := ExecCommand(p.Config.PostCommand, "")
		if err != nil {
			return errors.Wrapf(err, "failed to execute post-command, stdout: %s, stderr: %s", stdout, stderr)
		}
		p.logger.LogF("post-command executed %s", stdout)
	}
	return nil
}
