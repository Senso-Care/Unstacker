package communication

import (
	"crypto/tls"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/Senso-Care/Unstacker/internal/config"
	"github.com/Senso-Care/Unstacker/pkg/messages"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func parseAndInsert(bytes []byte, topic *string, inserter InsertData) {
	measure := messages.Measure{}
	log.Debug("Payload size: ", len(bytes))
	if err := proto.Unmarshal(bytes, &measure); err != nil {
		log.Error("Error while decoding protobuf message ", err)
	} else {
		sensor := path.Base(*topic)
		inserter.InsertMeasure(&measure, &sensor)
	}
}

func createConnectionOptions(configuration *config.MqServerConfiguration, broker *string, onMessageReceived *func(MQTT.Client, MQTT.Message)) *MQTT.ClientOptions {
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
		if token := c.Subscribe(configuration.Topic, byte(configuration.QOS), *onMessageReceived); token.Wait() && token.Error() != nil {
			log.Panicf("Error subscribing to topic %s", configuration.Topic)
			panic(token.Error())
		}
	}

	return connOpts
}

func Listen(configuration *config.MqServerConfiguration, inserter InsertData) *MQTT.Client {
	broker := "tcp://" + configuration.HostIp + ":" + strconv.Itoa(configuration.Port)
	onMessageReceived := func(client MQTT.Client, message MQTT.Message) {
		log.Debugf("Received message on topic: %s", message.Topic())
		topic := message.Topic()
		go parseAndInsert(message.Payload(), &topic, inserter)
	}
	connOpts := createConnectionOptions(configuration, &broker, &onMessageReceived)

	client := MQTT.NewClient(connOpts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Debug("Error connecting to server")
		panic(token.Error())
	} else {
		log.Infof("Connected to %s\n", broker)
	}

	return &client
}
