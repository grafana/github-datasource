package dfutil

import "github.com/grafana/grafana-plugin-sdk-go/data"

type Framer interface {
	Frame() data.Frames
}

type Fielder interface {
	Field() data.Field
}
