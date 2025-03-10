package aliyunoss

import (
	"context"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
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
	client, err := createSdkClient(p.Config.AccessKeyId, p.Config.AccessKeySecret, p.Config.Region)
	if err != nil {
		return errors.Wrap(err, "failed to create sdk client")
	}
	if p.Config.Bucket == "" {
		return errors.New("config `bucket` is required")
	}
	if p.Config.Domain == "" {
		return errors.New("config `domain` is required")
	}
	return client.PutBucketCnameWithCertificate(p.Config.Bucket, oss.PutBucketCname{
		Cname: p.Config.Domain,
		CertificateConfiguration: &oss.CertificateConfiguration{
			Certificate: string(certificate.Certificate),
			PrivateKey:  string(certificate.PrivateKey),
			Force:       true,
		},
	})
}
