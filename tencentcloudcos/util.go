package tencentcloudcos

import (
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tcSsl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
)

func createSdkClient(secretId, secretKey, region string) (*tcSsl.Client, error) {
	credential := common.NewCredential(secretId, secretKey)
	client, err := tcSsl.NewClient(credential, region, profile.NewClientProfile())
	if err != nil {
		return nil, err
	}

	return client, nil
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

func DeployCertificateInstance(sdkClient *tcSsl.Client, certId string, config *Config) (response *tcSsl.DeployCertificateInstanceResponse, err error) {
	deployCertificateInstanceReq := tcSsl.NewDeployCertificateInstanceRequest()
	deployCertificateInstanceReq.CertificateId = common.StringPtr(certId)
	deployCertificateInstanceReq.ResourceType = common.StringPtr("cos")
	deployCertificateInstanceReq.Status = common.Int64Ptr(1)
	deployCertificateInstanceReq.InstanceIdList = common.StringPtrs([]string{fmt.Sprintf("%s#%s#%s", config.Region, config.Bucket, config.Domain)})
	return sdkClient.DeployCertificateInstance(deployCertificateInstanceReq)
}
