package main

import (
	"github.com/Senso-Care/Unstacker/internal/communication"
	"github.com/Senso-Care/Unstacker/internal/config"
	log "github.com/sirupsen/logrus"
	"os"
	"runtime"
)

func main() {
	log.SetLevel(log.DebugLevel)

	configuration, err := config.LoadConfig()
	if err != nil {
		log.Fatal("error loading configuration: %s\n", err)
		os.Exit(1)
	}
	log.WithFields(log.Fields{
		"IP HOST": configuration.MqServer.HostIp,
		"PORT":    configuration.MqServer.Port,
	}).Debug("Server address loaded from configuration")
	log.WithField("GOMAXPROCS", configuration.Cores).Debug("Setting max number of cpus")
	runtime.GOMAXPROCS(configuration.Cores)
	communication.Listen(&configuration.MqServer)

}
