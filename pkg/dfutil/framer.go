package dfutil

import (
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

// Framer is an interface that allows any type to be treated as a data frame
type Framer interface {
	Frames() data.Frames
}

// FrameResponse creates a backend.DataResponse that contains the Framer's data.Frames
func FrameResponse(f Framer) backend.DataResponse {
	return backend.DataResponse{
		Frames: f.Frames(),
	}
}

// FrameResponseWithError creates a backend.DataResponse with the error's contents (if not nil), and the Framer's data.Frames
// This function is particularly useful if you have a function that returns `(Framer, error)`, which is a very common pattern
func FrameResponseWithError(f Framer, err error) backend.DataResponse {
	if err != nil {
		return backend.DataResponse{
			Error: err,
		}
	}

	return FrameResponse(f)
}
