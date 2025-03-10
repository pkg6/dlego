package local

type Config struct {
	// Shell 执行环境。
	// 前置命令。
	PreCommand string `json:"preCommand,omitempty"`
	// 后置命令。
	PostCommand string `json:"postCommand,omitempty"`

	// 输出证书文件路径。
	CertPath string `json:"certPath,omitempty"`
	// 输出私钥文件路径。
	KeyPath string `json:"keyPath,omitempty"`
}
