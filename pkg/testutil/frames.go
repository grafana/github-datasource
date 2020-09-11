package testutil

import (
	"fmt"
	"path/filepath"

	"github.com/grafana/github-datasource/pkg/dfutil"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/experimental"
)

// UpdateGoldenFiles defines whether or not to update the files locally after checking the responses for vailidity
const UpdateGoldenFiles bool = true

// CheckGoldenFramer checks the GoldenDataResponse and creates a standard file format for saving them
func CheckGoldenFramer(name string, f dfutil.Framer) error {
	dr := backend.DataResponse{
		Frames: f.Frames(),
	}
	return experimental.CheckGoldenDataResponse(filepath.Join("testdata", fmt.Sprintf("%s.golden.txt", name)), &dr, UpdateGoldenFiles)
}
