package communication

import (
	"fmt"
	"github.com/Senso-Care/Unstacker/internal/config"
	messages "github.com/Senso-Care/Unstacker/pkg/interface"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

type InfluxDBConnector struct {
	InfluxClient influxdb2.Client
	WriteApi     api.WriteAPI
}

func NewConnector(configuration *config.DatabaseConfiguration) *InfluxDBConnector {
	client := influxdb2.NewClient(configuration.ConnectionUri, fmt.Sprintf("%s:%s", configuration.Username, configuration.Password))
	writeAPI := client.WriteAPI("", configuration.DbName+configuration.RetentionPolicy)
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
	log.Infof("Inserting %s", measure)
	connector.WriteApi.WritePoint(point)
}

func MeasureToPoint(measure *messages.Measure, sensor *string) *write.Point {
	measurement := strings.Split(*sensor, "-")[0]
	point := influxdb2.NewPointWithMeasurement(measurement).
		AddTag("sensor", *sensor).
		AddField("v", fmt.Sprintf("%f", *measure.Value)).
		SetTime(time.Unix(int64(*measure.Timestamp), 0))
	return point
}
