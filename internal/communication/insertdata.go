package communication

import messages "github.com/Senso-Care/Unstacker/pkg/interface"

type InsertData interface {
	InsertMeasure(measure *messages.Measure, sensor *string)
}
