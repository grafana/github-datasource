package github

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	googlegithub "github.com/google/go-github/v84/github"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

type enrichedRepositoriesMockClient struct {
}

func (client *enrichedRepositoriesMockClient) GetOrganizationCustomProperties(context.Context, string) ([]*googlegithub.CustomProperty, *googlegithub.Response, error) {
	return []*googlegithub.CustomProperty{
		{PropertyName: googlegithub.Ptr("ownership"), ValueType: googlegithub.PropertyValueTypeSingleSelect},
		{PropertyName: googlegithub.Ptr("migration"), ValueType: googlegithub.PropertyValueTypeTrueFalse},
	}, &googlegithub.Response{}, nil
}

func (client *enrichedRepositoriesMockClient) ListRepositoriesByOrg(_ context.Context, _ string, options *googlegithub.RepositoryListByOrgOptions) ([]*googlegithub.Repository, *googlegithub.Response, error) {
	if options.Page == 1 {
		return []*googlegithub.Repository{repositoryFixture(1, "managed"), repositoryFixture(2, "filtered")}, &googlegithub.Response{}, nil
	}
	return nil, &googlegithub.Response{}, nil
}

func repositoryFixture(id int64, name string) *googlegithub.Repository {
	now := googlegithub.Timestamp{Time: time.Date(2026, time.October, 6, 0, 0, 0, 0, time.UTC)}
	customProperties := map[string]any{"ownership": "GFOG", "migration": nil}
	if name == "managed" {
		customProperties = map[string]any{
			"ownership": "Platform-Engineering",
			"migration": false,
			"audiences": []string{"internal", "platform"},
		}
	}
	return &googlegithub.Repository{
		ID: googlegithub.Ptr(id), Name: googlegithub.Ptr(name), FullName: googlegithub.Ptr("rwest/" + name),
		HTMLURL: googlegithub.Ptr("https://rwe.ghe.com/rwest/" + name), Visibility: googlegithub.Ptr("internal"),
		Private: googlegithub.Ptr(true), DefaultBranch: googlegithub.Ptr("main"), Language: googlegithub.Ptr("Go"),
		CreatedAt: &now, UpdatedAt: &now, PushedAt: &now,
		CustomProperties: customProperties,
	}
}

func TestGetEnrichedRepositoriesFiltersAndJoins(t *testing.T) {
	client := &enrichedRepositoriesMockClient{}
	inventory, err := getEnrichedRepositories(context.Background(), client, "rwest", "", "ownership", "Platform-Engineering")
	if err != nil {
		t.Fatal(err)
	}
	if len(inventory.Repositories) != 1 || inventory.Repositories[0].GetName() != "managed" {
		t.Fatalf("expected only managed repository after property filtering, got %v", inventory.Repositories)
	}

	frame := inventory.Frames()[0]
	if frame.Rows() != 1 {
		t.Fatalf("expected one inventory row, got %d", frame.Rows())
	}
	assertFrameStringPointer(t, frame, "ownership", 0, "Platform-Engineering")
	assertFrameStringPointer(t, frame, "migration", 0, "false")
}

func TestEnrichedRepositoriesWildcardIncludesAllRepositories(t *testing.T) {
	inventory, err := getEnrichedRepositories(context.Background(), &enrichedRepositoriesMockClient{}, "rwest", "", "ownership", "*")
	if err != nil {
		t.Fatal(err)
	}
	if len(inventory.Repositories) != 2 {
		t.Fatalf("expected wildcard to include both repositories, got %d", len(inventory.Repositories))
	}
}

func TestCustomPropertyMatchesMultiSelect(t *testing.T) {
	if !customPropertyMatches([]string{"internal", "platform"}, "platform") {
		t.Fatal("expected multi-select property to match one selected value")
	}
	if customPropertyMatches([]string{"internal", "platform"}, "external") {
		t.Fatal("unexpected multi-select property match")
	}
}

func TestRepositoryCustomPropertiesDecodeBoolean(t *testing.T) {
	var repository googlegithub.Repository
	if err := json.Unmarshal([]byte(`{"full_name":"rwest/demo","custom_properties":{"migration":false}}`), &repository); err != nil {
		t.Fatal(err)
	}
	value, ok := repository.CustomProperties["migration"].(bool)
	if !ok || value {
		t.Fatalf("expected false boolean custom property, got %#v", repository.CustomProperties["migration"])
	}
}

func assertFrameString(t *testing.T, frame *data.Frame, name string, row int, expected string) {
	t.Helper()
	field := requireFrameField(t, frame, name)
	if got := field.At(row).(string); got != expected {
		t.Fatalf("expected %s=%q, got %q", name, expected, got)
	}
}

func assertFrameStringPointer(t *testing.T, frame *data.Frame, name string, row int, expected string) {
	t.Helper()
	field := requireFrameField(t, frame, name)
	got := field.At(row).(*string)
	if got == nil || *got != expected {
		t.Fatalf("expected %s=%q, got %v", name, expected, got)
	}
}

func requireFrameField(t *testing.T, frame *data.Frame, name string) *data.Field {
	t.Helper()
	for _, field := range frame.Fields {
		if field.Name == name {
			return field
		}
	}
	t.Fatalf("field %q not found", name)
	return nil
}
