package github

import "github.com/shurcooL/githubv4"

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
				PageInfo   PageInfo
			} `graphql:"fields(first: 100)"`
			Items struct {
				// Edges
				TotalCount int64
				Nodes      []ProjectItem
				PageInfo   PageInfo
			} `graphql:"items(first: 100, after: $cursor)"`
		} `graphql:"projectV2(number: $number)"`
	} `graphql:"organization(login: $login)"`
}

type QueryProjectByUser struct {
	User struct {
		ProjectV2 struct {
			Fields struct {
				TotalCount int64
				Nodes      []Field
				PageInfo   PageInfo
			} `graphql:"fields(first: 100)"`
			Items struct {
				// Edges
				TotalCount int64
				Nodes      []ProjectItem
				PageInfo   PageInfo
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
	// Creator
	// Owner
	// Readme - The project's readme.
	// resourcePath (URI!) The HTTP path for this project.
	// fields []Field
}

type ProjectV2ItemContent struct {
	DraftIssue  Content `graphql:"... on DraftIssue"`
	Issue       Content `graphql:"... on Issue"`
	PullRequest Content `graphql:"... on PullRequest"`
}

type Content struct {
	Title     *string
	Body      *string
	CreatedAt *githubv4.DateTime
	Assignees *Assignees `graphql:"assignees(first: 10)"`
}

type Assignees struct {
	PageInfo   PageInfo
	TotalCount int64
	Nodes      []User
}

type ProjectItemsWithFields struct {
	Items  []ProjectItem
	Fields []Field
}

type FieldValues struct {
	PageInfo   PageInfo
	TotalCount int64
	Nodes      []FieldValue
}

type Field struct {
	Common Common `graphql:"... on ProjectV2FieldCommon"`
}

type FieldValue struct {
	DateValue   ProjectV2ItemFieldDateValue         `graphql:"... on ProjectV2ItemFieldDateValue"`
	TextValue   ProjectV2ItemFieldTextValue         `graphql:"... on ProjectV2ItemFieldTextValue"`
	SelectValue ProjectV2ItemFieldSingleSelectValue `graphql:"... on ProjectV2ItemFieldSingleSelectValue"`
}

type ProjectV2ItemFieldSingleSelectValue struct {
	Name  *string
	Field CommonField
}

type CommonField struct {
	Common ProjectV2FieldCommon `graphql:"... on ProjectV2FieldCommon"`
}

type ProjectV2FieldConfiguration struct {
	Field          ProjectV2Field `graphql:"... on ProjectV2Field"`
	IterationField ProjectV2Field `graphql:"... on ProjectV2IterationField"`
	SelectField    ProjectV2Field `graphql:"... on ProjectV2SingleSelectField"`
}

type ProjectV2Field struct {
	name string
}

type ProjectV2ItemFieldValueCommon struct {
	Field ProjectV2ItemFieldTextValue `graphql:"... on ProjectV2ItemFieldTextValue"`
}

type ProjectV2ItemFieldTextValue struct {
	Text  *string
	Field CommonField
}

type ProjectV2ItemFieldDateValue struct {
	CreatedAt githubv4.DateTime
	//Date      githubv4.DateTime
	Date      *string
	UpdatedAt githubv4.DateTime
	Field     CommonField
}

type ProjectV2FieldCommon struct {
	Name string
}

type Common struct {
	Name     string
	DataType string
}
