package main

import (
	"crypto/tls"
	"fmt"
	"github.com/Senso-Care/Unstacker/internal/config"
	messages "github.com/Senso-Care/Unstacker/pkg/interface"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"time"
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

	//MQTT.DEBUG = log.New(os.Stdout, "", 0)
	//MQTT.ERROR = log.New(os.Stdout, "", 0)
	hostname, _ := os.Hostname()
	clientid := hostname + strconv.Itoa(time.Now().Second())

	broker := "tcp://" + configuration.MqServer.HostIp + ":" + strconv.Itoa(configuration.MqServer.Port)

	connOpts := MQTT.NewClientOptions().AddBroker(broker).SetClientID(clientid).SetCleanSession(true)
	if configuration.MqServer.Username != "" {
		connOpts.SetUsername(configuration.MqServer.Username)
		if configuration.MqServer.Password != "" {
			connOpts.SetPassword(configuration.MqServer.Password)
		}
	}
	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	connOpts.SetTLSConfig(tlsConfig)

	client := MQTT.NewClient(connOpts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Panic("Error connecting to server")
		panic(token.Error())
		return
	}
	fmt.Printf("Connected to %s\n", broker)
	topics := [100]string{}
	for i := 0; i < 100; i++ {
		topics[i] = "temperature-" + strconv.Itoa(rand.Int())
	}
	i := 0
	//start := uint64(time.Now().Unix())
	//timeToWait := uint64(10)
	for {
		timestamp := uint64(time.Now().Unix())
		/*if timestamp-start >= timeToWait {
			break
		}*/
		value := rand.Float32() * 100
		measure := messages.Measure{
			Timestamp: &timestamp,
			Value:     &value,
		}
		if bytes, err := proto.Marshal(&measure); err != nil {
			log.Println("Error while unmarshalling: ", err)
		} else {
			topic := "/senso-care/sensors/" + topics[rand.Int63n(int64(len(topics)-1))]
			client.Publish(topic, byte(configuration.MqServer.QOS), false, bytes)
			i += 1
			log.Printf("Message sent to %s, nb %d", topic, i)
		}
		time.Sleep(time.Second)
	}
	//client.Disconnect(10)
	//fmt.Printf("n: %d, t: %d, msg/s: %f\n", i, timeToWait, float64(i)/float64(timeToWait))
}
