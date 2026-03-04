package system_test

import (
	"testing"

	"pkgstats-cli/internal/system"
)

func TestParseOSId(t *testing.T) {
	testCases := []struct {
		name         string
		content      string
		expectedOSId string
	}{
		{
			name: "should return Id from simple content",
			content: `
NAME="Test OS"
VERSION="1.0"
ID=testos
ID_LIKE=anotheros
`,
			expectedOSId: "testos",
		},
		{
			name: "should return Id from double-quoted syntax with whitespaces",
			content: `
NAME="Test OS"
VERSION="1.0"
 ID = "testos"
ID_LIKE=anotheros
`,
			expectedOSId: "testos",
		},
		{
			name: "should return Id with whitespaces",
			content: `
NAME="Test OS"
VERSION="1.0"
 ID = testos
ID_LIKE=anotheros
`,
			expectedOSId: "testos",
		},
		{
			name: "should return Id from single-quoted syntax",
			content: `
NAME="Test OS"
VERSION="1.0"
ID='testos'
ID_LIKE=anotheros
`,
			expectedOSId: "testos",
		},
		{
			name: "should return the last Id when duplicates exist",
			content: `
ID=firstid
NAME="Test OS"
ID=secondid
VERSION="1.0"
ID=lastid
`,
			expectedOSId: "lastid",
		},
		{
			name:         "should return empty string for empty content",
			content:      "",
			expectedOSId: "",
		},
		{
			name: "should return empty string for content with no Id",
			content: `
NAME="Test OS"
VERSION="1.0"
`,
			expectedOSId: "",
		},
		{
			name: "should ignore comments",
			content: `
#ID=commented
NAME="Test OS"
ID=actual
`,
			expectedOSId: "actual",
		},
		{
			name: "should handle Id with embedded equals",
			content: `
ID=some=value
`,
			expectedOSId: "some=value",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			osId := system.ParseOSId([]byte(tc.content))
			if osId != tc.expectedOSId {
				t.Errorf("expected OSId %q, got %q", tc.expectedOSId, osId)
			}
		})
	}
}
