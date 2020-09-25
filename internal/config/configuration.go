package config

type Configuration struct {
	MqServer MqServerConfiguration
	Database DatabaseConfiguration
	Cores  int
}