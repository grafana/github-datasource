package github

import (
	"context"
	"strings"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/shurcooL/githubv4"
)

// Member represents a team member
type Member struct {
	Login     string
	Name      string
	AvatarURL string
}

// QueryGetTeams is the GraphQL query for retrieving teams from an organization
type QueryGetTeams struct {
	Organization struct {
		Teams struct {
			Nodes []struct {
				ID          githubv4.ID
				Name        string
				Description string
				Privacy     githubv4.TeamPrivacy
				Members     struct {
					TotalCount int64
					Nodes      []struct {
						Login     string
						Name      string
						AvatarURL string `graphql:"avatarUrl"`
					}
				} `graphql:"members(first: 100)"`
				Repositories struct {
					TotalCount int64
				} `graphql:"repositories"`
				ParentTeam *struct {
					ID   githubv4.ID
					Name string
				}
				URL string
			}
			PageInfo models.PageInfo
		} `graphql:"teams(first: 100, after: $cursor)"`
	} `graphql:"organization(login: $login)"`
}

// QueryGetSingleTeam is the GraphQL query for retrieving a single team by slug
type QueryGetSingleTeam struct {
	Organization struct {
		Team struct {
			ID          githubv4.ID
			Name        string
			Description string
			Privacy     githubv4.TeamPrivacy
			Members     struct {
				TotalCount int64
				Nodes      []struct {
					Login     string
					Name      string
					AvatarURL string `graphql:"avatarUrl"`
				}
			} `graphql:"members(first: 100)"`
			Repositories struct {
				TotalCount int64
			} `graphql:"repositories"`
			ParentTeam *struct {
				ID   githubv4.ID
				Name string
			}
			URL string
		} `graphql:"team(slug: $teamSlug)"`
	} `graphql:"organization(login: $login)"`
}

// TeamEntry represents a single team entry
type TeamEntry struct {
	ID           string
	Name         string
	Description  string
	Privacy      string
	MembersCount int64
	Members      []Member
	ReposCount   int64
	ParentTeam   string
	URL          string
}

// Teams is a list of team entries
type Teams []TeamEntry

// Frames converts the list of teams to a Grafana DataFrame
func (t Teams) Frames() data.Frames {
	backend.Logger.Info("Creating teams data frame", "entries_count", len(t))

	names := make([]string, len(t))
	membersLogins := make([]string, len(t))

	for i, team := range t {
		names[i] = team.Name

		// Create comma-separated list of member logins
		memberLogins := make([]string, len(team.Members))
		for j, member := range team.Members {
			memberLogins[j] = member.Login
		}
		membersLogins[i] = strings.Join(memberLogins, ", ")
	}

	frame := data.NewFrame(
		"teams",
		data.NewField("name", nil, names),
		data.NewField("member_logins", nil, membersLogins),
	)

	return data.Frames{frame}
}

// GetTeams retrieves teams from an organization
func GetTeams(ctx context.Context, client models.Client, opts models.ListTeamsOptions) (Teams, error) {
	backend.Logger.Info("GetTeams called", "organization", opts.Organization, "query", opts.Query)

	// If query looks like a team slug (no spaces, alphanumeric with hyphens/underscores), try single team query first
	if opts.Query != "" && isTeamSlug(opts.Query) {
		team, err := getSingleTeam(ctx, client, opts.Organization, opts.Query)
		if err == nil && team != nil {
			backend.Logger.Info("Retrieved single team", "name", team.Name)
			return Teams{*team}, nil
		}
		// If single team query fails, fall back to list and filter
		backend.Logger.Info("Single team query failed, falling back to list and filter", "error", err)
	}

	// Default behavior: list all teams and filter
	return getAllTeams(ctx, client, opts)
}

// isTeamSlug checks if a string looks like a GitHub team slug
func isTeamSlug(query string) bool {
	// Team slugs typically don't contain spaces and are URL-friendly
	return !strings.Contains(query, " ") && len(query) > 0
}

// getSingleTeam retrieves a specific team by slug
func getSingleTeam(ctx context.Context, client models.Client, organization, teamSlug string) (*TeamEntry, error) {
	variables := map[string]interface{}{
		"login":    githubv4.String(organization),
		"teamSlug": githubv4.String(teamSlug),
	}

	q := &QueryGetSingleTeam{}
	if err := client.Query(ctx, q, variables); err != nil {
		return nil, err
	}

	team := q.Organization.Team
	parentTeamName := ""
	if team.ParentTeam != nil {
		parentTeamName = team.ParentTeam.Name
	}

	// Convert member nodes to Member structs
	members := make([]Member, len(team.Members.Nodes))
	for i, member := range team.Members.Nodes {
		members[i] = Member{
			Login:     member.Login,
			Name:      member.Name,
			AvatarURL: member.AvatarURL,
		}
	}

	teamEntry := &TeamEntry{
		ID:           team.ID.(string),
		Name:         team.Name,
		Description:  team.Description,
		Privacy:      string(team.Privacy),
		MembersCount: team.Members.TotalCount,
		Members:      members,
		ReposCount:   team.Repositories.TotalCount,
		ParentTeam:   parentTeamName,
		URL:          team.URL,
	}

	return teamEntry, nil
}

// getAllTeams retrieves all teams from an organization with optional filtering
func getAllTeams(ctx context.Context, client models.Client, opts models.ListTeamsOptions) (Teams, error) {
	var teams []TeamEntry
	var cursor *githubv4.String

	for {
		variables := map[string]interface{}{
			"login":  githubv4.String(opts.Organization),
			"cursor": cursor,
		}

		q := &QueryGetTeams{}
		if err := client.Query(ctx, q, variables); err != nil {
			backend.Logger.Error("Failed to query teams", "error", err)
			return nil, err
		}

		for _, team := range q.Organization.Teams.Nodes {
			// Filter by query if specified
			if opts.Query != "" {
				if !strings.Contains(strings.ToLower(team.Name), strings.ToLower(opts.Query)) &&
					!strings.Contains(strings.ToLower(team.Description), strings.ToLower(opts.Query)) {
					continue
				}
			}

			parentTeamName := ""
			if team.ParentTeam != nil {
				parentTeamName = team.ParentTeam.Name
			}

			// Convert member nodes to Member structs
			members := make([]Member, len(team.Members.Nodes))
			for i, member := range team.Members.Nodes {
				members[i] = Member{
					Login:     member.Login,
					Name:      member.Name,
					AvatarURL: member.AvatarURL,
				}
			}

			teamEntry := TeamEntry{
				ID:           team.ID.(string),
				Name:         team.Name,
				Description:  team.Description,
				Privacy:      string(team.Privacy),
				MembersCount: team.Members.TotalCount,
				Members:      members,
				ReposCount:   team.Repositories.TotalCount,
				ParentTeam:   parentTeamName,
				URL:          team.URL,
			}

			teams = append(teams, teamEntry)
		}

		if !q.Organization.Teams.PageInfo.HasNextPage {
			break
		}
		cursor = &q.Organization.Teams.PageInfo.EndCursor
	}

	backend.Logger.Info("Retrieved teams", "count", len(teams))
	return teams, nil
}
