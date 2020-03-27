package googlesheets

import (
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

func findTimeField(frame *data.Frame) int {
	for fieldIdx, f := range frame.Fields {
		ftype := f.Type()
		if ftype == data.FieldTypeTime || ftype == data.FieldTypeNullableTime {
			return fieldIdx
		}
	}
	return -1
}
