package github

import (
	"context"
	"log"
	"strings"

	"github.com/grafana/grafana-github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/pkg/errors"
	"github.com/shurcooL/githubv4"
)

// {
//   repository(name: "grafana", owner: "grafana") {
//     packages(names: "", packageType: "", first: 10) {
//       nodes {
//         id
//         name
//         packageType
//         statistics {
//           downloadsTotalCount
//         }
//         versions(first: 10) {
//           nodes {
//             preRelease
//             platform
//             version
//             statistics {
//               downloadsTotalCount
//             }
//           }
//           pageInfo {
//             hasNextPage
//             endCursor
//           }
//         }
//       }
//       totalCount
//       pageInfo {
//         endCursor
//         hasNextPage
//       }
//     }
//   }
// }
type QueryListPackages struct {
	Repository struct {
		Packages struct {
			Nodes []struct {
				Name        string
				PackageType githubv4.PackageType
				Statistics  PackageStatistics
				Versions    struct {
					Nodes    []PackageVersion
					PageInfo PageInfo
				} `graphql:"versions(first: 100, after: $versionsCursor)"`
			}
			PageInfo PageInfo
		} `graphql:"packages(names: $names, packageType: $packageType, first: 100, after: $cursor)"`
	} `graphql:"repository(name: $name, owner: $owner)"`
}

type PackageStatistics struct {
	DownloadsTotalCount int64
}

type PackageVersion struct {
	PreRelease bool
	Platform   string
	Version    string
	Statistics PackageStatistics
}

type Package struct {
	Name        string
	PackageType githubv4.PackageType
	Statistics  PackageStatistics
	Versions    []PackageVersion
}

type Packages []Package

func (p Packages) Frames() data.Frames {
	frame := data.NewFrame(
		"packages",
		data.NewField("name", nil, []string{}),
		data.NewField("platform", nil, []string{}),
		data.NewField("version", nil, []string{}),
		data.NewField("type", nil, []string{}),
		data.NewField("prerelease", nil, []bool{}),
		data.NewField("downlods", nil, []int64{}),
	)

	for _, pkg := range p {
		for _, version := range pkg.Versions {
			frame.AppendRow(
				pkg.Name,
				version.Platform,
				version.Version,
				string(pkg.PackageType),
				version.PreRelease,
				version.Statistics.DownloadsTotalCount,
			)
		}
	}

	return data.Frames{frame}
}

// GetAllPackages lists packages in a repository
func GetAllPackages(ctx context.Context, client Client, opts models.ListPackagesOptions) (Packages, error) {
	s := strings.Split(opts.Names, ",")
	names := make([]githubv4.String, len(s))
	for i, v := range s {
		names[i] = githubv4.String(strings.TrimSpace(v))
	}

	var (
		variables = map[string]interface{}{
			"cursor":         (*githubv4.String)(nil),
			"versionsCursor": (*githubv4.String)(nil),
			"owner":          githubv4.String(opts.Owner),
			"name":           githubv4.String(opts.Repository),
			"names":          names,
			"packageType":    opts.PackageType,
		}

		packages = Packages{}
	)

	log.Println(opts.PackageType)
	log.Println(opts.PackageType)
	log.Println(opts.PackageType)
	log.Println(opts.PackageType)
	log.Println(opts.PackageType)

	for {
		q := &QueryListPackages{}
		if err := client.Query(ctx, q, variables); err != nil {
			return nil, errors.WithStack(err)
		}

		p := make(Packages, len(q.Repository.Packages.Nodes))

		// Retrieve versions for each package
		for i, v := range q.Repository.Packages.Nodes {
			p[i] = Package{
				Name:        v.Name,
				PackageType: v.PackageType,
				Statistics:  v.Statistics,
			}

			for {
				p[i].Versions = append(p[i].Versions, v.Versions.Nodes...)
				if !v.Versions.PageInfo.HasNextPage {
					variables["versionsCursor"] = (*githubv4.String)(nil)
					break
				}
				variables["versionsCursor"] = v.Versions.PageInfo.EndCursor
				if err := client.Query(ctx, q, variables); err != nil {
					return nil, errors.WithStack(err)
				}
			}
		}

		packages = append(packages, p...)

		if !q.Repository.Packages.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = q.Repository.Packages.PageInfo.EndCursor
	}

	return packages, nil
}
