package main

import (
	"github.com/Senso-Care/Unstacker/internal/communication"
	"github.com/Senso-Care/Unstacker/internal/config"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"runtime"
	"syscall"
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
	connector := communication.NewConnector(&configuration.Database)
	client := communication.Listen(&configuration.MqServer, connector)
	waitForShutdown(client, connector)
}

func waitForShutdown(client *MQTT.Client, connector *communication.InfluxDBConnector) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	log.Info("Graceful shutdown")
	(*client).Disconnect(10)
	connector.Close()
}
