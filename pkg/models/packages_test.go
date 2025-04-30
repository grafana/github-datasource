package models

import (
	"testing"

	"github.com/shurcooL/githubv4"
	"github.com/stretchr/testify/require"
)

func Test_validatePackageType(t *testing.T) {
	tests := []struct {
		name        string
		packageType githubv4.PackageType
		wantErr     bool
		errMessage  string
	}{
		// Valid package types
		{
			name:        "valid PYPI package type",
			packageType: githubv4.PackageTypePypi,
			wantErr:     false,
		},
		{
			name:        "valid Docker package type",
			packageType: githubv4.PackageTypeDocker,
			wantErr:     false,
		},
		{
			name:        "valid Maven package type",
			packageType: githubv4.PackageTypeMaven,
			wantErr:     false,
		},
		{
			name:        "valid Debian package type",
			packageType: githubv4.PackageTypeDebian,
			wantErr:     false,
		},
		// Not supported package types by GraphQL API anymore
		{
			name:        "not supported NPM package type",
			packageType: githubv4.PackageTypeNpm,
			wantErr:     true,
			errMessage:  `package type "NPM" is not supported`,
		},
		{
			name:        "not supported Rubygems package type",
			packageType: githubv4.PackageTypeRubygems,
			wantErr:     true,
			errMessage:  `package type "RUBYGEMS" is not supported`,
		},
		{
			name:        "not supported Nuget package type",
			packageType: githubv4.PackageTypeNuget,
			wantErr:     true,
			errMessage:  `package type "NUGET" is not supported`,
		},
		// Invalid package types
		{
			name:        "invalid package type",
			packageType: "INVALID",
			wantErr:     true,
			errMessage:  `invalid package type "INVALID"`,
		},
		{
			name:        "empty package type",
			packageType: "",
			wantErr:     true,
			errMessage:  `invalid package type ""`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePackageType(tt.packageType)
			if tt.wantErr {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.errMessage)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
