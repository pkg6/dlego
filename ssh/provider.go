package ssh

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

func (p *Provider) WithLogger(logger dlego.ILogger) {
	p.logger = logger
}

func (p *Provider) Deploy(ctx context.Context, certificate *certificate.Resource) error {
	// 连接
	client, err := createSshClient(
		p.Config.SshHost,
		p.Config.SshPort,
		p.Config.SshUsername,
		p.Config.SshPassword,
		p.Config.SshKey,
		p.Config.SshKeyPassphrase,
	)
	if err != nil {
		return errors.Wrap(err, "failed to create ssh client")
	}
	defer client.Close()

	p.logger.LogF("SSH connected")

	// 执行前置命令
	if p.Config.PreCommand != "" {
		stdout, stderr, err := execSshCommand(client, p.Config.PreCommand)
		if err != nil {
			return errors.Wrapf(err, "failed to execute pre-command: stdout: %s, stderr: %s", stdout, stderr)
		}

		p.logger.LogF("SSH pre-command executed %s", stdout)
	}

	if err := writeFile(client, p.Config.UseSCP, p.Config.CertPath, certificate.Certificate); err != nil {
		return errors.Wrap(err, "failed to upload certificate file")
	}

	p.logger.LogF("certificate file uploaded")

	if err := writeFile(client, p.Config.UseSCP, p.Config.KeyPath, certificate.PrivateKey); err != nil {
		return errors.Wrap(err, "failed to upload private key file")
	}

	p.logger.LogF("private key file uploaded")

	// 执行后置命令
	if p.Config.PostCommand != "" {
		stdout, stderr, err := execSshCommand(client, p.Config.PostCommand)
		if err != nil {
			return errors.Wrapf(err, "failed to execute post-command, stdout: %s, stderr: %s", stdout, stderr)
		}

		p.logger.LogF("SSH post-command executed", stdout)
	}

	return nil
}
