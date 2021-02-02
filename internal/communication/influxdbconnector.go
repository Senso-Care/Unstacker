package communication

import (
	"fmt"
	"strings"
	"time"

	"github.com/Senso-Care/Unstacker/internal/config"
	"github.com/Senso-Care/Unstacker/pkg/messages"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	log "github.com/sirupsen/logrus"
)

type InfluxDBConnector struct {
	InfluxClient influxdb2.Client
	WriteApi     api.WriteAPI
}

func NewConnector(configuration *config.DatabaseConfiguration) *InfluxDBConnector {
	client := influxdb2.NewClient(configuration.ConnectionUri, fmt.Sprintf("%s:%s", configuration.Username, configuration.Password))
	writeAPI := client.WriteAPI("", configuration.DbName+"/"+configuration.RetentionPolicy)
	errorsCh := writeAPI.Errors()
	connector := InfluxDBConnector{
		InfluxClient: client,
		WriteApi:     writeAPI,
	}
	// Create go proc for reading and logging errors
	go func(errorChannel <-chan error) {
		for err := range errorChannel {
			log.Warn("write error: %s", err.Error())
		}
	}(errorsCh)
	return &connector
}

func (connector *InfluxDBConnector) Close() {
	connector.InfluxClient.Close()
}

func (connector *InfluxDBConnector) InsertMeasure(measure *messages.Measure, sensor *string) {
	point := MeasureToPoint(measure, sensor)
	if point != nil {
		log.Infof("Inserting %s", measure)
		connector.WriteApi.WritePoint(point)
	}
}

func MeasureToPoint(measure *messages.Measure, sensor *string) *write.Point {
	measurement := strings.ToLower(strings.Split(*sensor, "-")[0])
	var timestamp time.Time
	if measure.Timestamp == 0 {
		timestamp = time.Now()
	} else {
		timestamp = time.Unix(measure.Timestamp, 0)
	}
	var value interface{}
	switch tmpValue := measure.Value.(type) {
	case *messages.Measure_FValue:
		value = tmpValue.FValue
	case *messages.Measure_IValue:
		value = tmpValue.IValue
	default:
		log.Error("Error, no value present in message")
		return nil
	}
	point := influxdb2.NewPointWithMeasurement(measurement).
		AddTag("sensor", strings.ToLower(*sensor)).
		AddField("v", value).
		SetTime(timestamp)
	return point
}
