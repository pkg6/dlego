package tencentcloudclb

type ResourceType string

const (
	// 资源类型：通过 SSL 服务部署到云资源实例。
	RESOURCE_TYPE_VIA_SSLDEPLOY = ResourceType("ssl-deploy")
	// 资源类型：部署到指定负载均衡器。
	RESOURCE_TYPE_LOADBALANCER = ResourceType("loadbalancer")
	// 资源类型：部署到指定监听器。
	RESOURCE_TYPE_LISTENER = ResourceType("listener")
	// 资源类型：部署到指定转发规则域名。
	RESOURCE_TYPE_RULEDOMAIN = ResourceType("ruledomain")
)

type Config struct {
	// 腾讯云 SecretId。
	SecretId string `json:"secretId"`
	// 腾讯云 SecretKey。
	SecretKey string `json:"secretKey"`
	// 腾讯云地域。
	Region string `json:"region"`
	// 部署资源类型。
	ResourceType ResourceType `json:"resourceType"`
	// 负载均衡器 ID。
	// 部署资源类型为 [RESOURCE_TYPE_SSLDEPLOY]、[RESOURCE_TYPE_LOADBALANCER]、[RESOURCE_TYPE_RULEDOMAIN] 时必填。
	LoadbalancerId string `json:"loadbalancerId,omitempty"`
	// 负载均衡监听 ID。
	// 部署资源类型为 [RESOURCE_TYPE_SSLDEPLOY]、[RESOURCE_TYPE_LOADBALANCER]、[RESOURCE_TYPE_LISTENER]、[RESOURCE_TYPE_RULEDOMAIN] 时必填。
	ListenerId string `json:"listenerId,omitempty"`
	// SNI 域名或七层转发规则域名（支持泛域名）。
	// 部署资源类型为 [RESOURCE_TYPE_SSLDEPLOY] 时选填；部署资源类型为 [RESOURCE_TYPE_RULEDOMAIN] 时必填。
	Domain string `json:"domain,omitempty"`
}
