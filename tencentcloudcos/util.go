package tencentcloudcos

import (
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tcSsl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
)

func DeployCertificateInstance(sdkClient *tcSsl.Client, certId string, config *Config) (response *tcSsl.DeployCertificateInstanceResponse, err error) {
	deployCertificateInstanceReq := tcSsl.NewDeployCertificateInstanceRequest()
	deployCertificateInstanceReq.CertificateId = common.StringPtr(certId)
	deployCertificateInstanceReq.ResourceType = common.StringPtr("cos")
	deployCertificateInstanceReq.Status = common.Int64Ptr(1)
	deployCertificateInstanceReq.InstanceIdList = common.StringPtrs([]string{fmt.Sprintf("%s#%s#%s", config.Region, config.Bucket, config.Domain)})
	return sdkClient.DeployCertificateInstance(deployCertificateInstanceReq)
}
