package tencentcloudcdn

import (
	"github.com/pkg/errors"
	tcCdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tcSsl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
)

type Clients struct {
	ssl *tcSsl.Client
	cdn *tcCdn.Client
}

func createClients(secretId, secretKey string) (*Clients, error) {
	credential := common.NewCredential(secretId, secretKey)
	sslClient, err := tcSsl.NewClient(credential, "", profile.NewClientProfile())
	if err != nil {
		return nil, err
	}
	cdnClient, err := tcCdn.NewClient(credential, "", profile.NewClientProfile())
	if err != nil {
		return nil, err
	}
	return &Clients{
		ssl: sslClient,
		cdn: cdnClient,
	}, nil
}

func Upload(sdkClient *tcSsl.Client, certPem string, privkeyPem string) (response *tcSsl.UploadCertificateResponse, err error) {
	// 上传新证书
	// REF: https://cloud.tencent.com/document/product/400/41665
	uploadCertificateReq := tcSsl.NewUploadCertificateRequest()
	uploadCertificateReq.CertificatePublicKey = common.StringPtr(certPem)
	uploadCertificateReq.CertificatePrivateKey = common.StringPtr(privkeyPem)
	uploadCertificateReq.Repeatable = common.BoolPtr(false)
	return sdkClient.UploadCertificate(uploadCertificateReq)
}
func getDomainsByCertificateId(clients *Clients, cloudCertId string) ([]string, error) {
	// 获取证书中的可用域名
	// REF: https://cloud.tencent.com/document/product/228/42491
	describeCertDomainsReq := tcCdn.NewDescribeCertDomainsRequest()
	describeCertDomainsReq.CertId = common.StringPtr(cloudCertId)
	describeCertDomainsReq.Product = common.StringPtr("cdn")
	describeCertDomainsResp, err := clients.cdn.DescribeCertDomains(describeCertDomainsReq)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute sdk request 'cdn.DescribeCertDomains'")
	}

	domains := make([]string, 0)
	if describeCertDomainsResp.Response.Domains != nil {
		for _, domain := range describeCertDomainsResp.Response.Domains {
			domains = append(domains, *domain)
		}
	}
	return domains, nil
}
func getDeployedDomainsByCertificateId(clients *Clients, cloudCertId string) ([]string, error) {
	// 根据证书查询关联 CDN 域名
	// REF: https://cloud.tencent.com/document/product/400/62674
	describeDeployedResourcesReq := tcSsl.NewDescribeDeployedResourcesRequest()
	describeDeployedResourcesReq.CertificateIds = common.StringPtrs([]string{cloudCertId})
	describeDeployedResourcesReq.ResourceType = common.StringPtr("cdn")
	describeDeployedResourcesResp, err := clients.ssl.DescribeDeployedResources(describeDeployedResourcesReq)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute sdk request 'cdn.DescribeDeployedResources'")
	}

	domains := make([]string, 0)
	if describeDeployedResourcesResp.Response.DeployedResources != nil {
		for _, deployedResource := range describeDeployedResourcesResp.Response.DeployedResources {
			for _, resource := range deployedResource.Resources {
				domains = append(domains, *resource)
			}
		}
	}
	return domains, nil
}

func DeployCertificateInstance(clients *Clients, certId string, instanceIds []string) (response *tcSsl.DeployCertificateInstanceResponse, err error) {
	// 证书部署到 CDN 实例
	// REF: https://cloud.tencent.com/document/product/400/91667
	deployCertificateInstanceReq := tcSsl.NewDeployCertificateInstanceRequest()
	deployCertificateInstanceReq.CertificateId = common.StringPtr(certId)
	deployCertificateInstanceReq.ResourceType = common.StringPtr("cdn")
	deployCertificateInstanceReq.Status = common.Int64Ptr(1)
	deployCertificateInstanceReq.InstanceIdList = common.StringPtrs(instanceIds)
	return clients.ssl.DeployCertificateInstance(deployCertificateInstanceReq)
}
