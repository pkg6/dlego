package ssh

type Config struct {
	// SSH 主机。
	// 零值时默认为 "localhost"。
	SshHost string `json:"sshHost,omitempty"`
	// SSH 端口。
	// 零值时默认为 22。
	SshPort int32 `json:"sshPort,omitempty"`
	// SSH 登录用户名。
	SshUsername string `json:"sshUsername,omitempty"`
	// SSH 登录密码。
	SshPassword string `json:"sshPassword,omitempty"`
	// SSH 登录私钥。
	SshKey string `json:"sshKey,omitempty"`
	// SSH 登录私钥口令。
	SshKeyPassphrase string `json:"sshKeyPassphrase,omitempty"`
	// 是否回退使用 SCP。
	UseSCP bool `json:"useSCP,omitempty"`
	// 前置命令。
	PreCommand string `json:"preCommand,omitempty"`
	// 后置命令。
	PostCommand string `json:"postCommand,omitempty"`

	// 输出证书文件路径。
	CertPath string `json:"certPath,omitempty"`
	// 输出私钥文件路径。
	KeyPath string `json:"KeyPath,omitempty"`
}
