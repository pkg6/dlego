package dlego

import (
	"context"
	"github.com/go-acme/lego/v4/certificate"
)

type IProvider interface {
	SetCertificateName(certificate string)
	SetPrivateKeyName(privateKey string)
	SetLogger(logger ILogger)
	Deploy(ctx context.Context, certificate *certificate.Resource)
}
