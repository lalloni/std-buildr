package version

import (
	"regexp"
	"strconv"

	"github.com/Masterminds/semver"
	"github.com/pkg/errors"
)

type EventualVersion struct {
	TrackerID  string
	IssueID    string
	Version    uint64
	Prerelease string
}

var eventualVersionRegex = regexp.MustCompile(`^(.*)\-(\d+)\-(\d+)(\-\d+\-.*)?$`)

func ParseSemanticVersion(version string) (*semver.Version, error) {
	sv, err := semver.NewVersion(version)
	if err != nil {
		return nil, errors.Wrap(err, "tag name must be a valid semver 2 string prefixed with a 'v' character")
	}
	return sv, nil
}

func ParseEventualVersion(version string) (*EventualVersion, error) {
	v := eventualVersionRegex.FindStringSubmatch(version)
	if v == nil {
		return nil, errors.Errorf("eventual version must match regular expression %q", eventualVersionRegex.String())
	}
	vn, err := strconv.ParseUint(v[3], 10, 64)
	if err != nil {
		return nil, errors.Errorf("eventual version must be a number instead of %q", v[3])
	}
	return &EventualVersion{
		TrackerID:  v[1],
		IssueID:    v[2],
		Version:    vn,
		Prerelease: v[4],
	}, nil
}
