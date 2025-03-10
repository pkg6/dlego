package tencentcloudcdn

import (
	"context"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/pkg/errors"
	"github.com/pkg6/dlego"
	"slices"
	"strings"
)

type Provider struct {
	Config *Config
	logger dlego.ILogger
}

func (p *Provider) WithLogger(logger dlego.ILogger) {
	p.logger = logger
}

func (p *Provider) Deploy(ctx context.Context, certificate *certificate.Resource) error {
	clients, err := createClients(p.Config.SecretId, p.Config.SecretKey)
	if err != nil {
		return err
	}
	upres, err := Upload(clients.ssl, string(certificate.Certificate), string(certificate.PrivateKey))
	if err != nil {
		return errors.Wrap(err, "failed to upload certificate file")
	}
	p.logger.LogF("certificate file uploaded %s", upres)

	certId := *upres.Response.CertificateId
	// 获取待部署的 CDN 实例
	// 如果是泛域名，根据证书匹配 CDN 实例
	instanceIds := make([]string, 0)
	if strings.HasPrefix(p.Config.Domain, "*.") {
		domains, err := getDomainsByCertificateId(clients, certId)
		if err != nil {
			return err
		}
		instanceIds = domains
	} else {
		instanceIds = append(instanceIds, p.Config.Domain)
	}
	// 跳过已部署的 CDN 实例
	if len(instanceIds) > 0 {
		deployedDomains, err := getDeployedDomainsByCertificateId(clients, certId)
		if err != nil {
			return err
		}

		temp := make([]string, 0)
		for _, instanceId := range instanceIds {
			if !slices.Contains(deployedDomains, instanceId) {
				temp = append(temp, instanceId)
			}
		}
		instanceIds = temp
	}
	if len(instanceIds) == 0 {
		p.logger.LogF("已部署过或没有要部署的 CDN 实例")
	} else {
		// 证书部署到 CDN 实例
		// REF: https://cloud.tencent.com/document/product/400/91667
		deployedDomains, err := DeployCertificateInstance(clients, certId, instanceIds)
		if err != nil {
			return errors.Wrap(err, "failed to execute sdk request 'ssl.DeployCertificateInstance'")
		}

		p.logger.LogF("已部署证书到云资源实例 %d %d", deployedDomains.Response)
	}
	return nil
}
