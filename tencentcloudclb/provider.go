package tencentcloudclb

import (
	"context"
	"fmt"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/pkg/errors"
	"github.com/pkg6/dlego"
	"github.com/pkg6/dlego/common/tencentcloud"
)

type Provider struct {
	Config *Config
	logger dlego.ILogger
}

func (p *Provider) WithLogger(logger dlego.ILogger) {
	p.logger = logger
}

func (p *Provider) Deploy(ctx context.Context, certificate *certificate.Resource) error {
	clients, err := tencentcloud.NewClients(p.Config.SecretId, p.Config.SecretKey, p.Config.Region)
	if err != nil {
		return errors.Wrap(err, "failed to create sdk clients")
	}
	upres, err := tencentcloud.SSLUploadCertificate(clients.SSL, string(certificate.Certificate), string(certificate.PrivateKey))
	if err != nil {
		return errors.Wrap(err, "failed to upload certificate file")
	}
	p.logger.LogF("certificate file uploaded", upres)
	// 根据部署资源类型决定部署方式
	certId := *upres.Response.CertificateId
	switch p.Config.ResourceType {
	case RESOURCE_TYPE_VIA_SSLDEPLOY:
		if err := p.deployViaSslService(clients, p.Config, certId); err != nil {
			return err
		}
	case RESOURCE_TYPE_LOADBALANCER:
		if err := p.deployToLoadbalancer(clients, p.Config, certId); err != nil {
			return err
		}

	case RESOURCE_TYPE_LISTENER:
		if err := p.deployToListener(clients, p.Config, certId); err != nil {
			return err
		}
	case RESOURCE_TYPE_RULEDOMAIN:
		if err := p.deployToRuleDomain(clients, p.Config, certId); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported resource type: %s", p.Config.ResourceType)
	}
	return nil
}
