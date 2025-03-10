package tencentcloud

import (
	tcCdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tcSsl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
)

type Clients struct {
	SSL *tcSsl.Client
	CDN *tcCdn.Client
}

func NewClientsNoRegion(secretId, secretKey string) (*Clients, error) {
	credential := common.NewCredential(secretId, secretKey)
	sslClient, err := tcSsl.NewClient(credential, "", profile.NewClientProfile())
	if err != nil {
		return nil, err
	}
	cdnClient, err := tcCdn.NewClient(credential, "", profile.NewClientProfile())
	if err != nil {
		return nil, err
	}
	return &Clients{SSL: sslClient, CDN: cdnClient}, nil
}

func NewClients(secretId, secretKey, region string) (*Clients, error) {
	credential := common.NewCredential(secretId, secretKey)
	sslClient, err := tcSsl.NewClient(credential, region, profile.NewClientProfile())
	if err != nil {
		return nil, err
	}
	cdnClient, err := tcCdn.NewClient(credential, region, profile.NewClientProfile())
	if err != nil {
		return nil, err
	}
	return &Clients{SSL: sslClient, CDN: cdnClient}, nil
}
