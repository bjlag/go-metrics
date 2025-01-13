package cmd

import "fmt"

const noValue = "N/A"

type Build struct {
	version string
	date    string
	commit  string
}

func NewBuild(version, date, commit string) *Build {
	if version == "" {
		version = noValue
	}

	if date == "" {
		date = noValue
	}

	if commit == "" {
		commit = noValue
	}

	return &Build{
		version: version,
		date:    date,
		commit:  commit,
	}
}

func (b *Build) VersionString() string {
	return fmt.Sprintf("Build version: %s", b.version)
}

func (b *Build) DateString() string {
	return fmt.Sprintf("Build date: %s", b.date)
}

func (b *Build) CommitString() string {
	return fmt.Sprintf("Build commit: %s", b.commit)
}
