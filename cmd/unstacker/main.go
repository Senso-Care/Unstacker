package main

import (
	"crypto/tls"
	"fmt"
	"github.com/Senso-Care/Unstacker/internal/config"
	"github.com/Senso-Care/Unstacker/pkg/messages"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
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
	topics := [5]string{
		"temperature-bathroom",
		"temperature-kitchen",
		"temperature-wc",
		"temperature-livingroom",
		"temperature-bedroom",
	}
	i := 0
	start := uint64(time.Now().Unix())
	//timeToWait := uint64(10)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	for {
		timestamp := uint64(time.Now().Unix())
		for _, topic := range topics {
			topic = "/senso-care/sensors/" + topic
			value := float32(rand.Int31n(20)) + 10 + rand.Float32()
			measure := messages.Measure{
				Timestamp: &timestamp,
				Value:     &value,
			}
			if bytes, err := proto.Marshal(&measure); err != nil {
				log.Println("Error while unmarshalling: ", err)
			} else {
				go func(topic string, i int) {
					client.Publish(topic, byte(configuration.MqServer.QOS), false, bytes)
					log.Printf("Message sent to %s, nb %d", topic, i)
				}(topic, i)
			}
		}
		i += 1

		select {
		case msg := <-c:
			fmt.Println("Received shutdown signal", msg)
			shutdown(&client, i, timestamp-start)
		default:
		}
		time.Sleep(10 * time.Second)
	}

}

func shutdown(client *MQTT.Client, i int, duration uint64) {
	(*client).Disconnect(10)
	fmt.Printf("n: %d, t: %d, msg/s: %f\n", i, duration, float64(i)/float64(duration))
	os.Exit(0)
}
