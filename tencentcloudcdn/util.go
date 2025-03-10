package tencentcloudcdn

import (
	"github.com/pkg/errors"
	"github.com/pkg6/dlego/common/tencentcloud"
	tcCdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tcSsl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
)

func getDomainsByCertificateId(clients *tencentcloud.Clients, cloudCertId string) ([]string, error) {
	// 获取证书中的可用域名
	// REF: https://cloud.tencent.com/document/product/228/42491
	describeCertDomainsReq := tcCdn.NewDescribeCertDomainsRequest()
	describeCertDomainsReq.CertId = common.StringPtr(cloudCertId)
	describeCertDomainsReq.Product = common.StringPtr("cdn")
	describeCertDomainsResp, err := clients.CDN.DescribeCertDomains(describeCertDomainsReq)
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
func getDeployedDomainsByCertificateId(clients *tencentcloud.Clients, cloudCertId string) ([]string, error) {
	// 根据证书查询关联 CDN 域名
	// REF: https://cloud.tencent.com/document/product/400/62674
	describeDeployedResourcesReq := tcSsl.NewDescribeDeployedResourcesRequest()
	describeDeployedResourcesReq.CertificateIds = common.StringPtrs([]string{cloudCertId})
	describeDeployedResourcesReq.ResourceType = common.StringPtr("cdn")
	describeDeployedResourcesResp, err := clients.SSL.DescribeDeployedResources(describeDeployedResourcesReq)
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

func DeployCertificateInstance(clients *tencentcloud.Clients, certId string, instanceIds []string) (response *tcSsl.DeployCertificateInstanceResponse, err error) {
	// 证书部署到 CDN 实例
	// REF: https://cloud.tencent.com/document/product/400/91667
	deployCertificateInstanceReq := tcSsl.NewDeployCertificateInstanceRequest()
	deployCertificateInstanceReq.CertificateId = common.StringPtr(certId)
	deployCertificateInstanceReq.ResourceType = common.StringPtr("cdn")
	deployCertificateInstanceReq.Status = common.Int64Ptr(1)
	deployCertificateInstanceReq.InstanceIdList = common.StringPtrs(instanceIds)
	return clients.SSL.DeployCertificateInstance(deployCertificateInstanceReq)
}
