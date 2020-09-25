package communication

import (
	"crypto/tls"
	"encoding/json"
	"github.com/Senso-Care/daemons/internal/config"
	messages "github.com/Senso-Care/daemons/pkg/interface"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func onMessageReceived(client MQTT.Client, message MQTT.Message) {
	log.Debugf("Received message on topic: %s", message.Topic())
	go parse(message.Payload())
}

func parse(bytes []byte) {
	measure := messages.Measure{}
	if err := proto.Unmarshal(bytes, &measure); err != nil {
		log.Error("Error while decoding protobuf message ", err)
	} else {
		writeToDisk(&measure)
	}
}

func writeToDisk(measure *messages.Measure) {
	if bytes, err := json.Marshal(measure); err != nil {
		log.Error("Error while serializing to JSON ", err)
	} else {
		if err := ioutil.WriteFile("/tmp/go/" + strconv.Itoa(time.Now().Nanosecond()), bytes, 0644); err != nil {
			log.Error("Error while writing to disk", err)
		}
	}
}

func createConnectionOptions(configuration *config.MqServerConfiguration, broker *string) *MQTT.ClientOptions {
	hostname, _ := os.Hostname()
	clientid := hostname + strconv.Itoa(time.Now().Second())
	connOpts := MQTT.NewClientOptions().AddBroker(*broker).SetClientID(clientid).SetCleanSession(true)
	if configuration.Username != "" {
		connOpts.SetUsername(configuration.Username)
		if configuration.Password != "" {
			connOpts.SetPassword(configuration.Password)
		}
	}
	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	connOpts.SetTLSConfig(tlsConfig)

	connOpts.OnConnect = func(c MQTT.Client) {
		if token := c.Subscribe(configuration.Topic, byte(configuration.QOS), onMessageReceived); token.Wait() && token.Error() != nil {
			log.Panicf("Error subscribing to topic %s", configuration.Topic)
			panic(token.Error())
		}
	}

	return connOpts
}


func Listen(configuration *config.MqServerConfiguration) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	broker := "tcp://" + configuration.HostIp + ":" + strconv.Itoa(configuration.Port)

	connOpts := createConnectionOptions(configuration, &broker)

	client := MQTT.NewClient(connOpts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Panic("Error connecting to server")
		panic(token.Error())
	} else {
		log.Infof("Connected to %s\n", broker)
	}

	<-c
	log.Info("Graceful shutdown")
}