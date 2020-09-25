package config

type MqServerConfiguration struct {
	HostIp string
	Port int
	Topic string
	Username string
	Password string
	QOS int
}
