package config

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"os"
	"path"
	"strings"
)

const (
	MqHostIp          = "127.0.0.1"
	MqPort            = 1883
	MqTopic           = "/senso-care/sensors/+"
	MqQos             = 0
	MqUsername        = ""
	MqPassword        = ""
	DbConnectionUri   = "http://localhost:9999"
	DbUsername        = ""
	DbPassword        = ""
	DbName            = "sensocare"
	DbRetentionPolicy = "one_month"
	Cores             = 0
)

func defaultConfiguration() Configuration {
	return Configuration{
		MqServer: MqServerConfiguration{HostIp: MqHostIp, Port: MqPort, Topic: MqTopic, Username: MqUsername, Password: MqPassword, QOS: MqQos},
		Database: DatabaseConfiguration{ConnectionUri: DbConnectionUri, DbName: DbName, RetentionPolicy: DbRetentionPolicy, Username: DbUsername, Password: DbPassword},
		Cores:    Cores,
	}
}

func WriteExampleConfig() {
	yamlBytes, err := yaml.Marshal(defaultConfiguration())

	if err != nil {
		log.Panic("Error converting default configuration to yaml")
		panic(err)
	}

	viper.SetConfigType("yaml")
	if err := viper.ReadConfig(bytes.NewReader(yamlBytes)); err != nil {
		log.Panic("Error converting yaml default config to Viper")
		panic(err)
	}
	filepath := "./example-config.yaml"
	if err := viper.WriteConfigAs(filepath); err != nil {
		log.Panic("Error writing example configuration\n")
		panic(err)
	}
	log.Printf("Example configuration:\n%s\n", string(yamlBytes))
	log.Printf("Example configuration written to %s\n", filepath)
}

func initViper() *viper.Viper {
	viperConf := viper.New()
	viperConf.AddConfigPath(".")
	viperConf.SetDefault("MqServer.hostip", MqHostIp)
	viperConf.SetDefault("MqServer.port", MqPort)
	viperConf.SetDefault("MqServer.topic", MqTopic)
	viperConf.SetDefault("MqServer.username", MqUsername)
	viperConf.SetDefault("MqServer.password", MqPassword)
	viperConf.SetDefault("MqServer.QOS", MqQos)
	viperConf.SetDefault("Database.ConnectionUri", DbConnectionUri)
	viperConf.SetDefault("Database.Username", DbUsername)
	viperConf.SetDefault("Database.Password", DbPassword)
	viperConf.SetDefault("Database.DbName", DbName)
	viperConf.SetDefault("Database.RetentionPolicy", DbRetentionPolicy)
	viperConf.SetDefault("Cores", Cores)
	viperConf.SetEnvPrefix("senso_care")
	viperConf.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viperConf.AutomaticEnv()
	viperConf.SetConfigType("yaml")
	return viperConf
}

func createFlags() {
	generateDefaultConfiguration := pflag.Bool("default-config", false, "Generate an example configuration to ./example-config.yaml")
	pflag.StringP("config", "c", "", "Path to YAML config file")
	pflag.String("mq-hostip", MqHostIp, "Message queue Server host ip")
	pflag.Int("mq-port", MqPort, "Message queue Server port")
	pflag.String("mq-topic", MqTopic, "Message queue Server topic to listen on")
	pflag.String("mq-username", "", "Message queue Server username")
	pflag.String("mq-password", "", "Message queue Server password")
	pflag.Int("mq-qos", MqQos, "Message queue Server quality of service")
	pflag.String("db-connection-uri", DbConnectionUri, "Database connection uri")
	pflag.String("db-name", DbName, "Database name")
	pflag.String("db-retention-policy", DbRetentionPolicy, "Database retention policy")
	pflag.String("db-username", DbUsername, "Database username")
	pflag.String("db-password", DbPassword, "Database password")
	pflag.Int("cores", Cores, "Number of cores to use")
	pflag.Parse()
	if *generateDefaultConfiguration {
		WriteExampleConfig()
		// Job is only to write an example configuration, early exit
		os.Exit(0)
	}
}

func bindFlag(v *viper.Viper, nameInConfig string, nameInFlags string) {
	err := v.BindPFlag(nameInConfig, pflag.Lookup(nameInFlags))
	if err != nil {
		log.Panicf("Failed binding flag %s to config value %s", nameInFlags, nameInConfig)
		panic(err)
	}
}

func bindFlags(v *viper.Viper) {
	bindFlag(v, "MqServer.hostip", "mq-hostip")
	bindFlag(v, "MqServer.port", "mq-port")
	bindFlag(v, "MqServer.topic", "mq-topic")
	bindFlag(v, "MqServer.username", "mq-username")
	bindFlag(v, "MqServer.password", "mq-password")
	bindFlag(v, "MqServer.QOS", "mq-qos")
	bindFlag(v, "Database.ConnectionUri", "db-connection-uri")
	bindFlag(v, "Database.DbName", "db-name")
	bindFlag(v, "Database.RetentionPolicy", "db-retention-policy")
	bindFlag(v, "Database.Username", "db-username")
	bindFlag(v, "Database.Password", "db-password")
	bindFlag(v, "Cores", "cores")
}

func loadConfigFile(v *viper.Viper) {
	configurationPath := pflag.Lookup("config").Value.String()
	if len(configurationPath) >= 0 {
		file, err := os.Stat(configurationPath)
		if os.IsNotExist(err) {
			log.Debug("No configuration file given")
		} else {
			if file.IsDir() {
				v.SetConfigName("config")
				v.AddConfigPath(configurationPath)
			} else {
				v.AddConfigPath(path.Dir(configurationPath))
				filename := path.Base(configurationPath)
				v.SetConfigName(filename[:len(filename)-len(path.Ext(filename))])
			}
		}
	}
}

func LoadConfig() (Configuration, error) {
	createFlags()
	configuration := defaultConfiguration()
	v := initViper()
	loadConfigFile(v)
	bindFlags(v)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// it's okay if no config file is found
		} else {
			return configuration, err
		}
	}
	if err := v.Unmarshal(&configuration); err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
		return configuration, err
	}
	log.Debug("Configuration loaded")
	return configuration, nil
}
