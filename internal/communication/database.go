package communication

import (
	"fmt"
	"github.com/Senso-Care/Unstacker/internal/config"
	messages "github.com/Senso-Care/Unstacker/pkg/interface"
	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

func Connect(configuration *config.DatabaseConfiguration) (influxdb2.Client, api.WriteAPI) {
	client := influxdb2.NewClient(configuration.ConnectionUri, fmt.Sprintf("%s:%s", configuration.Username, configuration.Password))
	writeAPI := client.WriteAPI("", configuration.DbName+"/one_month")
	errorsCh := writeAPI.Errors()
	// Create go proc for reading and logging errors
	go errorLogger(errorsCh)
	return client, writeAPI
}

func errorLogger(errorChannel <-chan error) {
	for err := range errorChannel {
		log.Warn("write error: %s", err.Error())
	}
}

func MeasureToPoint(measure *messages.Measure, sensor *string) *write.Point {
	measurement := strings.Split(*sensor, "-")[0]

	point := influxdb2.NewPointWithMeasurement(measurement).
		AddTag("sensor", *sensor).
		AddField("v", measure.Value).
		SetTime(time.Unix(int64(*measure.Timestamp), 0))
	return point
}

func InsertMeasure(writeAPI api.WriteAPI, measure *messages.Measure, sensor *string) {
	point := MeasureToPoint(measure, sensor)
	log.Infof("Inserting %s", measure)
	writeAPI.WritePoint(point)
}
