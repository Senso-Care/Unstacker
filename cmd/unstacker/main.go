package main

import (
	"crypto/tls"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/Senso-Care/Unstacker/internal/config"
	"github.com/Senso-Care/Unstacker/pkg/messages"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
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

	topics := [1]string{
		//"sound-wc",
		//"sound-bathroom",
		//"sound-kitchen",
		"sound-livingroom",
		//"sound-bedroom",
	}
	i := 0
	start := time.Now().Unix()
	//timeToWait := uint64(10)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	nb := 7 * 24
	timestamp := time.Now().Unix() - (7 * 24 * 60 * 60)
	for count := 0; count < nb; count++ {
		timestamp = timestamp + (60 * 60)
		for _, topic := range topics {
			topic = "/senso-care/sensors/" + topic
			//value := float32(rand.Int31n(35-25)) + 25 + rand.Float32()
			value := rand.Int31n(35-25) + 25
			measure := messages.Measure{
				Timestamp: timestamp,
				Value:     &messages.Measure_IValue{IValue: value},
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

		/*select {
		case msg := <-c:
			fmt.Println("Received shutdown signal", msg)
			shutdown(&client, i, timestamp-start)
		default:
		}*/
		time.Sleep(10 * time.Millisecond)
	}
	shutdown(&client, i, timestamp-start)
}

func shutdown(client *MQTT.Client, i int, duration int64) {
	(*client).Disconnect(10)
	fmt.Printf("n: %d, t: %d, msg/s: %f\n", i, duration, float64(i)/float64(duration))
	os.Exit(0)
}
