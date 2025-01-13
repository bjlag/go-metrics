package cmd_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bjlag/go-metrics/cmd"
)

func TestBuild_Build(t *testing.T) {
	type fields struct {
		version string
		date    string
		commit  string
	}

	tests := []struct {
		name        string
		fields      fields
		wantVersion string
		wantDate    string
		wantCommit  string
	}{
		{
			name: "has data",
			fields: fields{
				version: "1.0.0",
				date:    "2025/01/13 19:41:21",
				commit:  "1656c36",
			},
			wantVersion: "Build version: 1.0.0",
			wantDate:    "Build date: 2025/01/13 19:41:21",
			wantCommit:  "Build commit: 1656c36",
		},
		{
			name: "no data",
			fields: fields{
				version: "",
				date:    "",
				commit:  "",
			},
			wantVersion: "Build version: N/A",
			wantDate:    "Build date: N/A",
			wantCommit:  "Build commit: N/A",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := cmd.NewBuild(tt.fields.version, tt.fields.date, tt.fields.commit)

			assert.Equal(t, tt.wantVersion, b.VersionString())
			assert.Equal(t, tt.wantDate, b.DateString())
			assert.Equal(t, tt.wantCommit, b.CommitString())

		})
	}
}
