package communication

import "github.com/Senso-Care/Unstacker/pkg/messages"

type InsertData interface {
	InsertMeasure(measure *messages.Measure, sensor *string)
}
