package tencentcloudcos

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
	if p.Config.Bucket == "" {
		return errors.New("config `bucket` is required")
	}
	if p.Config.Domain == "" {
		return errors.New("config `domain` is required")
	}
	client, err := createSdkClient(p.Config.SecretId, p.Config.SecretKey, p.Config.Region)
	if err != nil {
		return errors.Wrap(err, "failed to create sdk clients")
	}
	// 上传证书到 SSL
	upres, err := Upload(client, string(certificate.Certificate), string(certificate.PrivateKey))
	if err != nil {
		return errors.Wrap(err, "failed to upload certificate file")
	}
	p.logger.LogF("certificate file uploaded", upres)
	certId := *upres.Response.CertificateId
	response, err := DeployCertificateInstance(client, certId, p.Config)
	if err != nil {
		return errors.Wrap(err, "failed to deploy certificate instance")
	}
	p.logger.LogF("已部署证书到云资源实例", response.Response)
	return nil
}
