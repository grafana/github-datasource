package models

import "github.com/shurcooL/githubv4"

// ListPackagesOptions provides options when retrieving commits
type ListPackagesOptions struct {
	Repository  string               `json:"repository"`
	Owner       string               `json:"owner"`
	Names       string               `json:"names"`
	PackageType githubv4.PackageType `json:"packageType"`
}

// PackagesOptionsWithRepo adds Owner and Repo to a ListPackagesOptions. This is just for convenience
func PackagesOptionsWithRepo(opt ListPackagesOptions, owner string, repo string) ListPackagesOptions {
	return ListPackagesOptions{
		Owner:       owner,
		Repository:  repo,
		Names:       opt.Names,
		PackageType: opt.PackageType,
	}
}
