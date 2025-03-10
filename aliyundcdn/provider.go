package aliyundcdn

import (
	"context"
	"fmt"
	aliyunDcdn "github.com/alibabacloud-go/dcdn-20180115/v3/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/pkg/errors"
	"github.com/pkg6/dlego"
	"strings"
	"time"
)

type Provider struct {
	Config *Config
	logger dlego.ILogger
}

func (p *Provider) WithLogger(logger dlego.ILogger) {
	p.logger = logger
}

func (p *Provider) Deploy(ctx context.Context, certificate *certificate.Resource) error {
	client, err := createSdkClient(p.Config.AccessKeyId, p.Config.AccessKeySecret)
	if err != nil {
		return err
	}
	// "*.example.com" → ".example.com"，适配阿里云 DCDN 要求的泛域名格式
	domain := strings.TrimPrefix(p.Config.Domain, "*")
	// 配置域名证书
	// REF: https://help.aliyun.com/zh/edge-security-acceleration/dcdn/developer-reference/api-dcdn-2018-01-15-setdcdndomainsslcertificate
	setDcdnDomainSSLCertificateReq := &aliyunDcdn.SetDcdnDomainSSLCertificateRequest{
		DomainName:  tea.String(domain),
		CertName:    tea.String(fmt.Sprintf("certimate-%d", time.Now().UnixMilli())),
		CertType:    tea.String("upload"),
		SSLProtocol: tea.String("on"),
		SSLPub:      tea.String(string(certificate.Certificate)),
		SSLPri:      tea.String(string(certificate.PrivateKey)),
	}
	setDcdnDomainSSLCertificateResp, err := client.SetDcdnDomainSSLCertificate(setDcdnDomainSSLCertificateReq)
	if err != nil {
		return errors.Wrap(err, "failed to execute sdk request 'dcdn.SetDcdnDomainSSLCertificate'")
	}

	p.logger.LogF("已配置 DCDN 域名证书 %s", setDcdnDomainSSLCertificateResp)
	return nil
}
