package sink

import (
	"encoding/json"

	"github.com/ximura/teleprobe/internal/metric"
)

type JSONFormatter struct{}

func (f *JSONFormatter) Marshal(m *metric.Measurement) ([]byte, error) {
	return json.Marshal(m)
}
