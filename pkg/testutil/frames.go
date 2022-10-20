package testutil

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/grafana/github-datasource/pkg/dfutil"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/experimental"
)

// UpdateGoldenFiles defines whether or not to update the files locally after checking the responses for vailidity
const UpdateGoldenFiles bool = true

// CheckGoldenFramer checks the GoldenDataResponse and creates a standard file format for saving them
func CheckGoldenFramer(t *testing.T, name string, f dfutil.Framer) {
	dr := backend.DataResponse{
		Frames: f.Frames(),
	}
	experimental.CheckGoldenJSONResponse(t, filepath.Join("testdata"), fmt.Sprintf("%s.golden", name), &dr, UpdateGoldenFiles)
}
