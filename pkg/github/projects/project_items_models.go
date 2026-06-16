package projects

import (
	"github.com/grafana/github-datasource/pkg/models"
	"github.com/shurcooL/githubv4"
)

// QueryProject lists project items in a project
// organization(login: "grafana") {
// 	projectV2(number: 218) {
// 		items(first: 50) {
// 				totalCount
// 				nodes {
// 						id
// 						createdAt
// 				}
// 		},
// 		fields(first: 50) {
// 			totalCount
// 			nodes{
// 		 	... on ProjectV2FieldCommon {
// 				 name
// 				 dataType
// 		 	}
// 			}
// 		},
// 	}
// }
type QueryProject struct {
	Organization struct {
		ProjectV2 struct {
			Fields struct {
				TotalCount int64
				Nodes      []Field
				PageInfo   models.PageInfo
			} `graphql:"fields(first: 100)"`
			Items struct {
				// Edges
				TotalCount int64
				Nodes      []ProjectItem
				PageInfo   models.PageInfo
			} `graphql:"items(first: 100, after: $cursor)"`
		} `graphql:"projectV2(number: $number)"`
	} `graphql:"organization(login: $login)"`
}

// QueryProjectByUser lists GitHub projects by User
type QueryProjectByUser struct {
	User struct {
		ProjectV2 struct {
			Fields struct {
				TotalCount int64
				Nodes      []Field
				PageInfo   models.PageInfo
			} `graphql:"fields(first: 100)"`
			Items struct {
				// Edges
				TotalCount int64
				Nodes      []ProjectItem
				PageInfo   models.PageInfo
			} `graphql:"items(first: 100, after: $cursor)"`
		} `graphql:"projectV2(number: $number)"`
	} `graphql:"user(login: $login)"`
}

// ProjectItem is a GitHub project item
type ProjectItem struct {
	Content     ProjectV2ItemContent
	FieldValues FieldValues `graphql:"fieldValues(first: 100)"`
	ID          string
	IsArchived  bool
	Type        string
	CreatedAt   githubv4.DateTime
	UpdatedAt   githubv4.DateTime
}

// ProjectV2ItemContent contains Content for a ProjectItem
type ProjectV2ItemContent struct {
	DraftIssue  Content      `graphql:"... on DraftIssue"`
	Issue       IssueContent `graphql:"... on Issue"`
	PullRequest IssueContent `graphql:"... on PullRequest"`
}

// Content of the ProjectItem
type Content struct {
	Title     *string
	Body      *string
	CreatedAt *githubv4.DateTime
	Assignees *Assignees `graphql:"assignees(first: 10)"`
}

// IssueContent of the ProjectItem
type IssueContent struct {
	Title     *string
	Body      *string
	CreatedAt *githubv4.DateTime
	Assignees *Assignees `graphql:"assignees(first: 10)"`
	Milestone *models.Milestone
	ClosedAt  *githubv4.DateTime
}

// Assignees to the ProjectItem
type Assignees struct {
	PageInfo   models.PageInfo
	TotalCount int64
	Nodes      []models.User
}

// ProjectItemsWithFields ...
type ProjectItemsWithFields struct {
	Items   []ProjectItem
	Fields  []Field
	Filters []models.Filter
}

// FieldValues are the values of each Field of a ProjectItem
type FieldValues struct {
	PageInfo   models.PageInfo
	TotalCount int64
	Nodes      []FieldValue
}

// Field is a field on a ProjectItem
type Field struct {
	Common ProjectV2FieldCommon `graphql:"... on ProjectV2FieldCommon"`
}

// FieldValue is a value for a Field
type FieldValue struct {
	DateValue      ProjectV2ItemFieldDateValue         `graphql:"... on ProjectV2ItemFieldDateValue"`
	TextValue      ProjectV2ItemFieldTextValue         `graphql:"... on ProjectV2ItemFieldTextValue"`
	SelectValue    ProjectV2ItemFieldSingleSelectValue `graphql:"... on ProjectV2ItemFieldSingleSelectValue"`
	IterationValue ProjectV2ItemFieldIterationValue    `graphql:"... on ProjectV2ItemFieldIterationValue"`
	LabelsValue    ProjectV2ItemFieldLabelValue        `graphql:"... on ProjectV2ItemFieldLabelValue"`
	NumberValue    ProjectV2ItemFieldNumberValue       `graphql:"... on ProjectV2ItemFieldNumberValue"`
	ReviewerValue  ProjectV2ItemFieldReviewerValue     `graphql:"... on ProjectV2ItemFieldReviewerValue"`
	RepoValue      ProjectV2ItemFieldRepositoryValue   `graphql:"... on ProjectV2ItemFieldRepositoryValue"`
}

// ProjectV2ItemFieldRepositoryValue ...
type ProjectV2ItemFieldRepositoryValue struct {
	Repository models.Repository
	Field      CommonField
}

// ProjectV2ItemFieldReviewerValue ...
type ProjectV2ItemFieldReviewerValue struct {
	Reviewers `graphql:"reviewers(first: 10)"`
	Field     CommonField
}

// Reviewers ...
type Reviewers struct {
	Nodes []Reviewer
}

// Reviewer ...
type Reviewer struct {
	models.User `graphql:"... on User"`
}

// ProjectV2ItemFieldNumberValue is a value for a Number field
type ProjectV2ItemFieldNumberValue struct {
	Number *float64
	Field  CommonField
}

// ProjectV2ItemFieldLabelValue is a value for a Labels field
type ProjectV2ItemFieldLabelValue struct {
	ProjectLabels `graphql:"labels(first: 10)"`
	Field         CommonField
}

// ProjectLabels ...
type ProjectLabels struct {
	Nodes []ProjectLabel
}

// ProjectLabel ...
type ProjectLabel struct {
	Name string
}

// ProjectV2ItemFieldIterationValue is a value for an Iteration field
type ProjectV2ItemFieldIterationValue struct {
	Title *string
	Field CommonField
}

// ProjectV2ItemFieldSingleSelectValue is a value for a SingleSelect field
type ProjectV2ItemFieldSingleSelectValue struct {
	Name  *string
	Field CommonField
}

// CommonField ...
type CommonField struct {
	Common ProjectV2FieldCommon `graphql:"... on ProjectV2FieldCommon"`
}

// ProjectV2ItemFieldTextValue is a value for a Text field
type ProjectV2ItemFieldTextValue struct {
	Text  *string
	Field CommonField
}

// ProjectV2ItemFieldDateValue is a value for a Date field
type ProjectV2ItemFieldDateValue struct {
	CreatedAt githubv4.DateTime
	Date      *string
	UpdatedAt githubv4.DateTime
	Field     CommonField
}

// ProjectV2FieldCommon is common to fields
type ProjectV2FieldCommon struct {
	Name     string
	DataType string
}
