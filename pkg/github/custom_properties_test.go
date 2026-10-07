package github

import (
	"testing"

	googlegithub "github.com/google/go-github/v84/github"
)

func TestCustomPropertiesFrames(t *testing.T) {
	properties := CustomProperties{
		{
			PropertyName:  googlegithub.Ptr("ownership"),
			ValueType:     googlegithub.PropertyValueTypeSingleSelect,
			Required:      googlegithub.Ptr(true),
			DefaultValue:  "Unknown",
			AllowedValues: []string{"Unknown", "Platform-Engineering"},
		},
		{
			PropertyName: googlegithub.Ptr("migration"),
			ValueType:    googlegithub.PropertyValueTypeTrueFalse,
		},
	}

	frame := properties.Frames()[0]
	if frame.Rows() != 3 {
		t.Fatalf("expected 3 rows, got %d", frame.Rows())
	}
	if got := frame.Fields[0].At(0).(string); got != "ownership" {
		t.Fatalf("expected ownership property, got %q", got)
	}
	if got := frame.Fields[6].At(1).(*string); got == nil || *got != "Platform-Engineering" {
		t.Fatalf("expected allowed value Platform-Engineering, got %v", got)
	}
	if got := frame.Fields[6].At(2).(*string); got != nil {
		t.Fatalf("expected no allowed value for boolean property, got %v", got)
	}
}
