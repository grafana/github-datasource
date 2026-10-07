package github

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	googlegithub "github.com/google/go-github/v84/github"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

type organizationPropertiesClient interface {
	GetOrganizationCustomProperties(context.Context, string) ([]*googlegithub.CustomProperty, *googlegithub.Response, error)
}

// CustomProperties is the organization custom-property schema.
type CustomProperties []*googlegithub.CustomProperty

// Frames converts custom-property definitions into a variable- and table-friendly frame.
func (properties CustomProperties) Frames() data.Frames {
	frame := data.NewFrame(
		"organization_custom_properties",
		data.NewField("property_name", nil, []string{}),
		data.NewField("value_type", nil, []string{}),
		data.NewField("required", nil, []bool{}),
		data.NewField("default_value", nil, []*string{}),
		data.NewField("description", nil, []*string{}),
		data.NewField("values_editable_by", nil, []*string{}),
		data.NewField("allowed_value", nil, []*string{}),
	)

	for _, property := range properties {
		allowedValues := property.AllowedValues
		if len(allowedValues) == 0 {
			appendCustomPropertyRow(frame, property, nil)
			continue
		}
		for _, allowedValue := range allowedValues {
			value := allowedValue
			appendCustomPropertyRow(frame, property, &value)
		}
	}

	return data.Frames{frame}
}

func appendCustomPropertyRow(frame *data.Frame, property *googlegithub.CustomProperty, allowedValue *string) {
	frame.AppendRow(
		property.GetPropertyName(),
		string(property.ValueType),
		property.GetRequired(),
		formatCustomPropertyValue(property.DefaultValue),
		property.Description,
		property.ValuesEditableBy,
		allowedValue,
	)
}

func formatCustomPropertyValue(value any) *string {
	if value == nil {
		return nil
	}
	if stringValue, ok := value.(string); ok {
		return &stringValue
	}
	encoded, err := json.Marshal(value)
	if err != nil {
		fallback := fmt.Sprint(value)
		return &fallback
	}
	formatted := string(encoded)
	return &formatted
}

func getCustomProperties(ctx context.Context, client organizationPropertiesClient, owner, propertyName string) (CustomProperties, error) {
	properties, _, err := client.GetOrganizationCustomProperties(ctx, owner)
	if err != nil {
		return nil, err
	}
	filtered := make(CustomProperties, 0, len(properties))
	for _, property := range properties {
		if propertyName == "" || property.GetPropertyName() == propertyName {
			filtered = append(filtered, property)
		}
	}
	sort.SliceStable(filtered, func(i, j int) bool {
		return filtered[i].GetPropertyName() < filtered[j].GetPropertyName()
	})
	return filtered, nil
}
