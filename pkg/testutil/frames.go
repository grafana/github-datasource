package testutil

import (
	"fmt"
	"path/filepath"

	"github.com/grafana/grafana-github-datasource/pkg/dfutil"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/experimental"
)

const UpdateGoldenFiles bool = true

func CheckGoldenFramer(name string, f dfutil.Framer) error {
	dr := backend.DataResponse{
		Frames: f.Frame(),
	}
	return experimental.CheckGoldenDataResponse(filepath.Join("testdata", fmt.Sprintf("%s.golden.txt", name)), &dr, UpdateGoldenFiles)
}
