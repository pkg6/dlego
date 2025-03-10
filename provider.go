package dlego

import (
	"context"
	"github.com/go-acme/lego/v4/certificate"
)

type IProvider interface {
	WithLogger(logger ILogger)
	Deploy(ctx context.Context, certificate *certificate.Resource) error
}
