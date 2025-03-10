package tencentcloudclb

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/pkg6/dlego/common/tencentcloud"
	tcClb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tcSsl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
)

func (p *Provider) deployViaSslService(clients *tencentcloud.Clients, config *Config, cloudCertId string) error {
	if config.LoadbalancerId == "" {
		return errors.New("config `loadbalancerId` is required")
	}
	if config.ListenerId == "" {
		return errors.New("config `listenerId` is required")
	}
	// 证书部署到 CLB 实例
	// REF: https://cloud.tencent.com/document/product/400/91667
	deployCertificateInstanceReq := tcSsl.NewDeployCertificateInstanceRequest()
	deployCertificateInstanceReq.CertificateId = common.StringPtr(cloudCertId)
	deployCertificateInstanceReq.ResourceType = common.StringPtr("clb")
	deployCertificateInstanceReq.Status = common.Int64Ptr(1)
	if config.Domain == "" {
		// 未指定 SNI，只需部署到监听器
		deployCertificateInstanceReq.InstanceIdList = common.StringPtrs([]string{fmt.Sprintf("%s|%s", config.LoadbalancerId, config.ListenerId)})
	} else {
		// 指定 SNI，需部署到域名
		deployCertificateInstanceReq.InstanceIdList = common.StringPtrs([]string{fmt.Sprintf("%s|%s|%s", config.LoadbalancerId, config.ListenerId, config.Domain)})
	}
	deployCertificateInstanceResp, err := clients.SSL.DeployCertificateInstance(deployCertificateInstanceReq)
	if err != nil {
		return errors.Wrap(err, "failed to execute sdk request 'ssl.DeployCertificateInstance'")
	}

	p.logger.LogF("已部署证书到云资源实例", deployCertificateInstanceResp.Response)

	return nil
}
func (p *Provider) deployToListener(clients *tencentcloud.Clients, config *Config, cloudCertId string) error {
	if config.LoadbalancerId == "" {
		return errors.New("config `loadbalancerId` is required")
	}
	if config.ListenerId == "" {
		return errors.New("config `listenerId` is required")
	}
	// 更新监听器证书
	if err := p.modifyListenerCertificate(clients, config.LoadbalancerId, config.ListenerId, cloudCertId); err != nil {
		return err
	}
	return nil
}
func (d *Provider) deployToRuleDomain(clients *tencentcloud.Clients, config *Config, cloudCertId string) error {
	if config.LoadbalancerId == "" {
		return errors.New("config `loadbalancerId` is required")
	}
	if config.ListenerId == "" {
		return errors.New("config `listenerId` is required")
	}
	if config.Domain == "" {
		return errors.New("config `domain` is required")
	}
	// 修改负载均衡七层监听器转发规则的域名级别属性
	// REF: https://cloud.tencent.com/document/api/214/38092
	modifyDomainAttributesReq := tcClb.NewModifyDomainAttributesRequest()
	modifyDomainAttributesReq.LoadBalancerId = common.StringPtr(config.LoadbalancerId)
	modifyDomainAttributesReq.ListenerId = common.StringPtr(config.ListenerId)
	modifyDomainAttributesReq.Domain = common.StringPtr(config.Domain)
	modifyDomainAttributesReq.Certificate = &tcClb.CertificateInput{
		SSLMode: common.StringPtr("UNIDIRECTIONAL"),
		CertId:  common.StringPtr(cloudCertId),
	}
	modifyDomainAttributesResp, err := clients.CLB.ModifyDomainAttributes(modifyDomainAttributesReq)
	if err != nil {
		return errors.Wrap(err, "failed to execute sdk request 'clb.ModifyDomainAttributes'")
	}
	d.logger.LogF("已修改七层监听器转发规则的域名级别属性", modifyDomainAttributesResp.Response)
	return nil
}
func (p *Provider) deployToLoadbalancer(clients *tencentcloud.Clients, config *Config, cloudCertId string) error {
	if config.LoadbalancerId == "" {
		return errors.New("config `loadbalancerId` is required")
	}

	// 查询监听器列表
	// REF: https://cloud.tencent.com/document/api/214/30686
	listenerIds := make([]string, 0)
	describeListenersReq := tcClb.NewDescribeListenersRequest()
	describeListenersReq.LoadBalancerId = common.StringPtr(config.LoadbalancerId)
	describeListenersResp, err := clients.CLB.DescribeListeners(describeListenersReq)
	if err != nil {
		return errors.Wrap(err, "failed to execute sdk request 'clb.DescribeListeners'")
	} else {
		if describeListenersResp.Response.Listeners != nil {
			for _, listener := range describeListenersResp.Response.Listeners {
				if listener.Protocol == nil || (*listener.Protocol != "HTTPS" && *listener.Protocol != "TCP_SSL" && *listener.Protocol != "QUIC") {
					continue
				}

				listenerIds = append(listenerIds, *listener.ListenerId)
			}
		}
	}

	p.logger.LogF("已查询到负载均衡器下的监听器", listenerIds)
	// 遍历更新监听器证书
	if len(listenerIds) == 0 {
		return errors.New("listener not found")
	} else {
		var errs []error
		for _, listenerId := range listenerIds {
			if err := p.modifyListenerCertificate(clients, config.LoadbalancerId, listenerId, cloudCertId); err != nil {
				errs = append(errs, err)
			}
		}
		if len(errs) > 0 {
			return fmt.Errorf("%v", errs)
		}
	}
	return nil
}

func (d *Provider) modifyListenerCertificate(clients *tencentcloud.Clients, cloudLoadbalancerId, cloudListenerId, cloudCertId string) error {
	// 查询监听器列表
	// REF: https://cloud.tencent.com/document/api/214/30686
	describeListenersReq := tcClb.NewDescribeListenersRequest()
	describeListenersReq.LoadBalancerId = common.StringPtr(cloudLoadbalancerId)
	describeListenersReq.ListenerIds = common.StringPtrs([]string{cloudListenerId})
	describeListenersResp, err := clients.CLB.DescribeListeners(describeListenersReq)
	if err != nil {
		return errors.Wrap(err, "failed to execute sdk request 'clb.DescribeListeners'")
	}
	if len(describeListenersResp.Response.Listeners) == 0 {
		return errors.New("listener not found")
	}
	d.logger.LogF("已查询到监听器属性", describeListenersResp.Response)
	// 修改监听器属性
	// REF: https://cloud.tencent.com/document/product/214/30681
	modifyListenerReq := tcClb.NewModifyListenerRequest()
	modifyListenerReq.LoadBalancerId = common.StringPtr(cloudLoadbalancerId)
	modifyListenerReq.ListenerId = common.StringPtr(cloudListenerId)
	modifyListenerReq.Certificate = &tcClb.CertificateInput{CertId: common.StringPtr(cloudCertId)}
	if describeListenersResp.Response.Listeners[0].Certificate != nil && describeListenersResp.Response.Listeners[0].Certificate.SSLMode != nil {
		modifyListenerReq.Certificate.SSLMode = describeListenersResp.Response.Listeners[0].Certificate.SSLMode
		modifyListenerReq.Certificate.CertCaId = describeListenersResp.Response.Listeners[0].Certificate.CertCaId
	} else {
		modifyListenerReq.Certificate.SSLMode = common.StringPtr("UNIDIRECTIONAL")
	}
	modifyListenerResp, err := clients.CLB.ModifyListener(modifyListenerReq)
	if err != nil {
		return errors.Wrap(err, "failed to execute sdk request 'clb.ModifyListener'")
	}
	d.logger.LogF("已修改监听器属性", modifyListenerResp.Response)
	return nil
}
