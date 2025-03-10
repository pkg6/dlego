package tencentcloud

import (
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tcSsl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
)

func SSLUploadCertificate(client *tcSsl.Client, certPem string, privkeyPem string) (response *tcSsl.UploadCertificateResponse, err error) {
	// 上传新证书
	// REF: https://cloud.tencent.com/document/product/400/41665
	uploadCertificateReq := tcSsl.NewUploadCertificateRequest()
	uploadCertificateReq.CertificatePublicKey = common.StringPtr(certPem)
	uploadCertificateReq.CertificatePrivateKey = common.StringPtr(privkeyPem)
	uploadCertificateReq.Repeatable = common.BoolPtr(false)
	return client.UploadCertificate(uploadCertificateReq)
}
