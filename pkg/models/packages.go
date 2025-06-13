package models

import (
	"fmt"
	"slices"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/shurcooL/githubv4"
)

// ListPackagesOptions provides options when retrieving commits
type ListPackagesOptions struct {
	Repository  string               `json:"repository"`
	Owner       string               `json:"owner"`
	Names       string               `json:"names"`
	PackageType githubv4.PackageType `json:"packageType"`
}

// PackagesOptionsWithRepo adds Owner and Repo to a ListPackagesOptions. This is just for convenience
func PackagesOptionsWithRepo(opt ListPackagesOptions, owner string, repo string) (ListPackagesOptions, error) {
	err := validatePackageType(opt.PackageType)
	if err != nil {
		return ListPackagesOptions{}, err
	}

	return ListPackagesOptions{
		Owner:       owner,
		Repository:  repo,
		Names:       opt.Names,
		PackageType: opt.PackageType,
	}, nil
}

// validPackageTypes is a list of valid package types that are supported by the GitHub graphql API that we are using
var validPackageTypes = []githubv4.PackageType{
	githubv4.PackageTypeMaven,
	githubv4.PackageTypeDocker,
	githubv4.PackageTypeDebian,
	githubv4.PackageTypePypi,
}

// notSupportedPackageTypes is a list of package types that are not supported by the GitHub graphql API
// They were supported in the past but are not supported anymore and we want to return an error if they are used
var notSupportedPackageTypes = []githubv4.PackageType{
	githubv4.PackageTypeNpm,
	githubv4.PackageTypeRubygems,
	githubv4.PackageTypeNuget,
}

func validatePackageType(packageType githubv4.PackageType) error {
	if slices.Contains(validPackageTypes, packageType) {
		return nil
	}

	if slices.Contains(notSupportedPackageTypes, packageType) {
		return backend.DownstreamError(fmt.Errorf("package type %q is not supported. Valid types are: MAVEN, DOCKER, DEBIAN, PYPI", packageType))
	}
	return backend.DownstreamError(fmt.Errorf("invalid package type %q. Valid types are: MAVEN, DOCKER, DEBIAN, PYPI", packageType))
}
