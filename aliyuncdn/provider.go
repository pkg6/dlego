package aliyuncdn

import (
	"context"
	"fmt"
	aliyunCdn "github.com/alibabacloud-go/cdn-20180510/v5/client"
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
	domain := strings.TrimPrefix(p.Config.Domain, "*")
	client, err := createSdkClient(p.Config.AccessKeyId, p.Config.AccessKeySecret)
	if err != nil {
		return errors.Wrap(err, "failed to create sdk client")
	}
	setCdnDomainSSLCertificateReq := &aliyunCdn.SetCdnDomainSSLCertificateRequest{
		DomainName:  tea.String(domain),
		CertName:    tea.String(fmt.Sprintf("certimate-%d", time.Now().UnixMilli())),
		CertType:    tea.String("upload"),
		SSLProtocol: tea.String("on"),
		SSLPub:      tea.String(string(certificate.Certificate)),
		SSLPri:      tea.String(string(certificate.PrivateKey)),
	}
	setCdnDomainSSLCertificateResp, err := client.SetCdnDomainSSLCertificate(setCdnDomainSSLCertificateReq)
	if err != nil {
		return errors.Wrap(err, "failed to execute sdk request 'cdn.SetCdnDomainSSLCertificate'")
	}
	p.logger.LogF("已设置 CDN 域名证书 %s", setCdnDomainSSLCertificateResp)
	return nil
}
